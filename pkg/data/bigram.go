package data

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Bigram struct {
	pair string
	count int
}

func (b Bigram) Pair() string {
	return b.pair
}

func (b Bigram) Count() int {
	return b.count
}

const (
	startToken = 26
	endToken   = 27
)

func CharIndxToString(i int) string {
	switch i {
	case startToken:
		return "<S>"
	case endToken:
	 	return "<E>"
	default:
		return string(byte(i) + 'a')
	}
}

func Names() ([][]uint16, []Bigram) {
	f, err := os.Open("./names.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	counter := make([][]uint16, 28)
	for i := range counter {
		counter[i] = make([]uint16, 28)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name := scanner.Bytes()

		counter[startToken][name[0] - 'a'] += 1

		for i, b := range name[:len(name)-1] {
			ch1 := b - 'a'
			ch2 := name[i+1] - 'a'

			counter[ch1][ch2] += 1
		}

		counter[name[len(name)-1] - 'a'][endToken] += 1
	}

	var pairs []Bigram
	for i, ps := range counter {
		for j, count := range ps {
			if count < 1 {
				continue
			}
			pairs = append(pairs, Bigram{
				pair: CharIndxToString(i) + CharIndxToString(j),
				count: int(count),
			})
		}
	}

	slices.SortFunc(pairs, func(a Bigram, b Bigram) int {
		return a.count - b.count
	})

	for i := range pairs {
		if i == 25 {
			fmt.Printf("...\n")
			continue
		} else if i > 25 && i < len(pairs) - 25 {
			continue
		} 

		// fmt.Printf("%s => %d\n", p.pair, p.count)
	}

	os.Getwd()

	return counter, pairs
}
