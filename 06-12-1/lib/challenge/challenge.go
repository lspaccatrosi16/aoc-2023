package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

const ACCEL = 1

type Race struct {
	Time       int
	RecordDist int
}

func (r *Race) Combinations() int {
	winningCombos := 0
	for i := 0; i < r.Time; i++ {
		v := float64(i)
		t := float64(r.Time - i)

		s := v * t
		if s > float64(r.RecordDist) {
			winningCombos++
		}
	}
	return winningCombos
}

func Challenge(input string) string {
	races := []Race{}

	for _, l := range strings.Split(input, "\n") {
		cols := strings.Split(l, " ")
		nums := []int{}

		for i, c := range cols {
			if i < 1 || c == "" {
				continue
			}
			vInt, err := strconv.Atoi(c)
			if err != nil {
				fmt.Println(c)
				panic("int conv error")
			}
			nums = append(nums, vInt)
		}

		for i, n := range nums {
			if i == len(races) {
				races = append(races, Race{})
			}

			if races[i].Time != 0 {
				races[i].RecordDist = n
			} else {
				races[i] = Race{Time: n}
			}
		}
	}

	tally := 1

	for _, r := range races {
		fmt.Println("Race", r.RecordDist, r.Time)
		tally *= r.Combinations()
	}

	return strconv.Itoa(tally)
}
