package finalseg

import (
	"regexp"
	"sort"

	"github.com/zhongxic/sellbot/pkg/regex"
)

const minFloat = -3.14e100

var (
	reHan      = regexp.MustCompile(`(\p{Han}+)`)
	reSkip     = regexp.MustCompile(`([a-zA-Z0-9]+(?:\.\d+)?%?)`)
	prevStatus = map[byte][]byte{
		'B': {'E', 'S'},
		'M': {'M', 'B'},
		'S': {'S', 'E'},
		'E': {'B', 'M'}}
)

type stateProb struct {
	state byte
	p     float64
}

type probSlice []stateProb

func (s probSlice) Len() int {
	return len(s)
}

func (s probSlice) Less(i, j int) bool {
	return s[i].p < s[j].p
}

func (s probSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func Cut(sentence string) []string {
	res := make([]string, 0)
	blocks := regex.Split(sentence, reHan)
	for _, blk := range blocks {
		if reHan.MatchString(blk) {
			words := cut(blk)
			res = append(res, words...)
		} else {
			ss := regex.Split(blk, reSkip)
			for _, s := range ss {
				if s == "" {
					continue
				}
				res = append(res, s)
			}
		}
	}
	return res
}

func cut(sentence string) []string {
	res := make([]string, 0, 10)

	runes := []rune(sentence)
	_, positions := viterbi(runes, []byte{'B', 'M', 'E', 'S'})
	begin, next := 0, 0

	for i, char := range runes {
		pos := positions[i]
		switch pos {
		case 'B':
			begin = i
		case 'E':
			res = append(res, string(runes[begin:i+1]))
			next = i + 1
		case 'S':
			res = append(res, string(char))
			next = i + 1
		}
	}

	if next < len(runes) {
		res = append(res, string(runes[next:]))
	}

	return res
}

func viterbi(obs []rune, states []byte) (prob float64, positions []byte) { // NOSONAR
	V := make([]map[byte]float64, len(obs))
	path := make(map[byte][]byte)

	V[0] = make(map[byte]float64)
	for _, y := range states {
		emit, ok := probEmit[y][obs[0]]
		if !ok {
			emit = minFloat
		}
		V[0][y] = probStart[y] + emit
		path[y] = []byte{y}
	}

	N := len(obs)
	for t := 1; t < N; t++ {
		V[t] = make(map[byte]float64)
		newPath := make(map[byte][]byte)
		for _, y := range states {
			emit, ok := probEmit[y][obs[t]]
			if !ok {
				emit = minFloat
			}
			ps := make(probSlice, 0)
			for _, x := range prevStatus[y] {
				trans, ok := probTrans[x][y]
				if !ok {
					trans = minFloat
				}
				p := stateProb{x, V[t-1][x] + trans + emit}
				ps = append(ps, p)
			}
			sort.Sort(ps)
			mp := ps[len(ps)-1]
			V[t][y] = mp.p
			tp := make([]byte, len(path[mp.state]))
			copy(tp, path[mp.state])
			newPath[y] = append(tp, y)
		}
		path = newPath
	}

	ps := make(probSlice, 0)
	stats := []byte{'E', 'S'}
	for _, y := range stats {
		ps = append(ps, stateProb{y, V[N-1][y]})
	}
	sort.Sort(ps)
	mp := ps[len(ps)-1]
	return mp.p, path[mp.state]
}
