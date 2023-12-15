package challenge

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Boxes [256][]*Lense

func (b *Boxes) Add(idx int, l *Lense) {
	found := false
	for i, fl := range b[idx] {
		if fl == nil {
			continue
		}
		if l.Label == fl.Label {
			b[idx][i] = l
			found = true
			break
		}
	}

	if !found {
		b[idx] = append(b[idx], l)
	}
	b.reslice(idx)
}

func (b *Boxes) Remove(idx int, label string) {
	for i, l := range b[idx] {
		if l == nil {
			continue
		}
		if l.Label == label {
			b[idx][i] = nil
		}
	}

	b.reslice(idx)
}

func (b *Boxes) reslice(idx int) {
	newS := []*Lense{}

	for _, l := range b[idx] {
		if l != nil {
			newS = append(newS, l)
		}
	}

	b[idx] = newS
}

func (b *Boxes) String() string {
	buf := bytes.NewBuffer(nil)
	for i := 0; i < len(b); i++ {
		if len(b[i]) == 0 {
			continue
		}
		fmt.Fprintf(buf, "%d: [", i)
		for j := 0; j < len(b[i]); j++ {
			fmt.Fprintf(buf, "(%s, %d)", b[i][j].Label, b[i][j].FocalLength)
		}
		fmt.Fprintln(buf, "]")
	}
	return buf.String()
}

type Lense struct {
	Label       string
	FocalLength int
}

func (b Boxes) FocusingPower() int {
	totalFp := 0
	for i := 0; i < len(b); i++ {
		b.reslice(i)
		for j := 0; j < len(b[i]); j++ {
			totalFp += (i + 1) * (j + 1) * b[i][j].FocalLength
		}
	}
	return totalFp
}

func Challenge(input string) string {
	items := strings.Split(input, ",")

	boxes := Boxes{}

	for _, item := range items {
		letterCode := ""

		restStartIdx := 0

		for _, c := range item {
			if c >= 'a' && c <= 'z' {
				letterCode += string(c)
				restStartIdx++
			} else {
				break
			}
		}
		boxCode := hashAlgorithm(letterCode)
		switch item[restStartIdx] {
		case '=':
			numberStr := string(item[restStartIdx+1:])
			vInt, err := strconv.Atoi(numberStr)
			if err != nil {
				panic("int conv error")
			}
			lense := &Lense{
				FocalLength: vInt,
				Label:       letterCode,
			}
			boxes.Add(boxCode, lense)
		case '-':
			boxes.Remove(boxCode, letterCode)
		}
	}

	// fmt.Println(boxes.String())
	return strconv.Itoa(boxes.FocusingPower())
}

func hashAlgorithm(input string) int {
	var val int

	for _, c := range input {
		if c == '\n' {
			continue
		}
		val += int(c)
		val *= 17
		val = val % 256
	}

	return val
}
