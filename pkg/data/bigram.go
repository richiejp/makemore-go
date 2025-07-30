package data

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Bigram struct {
	Ch1   int
	Ch2   int
	Count int
}

func (b Bigram) String() string {
	return fmt.Sprintf("%v%v", CharIndxToString(b.Ch1), CharIndxToString(b.Ch2))
}

const (
	StartToken = 26
	EndToken   = 27
)

func CharIndxToString(i int) string {
	switch i {
	case StartToken:
		return "<S>"
	case EndToken:
		return "<E>"
	default:
		return string(byte(i) + 'a')
	}
}

func NamesScanner() (*os.File, *bufio.Scanner) {
	f, err := os.Open("./names.txt")
	if err != nil {
		panic(err)
	}

	return f, bufio.NewScanner(f)
}

func BigramVisitor(fn func(ch1, ch2 int) bool) {
	f, s := NamesScanner()
	defer f.Close()

	for s.Scan() {
		name := s.Bytes()

		ch1 := int(StartToken)
		ch2 := int(name[0] - 'a')
		for j := range len(name) + 1 {
			if !fn(ch1, ch2) {
				return
			}

			ch1 = ch2
			if j+1 < len(name) {
				ch2 = int(name[j+1] - 'a')
			} else {
				ch2 = EndToken
			}
		}
	}
}

func Names(smoothing int) ([][]uint16, []Bigram) {
	counter := make([][]uint16, 28)
	for i := range counter {
		counter[i] = make([]uint16, 28)

		for j := range counter[i] {
			counter[i][j] = uint16(smoothing)
		}
	}

	BigramVisitor(func(ch1, ch2 int) bool {
		counter[ch1][ch2] += 1

		return true
	})

	var pairs []Bigram
	for i, ps := range counter {
		for j, count := range ps {
			if count < 1 {
				continue
			}
			pairs = append(pairs, Bigram{
				Ch1:   i,
				Ch2:   j,
				Count: int(count),
			})
		}
	}

	slices.SortFunc(pairs, func(a Bigram, b Bigram) int {
		return a.Count - b.Count
	})

	for i := range pairs {
		if i == 25 {
			fmt.Printf("...\n")
			continue
		} else if i > 25 && i < len(pairs)-25 {
			continue
		}

		// fmt.Printf("%s => %d\n", p.pair, p.count)
	}

	os.Getwd()

	return counter, pairs
}
