package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regSplitByDelimiters = regexp.MustCompile(`(?m)([\p{L}][\p{L}-]*)`)

type word struct {
	w string
	n int
}

func Top10(str string) []string {
	var result []string

	m := make(map[string]int)

	result = regSplitByDelimiters.FindAllString(str, -1)
	for i := 0; i < len(result); i++ {
		m[strings.ToLower(result[i])]++
	}

	sw := make([]word, 0)
	for key, value := range m {
		sw = append(sw, word{w: key, n: value})
	}

	sort.Slice(sw, func(i, j int) bool {
		r := false
		if sw[i].n > sw[j].n {
			r = true
		} else if sw[i].n == sw[j].n && sw[i].w < sw[j].w {
			r = true
		}
		return r
	})

	result = nil
	for i := 0; i < 10 && i < len(sw); i++ {
		result = append(result, sw[i].w)
	}

	return result
}
