package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gammazero/workerpool"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strconv"
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

func main() {
	name := os.Args[0]
	args := os.Args[1:]

	if len(args) != 2 {
		errorAndExit(name, fmt.Errorf("exactly 2 command-line arguments expected, %d received", len(args)))
	}

	plex, err := NewPlex(args[0], args[1])

	if err != nil {
		errorAndExit(name, err)
	}

	wp := workerpool.New(10)
	shows, err := plex.Shows()

	if err != nil {
		errorAndExit(name, err)
	}

	s := spinner.New(
		spinner.CharSets[14],
		250*time.Millisecond,
		spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Fetching metadata..."),
		spinner.WithFinalMSG("Metadata fetched.\n"),
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
		errorAndExit(name, err)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetHeader([]string{"Show", "Season", "Detected", "Comment"})
	table.SetRowLine(true)

	for _, result := range results {
		var status string
		var color int
		var missing string

		if result.DetectedEpisodesCount == 0 {
			status = "No"
			color = tablewriter.FgRedColor
		} else if result.DetectedEpisodesCount == result.TotalEpisodesCount {
			status = "Yes"
			color = tablewriter.FgGreenColor
		} else {
			status = "Partial"
			color = tablewriter.FgYellowColor
			missing = missingEpisodeString(result.MissingEpisodesList)
		}

		table.Rich(
			[]string{result.Show.Title(), strconv.Itoa(result.Season.Number()), status, missing},
			[]tablewriter.Colors{{}, {}, {tablewriter.Normal, color}},
		)
	}

	s.Stop()
	table.Render()
}
