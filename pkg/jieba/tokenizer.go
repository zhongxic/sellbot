package jieba

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zhongxic/sellbot/pkg/container"
)

//go:embed dict.txt
var dict string

type Tokenizer struct {
	freq  *container.ConcurrentMap[string, int64]
	total int64
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
	var lfreq = container.NewConcurrentMap[string, int64]()
	var ltotal int64 = 0
	var start = time.Now()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		split := strings.Split(line, " ")
		word := split[0]
		freq, err := strconv.ParseInt(split[1], 10, 64)
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
	slog.Info("tokenizer initialized", "cost", time.Since(start).Milliseconds())
	seg := &Tokenizer{freq: lfreq, total: ltotal}
	return seg, nil
}

// AddWord add a word with specific frequency into the dict held by this tokenizer.
func (t *Tokenizer) AddWord(word string, frequency int64) {
	freq, _ := t.freq.Get(word)
	if freq+frequency <= 0 {
		t.freq.Put(word, 0)
		t.total -= freq
	} else {
		t.freq.Put(word, freq+frequency)
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

func (t *Tokenizer) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	freq := make(map[string]int64)
	t.freq.Range(func(key, value any) bool {
		freq[key.(string)] = value.(int64)
		return true
	})
	m["freq"] = freq
	m["total"] = t.total
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *Tokenizer) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	m := map[string]any{}
	if err := dec.Decode(&m); err != nil {
		return err
	}
	wfreq := m["freq"].(map[string]any)
	lfreq := container.NewConcurrentMap[string, int64]()
	for word, freq := range wfreq {
		n, err := freq.(json.Number).Int64()
		if err != nil {
			return err
		}
		lfreq.Put(word, n)
	}
	ltotal, err := m["total"].(json.Number).Int64()
	if err != nil {
		return err
	}
	t.freq = lfreq
	t.total = ltotal
	return nil
}
