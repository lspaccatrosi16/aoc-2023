package challenge

import (
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/algorithm"
)

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	tally := 0

	for _, l := range lines {
		terms := strings.Split(l, " ")

		nums := []int{}

		for _, t := range terms {
			if t == "" {
				continue
			}

			vInt, err := strconv.Atoi(t)
			if err != nil {
				panic("int conv error")
			}

			nums = append(nums, vInt)
		}
		seq := algorithm.SolveSequence(nums...)

		newVal := seq.Get(len(nums) + 1)

		tally += newVal
	}

	return strconv.Itoa(tally)
}
