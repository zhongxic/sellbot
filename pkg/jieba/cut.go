package jieba

import (
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/zhongxic/sellbot/pkg/regex"
)

var reEng = regexp.MustCompile(`[a-zA-Z0-9]`)

type edge struct {
	weight float64
	index  int
}

type cutFunc func(string) []string

type edgeSlice []edge

func (s edgeSlice) Len() int {
	return len(s)
}

func (s edgeSlice) Less(i, j int) bool {
	return s[i].weight < s[j].weight
}

func (s edgeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// cut is a preprocessing method that uses regular expressions to split this sentence into blocks.
//
// Cutting block into chinese words is delegated to cutFunc.
func (t *Tokenizer) cut(sentence string, reHan, reSkip *regexp.Regexp, cutFunc cutFunc) []string { //NOSONAR
	res := make([]string, 0)
	blocks := regex.Split(sentence, reHan)
	for _, blk := range blocks {
		if blk == "" {
			continue
		}
		if reHan.MatchString(blk) {
			words := cutFunc(blk)
			res = append(res, words...)
		} else {
			ss := regex.Split(blk, reSkip)
			for _, s := range ss {
				if reSkip.MatchString(s) {
					res = append(res, s)
				} else {
					for _, x := range s {
						res = append(res, string(x))
					}
				}
			}
		}
	}
	return res
}

func (t *Tokenizer) cutAll(sentence string) []string { //NOSONAR
	words := make([]string, 0)

	DAG := t.getDAG(sentence)
	shadow := -1
	engScan := false
	engBuf := &strings.Builder{}
	runes := []rune(sentence)
	N := len(runes)
	for k := 0; k < N; k++ {
		if engScan && !reEng.MatchString(string(runes[k])) {
			engScan = false
			words = append(words, engBuf.String())
			engBuf.Reset()
		}
		L := DAG[k]
		if len(L) == 1 && k > shadow {
			word := string(runes[k : L[0]+1])
			if reEng.MatchString(word) {
				if !engScan {
					engScan = true
				}
				engBuf.WriteString(word)
			} else {
				words = append(words, word)
			}
			shadow = L[0]
		} else {
			for _, j := range L {
				if j > k {
					word := string(runes[k : j+1])
					words = append(words, word)
					shadow = j
				}
			}
		}
	}
	if engScan {
		words = append(words, engBuf.String())
	}
	return words
}

func (t *Tokenizer) cutDAGNoHHM(sentence string) []string {
	words := make([]string, 0)
	DAG := t.getDAG(sentence)
	route := t.calc(sentence, DAG)
	engBuf := &strings.Builder{}
	runes := []rune(sentence)
	N := len(runes)
	for x := 0; x < N; {
		y := route[x].index + 1
		word := string(runes[x:y])
		if reEng.MatchString(word) {
			engBuf.WriteString(word)
		} else {
			if engBuf.Len() > 0 {
				words = append(words, engBuf.String())
				engBuf.Reset()
			}
			words = append(words, word)
		}
		x = y
	}
	return words
}

func (t *Tokenizer) getDAG(sentence string) map[int][]int {
	DAG := map[int][]int{}
	runes := []rune(sentence)
	N := len(runes)
	for k := 0; k < N; k++ {
		tmplist := make([]int, 0)
		i := k
		frag := string(runes[k])
		freq, ok := t.freq[frag]
		for i < N && ok {
			if freq > 0 {
				tmplist = append(tmplist, i)
			}
			i += 1
			if i == N {
				break
			}
			frag = string(runes[k : i+1])
			freq, ok = t.freq[frag]
		}
		if len(tmplist) == 0 {
			tmplist = append(tmplist, k)
		}
		DAG[k] = tmplist
	}
	return DAG
}

func (t *Tokenizer) calc(sentence string, DAG map[int][]int) map[int]edge {
	runes := []rune(sentence)
	N := len(runes)
	route := map[int]edge{
		N: {0, 0},
	}
	logtotal := math.Log(float64(t.total))
	for idx := N - 1; idx >= 0; idx-- {
		edges := make(edgeSlice, 0)
		for _, x := range DAG[idx] {
			frag := string(runes[idx : x+1])
			freq, ok := t.freq[frag]
			if !ok {
				freq = 1
			}
			weight := math.Log(float64(freq)) - logtotal + route[x+1].weight
			edges = append(edges, edge{weight, x})
		}
		sort.Sort(edges)
		route[idx] = edges[len(edges)-1]
	}
	return route
}
