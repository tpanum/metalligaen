package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	clockParser = regexp.MustCompile(`([0-9]{2}):([0-9]{2})`)
)

func GameTimeToSeconds(s string) (uint, error) {
	matches := clockParser.FindStringSubmatch(s)

	if len(matches) == 3 {
		mins, _ := strconv.Atoi(matches[1])
		secs, _ := strconv.Atoi(matches[2])

		var duration uint = uint(mins)*60 + uint(secs)

		return duration, nil
	}

	return 0, fmt.Errorf("Unable to parse game time (%s)", s)
}
