package jieba

import (
	"regexp"
	"strings"
)

var (
	reHan  = regexp.MustCompile(`([\p{Han}a-zA-Z0-9+#&._%\-]+)`)
	reEng  = regexp.MustCompile(`[a-zA-Z0-9]`)
	reSkip = regexp.MustCompile(`(\r\n|\s)`)
)

// CutAll slices sentence into separated words.
//
// This method try to match all words contained in the dictionary,
// those parts who not appear in dictionary will be cut into single character.
func (t *Tokenizer) CutAll(sentence string) []string { //NOSONAR
	res := make([]string, 0)
	blocks := Split(sentence, reHan)
	for _, blk := range blocks {
		if blk == "" {
			continue
		}
		if reHan.MatchString(blk) {
			words := t.cutAll(blk)
			for _, word := range words {
				res = append(res, word)
			}
		} else {
			ss := Split(blk, reSkip)
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
