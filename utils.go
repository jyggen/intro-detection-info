package main

import (
	"fmt"
)

func missingEpisodeString(episodes []*Episode) string {
	l := len(episodes)

	if l == 0 {
		return "n/a"
	}

	if l == 1 {
		return fmt.Sprintf("Episode %d does not have an intro detected.", episodes[0].Number())
	}

	if l == 2 {
		return fmt.Sprintf("Episodes %d and %d do not have intros detected.", episodes[0].Number(), episodes[1].Number())
	}

	result := fmt.Sprintf("Episodes %d", episodes[0].Number())

	for _, episode := range episodes[1 : l-2] {
		result += fmt.Sprintf(", %d", episode.Number())
	}

	return result + fmt.Sprintf(" and %d do not have intros detected.", episodes[l-1].Number())
}
