package utils

import (
	"strconv"
	"strings"

	"github.com/tx7do/go-utils/sliceutil"
)

func FilterBlacklist(data []string, blacklist []string) []string {
	bm := make(map[string]struct{}, len(blacklist))
	for _, s := range blacklist {
		bm[s] = struct{}{}
	}

	n := 0
	for _, x := range data {
		if _, found := bm[x]; !found {
			data[n] = x
			n++
		}
	}
	return data[:n]
}

func NumberSliceToString(numbers []int64) string {
	return strings.Join(
		sliceutil.Map(numbers, func(value int64, _ int, _ []int64) string { return strconv.FormatInt(value, 10) }),
		",",
	)
}
