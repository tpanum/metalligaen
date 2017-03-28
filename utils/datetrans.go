package utils

import (
	"strings"
)

var dic = map[string]string{
	"January":  "Januar",
	"February": "Februar",
	"March":    "Marts",
	"May":      "Maj",
	"June":     "Juni",
	"July":     "Juli",
	"October":  "Oktober",

	"Monday":    "Mandag",
	"Tuesday":   "Tirsdag",
	"Wednesday": "Onsday",
	"Thursday":  "Torsdag",
	"Friday":    "Fredag",
	"Saturday":  "Lørdag",
	"Sunday":    "Søndag",
}

func TranslateTimeDA(day string) string {
	for search, replace := range dic {
		day = strings.Replace(day, search, replace, -1)
	}

	return day
}
