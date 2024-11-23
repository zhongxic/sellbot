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
)

//go:embed dict.txt
var dict string

type Tokenizer struct {
	freq  map[string]int64
	total int64
}

// NewTokenizer create a tokenizer with specific dict
func NewTokenizer(dict string) (tokenizer *Tokenizer, err error) {
	f, err := os.Open(dict)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	return initialize(scanner)
}

// NewDefaultTokenizer create a tokenizer with embedded dict
func NewDefaultTokenizer() (tokenizer *Tokenizer, err error) {
	scanner := bufio.NewScanner(strings.NewReader(dict))
	return initialize(scanner)
}

func initialize(scanner *bufio.Scanner) (tokenizer *Tokenizer, err error) {
	var lfreq = map[string]int64{}
	var ltotal int64 = 0
	var start = time.Now()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, " ")
		word := split[0]
		freq, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return nil, err
		}
		lfreq[word] = freq
		ltotal += freq
		runes := []rune(word)
		length := len(runes)
		for i := 0; i < length; i++ {
			wfrag := string(runes[:i+1])
			if _, ok := lfreq[wfrag]; !ok {
				lfreq[wfrag] = 0
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

func (t *Tokenizer) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	m["freq"] = t.freq
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
	lfreq := make(map[string]int64, len(wfreq))
	for word, freq := range wfreq {
		n, err := freq.(json.Number).Int64()
		if err != nil {
			return err
		}
		lfreq[word] = n
	}
	ltotal, err := m["total"].(json.Number).Int64()
	if err != nil {
		return err
	}
	t.freq = lfreq
	t.total = ltotal
	return nil
}
