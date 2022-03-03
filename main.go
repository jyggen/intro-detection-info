package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gammazero/workerpool"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

type Result struct {
	Show                  *Show
	Season                *Season
	TotalEpisodesCount    int
	DetectedEpisodesCount int
	MissingEpisodesList   []*Episode
}

var colors bool
var format string
var timeout int

func main() {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !isValidFormat(format) {
				return fmt.Errorf("\"%s\" is not a known output format, expected one of %s", format, formats)
			}

			return nil
		},
		RunE:    run,
		Use:     os.Args[0],
		Version: getBuildInfo(),
	}

	cmd.PersistentFlags().BoolVar(&colors, "colors", true, "whether colors should be used in output or not")
	cmd.PersistentFlags().StringVar(&format, "format", "ascii", fmt.Sprintf("preferred output format, should be one of %s", formats))
	cmd.PersistentFlags().IntVar(&timeout, "timeout", 10, "timeout for all HTTP requests, in seconds")

	err := cmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	plex, err := NewPlex(args[0], args[1], time.Second*time.Duration(timeout))

	if err != nil {
		return err
	}

	wp := workerpool.New(10)
	shows, err := plex.Shows()

	if err != nil {
		return err
	}

	var writer io.Writer

	if colors {
		writer = colorable.NewColorableStdout()
	} else {
		writer = colorable.NewNonColorable(os.Stdout)
	}

	s := spinner.New(
		spinner.CharSets[14],
		250*time.Millisecond,
		spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Fetching metadata..."),
		spinner.WithFinalMSG("Metadata fetched.\n"),
		spinner.WithWriter(writer),
	)

	s.Start()

	errorsChan := make(chan error)
	resultsChan := make(chan *Result)
	errors := make([]error, 0)
	results := make([]*Result, 0)

	var channels sync.WaitGroup

	channels.Add(2)

	go func() {
		defer channels.Done()

		for err := range errorsChan {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer channels.Done()

		for result := range resultsChan {
			results = append(results, result)
		}
	}()

	for _, show := range shows {
		show := show
		wp.Submit(func() {
			seasons, err := show.Seasons()

			if err != nil {
				errorsChan <- err
				return
			}

			for _, season := range seasons {
				episodes, err := season.Episodes()

				if err != nil {
					errorsChan <- err
					continue
				}

				result := &Result{
					Show:                show,
					Season:              season,
					MissingEpisodesList: make([]*Episode, 0),
				}

				for _, episode := range episodes {
					result.TotalEpisodesCount++

					if episode.HasIntroMarker() {
						result.DetectedEpisodesCount++
					} else {
						result.MissingEpisodesList = append(result.MissingEpisodesList, episode)
					}
				}

				resultsChan <- result
			}
		})
	}

	wp.StopWait()

	close(errorsChan)
	close(resultsChan)

	channels.Wait()

	for _, err := range errors {
		return err
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Show.SortTitle() == results[j].Show.SortTitle() {
			if results[i].Show.RatingKey() == results[j].Show.RatingKey() {
				return results[i].Season.Number() < results[j].Season.Number()
			}

			return results[i].Show.RatingKey() < results[j].Show.RatingKey()
		}

		return results[i].Show.SortTitle() < results[j].Show.SortTitle()
	})

	s.Stop()

	if format == "ascii" {
		err = outputAsAscii(results)
	} else if format == "csv" {
		err = outputAsCsv(results)
	}

	return err
}
