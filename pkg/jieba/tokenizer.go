package jieba

import (
	"bufio"
	_ "embed"
	"os"
	"strconv"
	"strings"

	"github.com/zhongxic/sellbot/pkg/container"
)

//go:embed dict.txt
var dict string

type Tokenizer struct {
	freq  *container.ConcurrentMap[string, int]
	total int
}

// NewTokenizer create a tokenizer with specific dict.
func NewTokenizer(dict string) (tokenizer *Tokenizer, err error) {
	f, err := os.Open(dict)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	return initialize(scanner)
}

// NewDefaultTokenizer create a tokenizer with embedded dict.
func NewDefaultTokenizer() (tokenizer *Tokenizer, err error) {
	scanner := bufio.NewScanner(strings.NewReader(dict))
	return initialize(scanner)
}

func initialize(scanner *bufio.Scanner) (tokenizer *Tokenizer, err error) {
	lfreq := container.NewConcurrentMap[string, int]()
	ltotal := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		split := strings.Split(line, " ")
		word := split[0]
		freq, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}
		lfreq.Put(word, freq)
		ltotal += freq
		runes := []rune(word)
		length := len(runes)
		for i := 0; i < length; i++ {
			wfrag := string(runes[:i+1])
			if _, ok := lfreq.Get(wfrag); !ok {
				lfreq.Put(wfrag, 0)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	seg := &Tokenizer{freq: lfreq, total: ltotal}
	return seg, nil
}

// AddWord add a word with specific frequency into the dict held by this tokenizer.
func (t *Tokenizer) AddWord(word string, frequency int) {
	freq, _ := t.freq.Get(word)
	nfreq := freq + frequency
	if nfreq <= 0 {
		t.freq.Put(word, 0)
		t.total -= freq
	} else {
		t.freq.Put(word, nfreq)
		t.total += frequency
	}
	runes := []rune(word)
	N := len(runes)
	for i := 0; i < N; i++ {
		frag := string(runes[:i+1])
		if _, ok := t.freq.Get(frag); !ok {
			t.freq.Put(frag, 0)
		}
	}
}

// DelWord del a word from the dict held by this tokenizer.
//
// If you want to decrease the frequency of this word rather than remove it from the dict entirely,
// please call AddWord with a negative frequency argument.
func (t *Tokenizer) DelWord(word string) {
	freq, _ := t.freq.Get(word)
	t.freq.Put(word, 0)
	t.total -= freq
}
