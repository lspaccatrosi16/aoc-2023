package challenge

import (
	"fmt"
	"strconv"
	"strings"
)

type ValMap struct {
	Rules []ValStruct
}

func (v *ValMap) Calc(i int) int {
	for _, r := range v.Rules {
		prod := r.Calc(i)
		if prod != i {
			return prod
		}
	}
	return i
}

type ValStruct struct {
	DstStart int
	SrcStart int
	RangeLen int
}

func (v *ValStruct) Calc(i int) int {
	srcEnd := v.SrcStart + v.RangeLen
	if i < v.SrcStart {
		return i
	} else if i < srcEnd {
		return (i - v.SrcStart) + v.DstStart
	} else {
		return i
	}
}

func Challenge(input string) string {

	overallMap := map[string]*ValMap{}

	lines := strings.Split(input, "\n")

	currentKey := ""
	seeds := []int{}

	mapList := []string{}

	for _, l := range lines {
		// fmt.Printf("Line %d\n", i+1)
		if l == "" {
			continue
		}
		if strings.Contains(l, ":") {
			// fmt.Println("Key Line")
			if strings.HasPrefix(l, "seeds:") {
				parts := strings.Split(l, ":")
				nums := strings.Split(parts[1], " ")

				for _, n := range nums {
					if n == "" {
						continue
					}

					vInt, err := strconv.Atoi(n)
					if err != nil {
						panic("int conv error")
					}
					seeds = append(seeds, vInt)
				}
			} else {
				parts := strings.Split(l, ":")
				words := strings.Split(parts[0], " ")
				if words[1] != "map" {
					panic("expected map")
				}
				names := strings.Split(words[0], "-")
				if names[1] != "to" {
					panic("expected to")
				}
				currentKey = fmt.Sprintf("%s-%s", names[0], names[2])
				mapList = append(mapList, currentKey)
			}
		} else {
			// fmt.Println("Number Line")
			nums := strings.Split(l, " ")
			// fmt.Println(nums)
			n := []int{}
			for _, num := range nums {
				if num == "" {
					continue
				}
				vInt, err := strconv.Atoi(num)
				if err != nil {
					panic("int conv error")
				}
				n = append(n, vInt)
			}

			if len(n) != 3 {
				panic("incorrect len number list")
			}

			existing := overallMap[currentKey]
			overallMap[currentKey] = produceRange(n[0], n[1], n[2], existing)
		}
	}

	path := findPath(mapList, "seed-location")

	locs := []int{}
	minLoc := 1000000000000000000

	for _, seed := range seeds {
		locs = append(locs, followPath(overallMap, path, seed))
	}

	for _, loc := range locs {
		if loc < minLoc {
			minLoc = loc
		}
	}

	return strconv.Itoa(minLoc)
}

func followPath(m map[string]*ValMap, path []string, input int) int {
	val := input

	for _, p := range path {
		val = m[p].Calc(val)
	}

	return val
}

func findPath(list []string, required string) []string {
	order := []string{}

	for _, path := range list {
		if path == required {
			order = append(order, path)
			return order
		}
	}

	comps := strings.Split(required, "-")
	src := comps[0]
	dst := comps[1]

	for _, path := range list {
		if path == required {
			order = append(order, path)
			return order
		}

		subComps := strings.Split(path, "-")

		subSrc := subComps[0]
		subDst := subComps[1]

		if subDst == dst {
			newKey := fmt.Sprintf("%s-%s", src, subSrc)
			order = append(findPath(list, newKey), path)
			return order
		}
	}

	panic("could not find a path")
}

func produceRange(destRange, sourceRange, rangeLen int, m *ValMap) *ValMap {
	if m == nil {
		m = &ValMap{}
	}

	rule := ValStruct{destRange, sourceRange, rangeLen}

	m.Rules = append(m.Rules, rule)

	return m
}
