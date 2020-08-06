package main

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func outputAsTable(results []*Result) {
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

	table.Render()
}
