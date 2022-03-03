package main

import (
	"encoding/csv"
	"github.com/mattn/go-colorable"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

var formats = []string{"csv", "ascii"}

func isValidFormat(format string) bool {
	for _, allowed := range formats {
		if allowed == format {
			return true
		}
	}

	return false
}

func outputAsAscii(results []*Result) error {
	var table *tablewriter.Table

	if colors {
		table = tablewriter.NewWriter(colorable.NewColorableStdout())
	} else {
		table = tablewriter.NewWriter(colorable.NewNonColorable(os.Stdout))
	}

	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetHeader([]string{"Show", "Season", "Detected", "Comment"})
	table.SetRowLine(true)

	for _, result := range results {
		var color int

		if result.DetectedEpisodesCount == 0 {
			color = tablewriter.FgRedColor
		} else if result.DetectedEpisodesCount == result.TotalEpisodesCount {
			color = tablewriter.FgGreenColor
		} else {
			color = tablewriter.FgYellowColor
		}

		table.Rich(
			[]string{
				result.Show.Title(),
				strconv.Itoa(result.Season.Number()),
				seasonStatusString(result),
				missingEpisodeString(result.MissingEpisodesList),
			},
			[]tablewriter.Colors{{}, {}, {tablewriter.Normal, color}},
		)
	}

	table.Render()

	return nil
}

func outputAsCsv(results []*Result) error {
	rows := make([][]string, len(results)+1)
	rows[0] = []string{"Show", "Season", "Detected", "Comment"}

	for index, result := range results {
		rows[index+1] = []string{
			result.Show.Title(),
			strconv.Itoa(result.Season.Number()),
			seasonStatusString(result),
			missingEpisodeString(result.MissingEpisodesList),
		}
	}

	w := csv.NewWriter(os.Stdout)

	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}
