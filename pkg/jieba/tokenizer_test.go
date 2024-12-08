package jieba

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestBefore generate dict in testdata dir before test.
func TestBefore(t *testing.T) {
	dict := `
	我 123
	在 234
	学习 456
	结巴 345
	分词 456
	结巴分词 23
	学 2344
	分 23
	结 234
	`
	filename := filepath.Join("testdata", "dict.txt")
	err := os.MkdirAll(filepath.Dir(filename), 0644)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(strings.NewReader(dict))
	for scanner.Scan() {
		if _, err := fmt.Fprintln(f, strings.TrimSpace(scanner.Text())); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNewTokenizer(t *testing.T) {
	TestBefore(t)
	freq := map[string]int64{
		"我":    123,
		"在":    234,
		"学":    2344,
		"学习":   456,
		"结":    234,
		"结巴":   345,
		"结巴分":  0,
		"结巴分词": 23,
		"分":    23,
		"分词":   456,
	}
	var total int64
	for _, v := range freq {
		total += v
	}
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if tokenizer == nil {
		t.Fatal("tokenizer expected not nil")
	}
	if tokenizer.freq == nil {
		t.Error("freq in tokenizer expected not nil")
	}
	tokenizer.freq.Range(func(key, value any) bool {
		if value.(int64) != freq[key.(string)] {
			t.Errorf("expected freq of word [%v] is [%v] actual [%v]", key, freq[key.(string)], value)
		}
		return true
	})
	if tokenizer.total != total {
		t.Errorf("expected total in tokenizer [%v] actual [%v]", total, tokenizer.total)
	}
}

func TestNewDefaultTokenizer(t *testing.T) {
	tokenizer, err := NewDefaultTokenizer()
	if err != nil {
		t.Fatal(err)
	}
	if tokenizer == nil {
		t.Fatal("tokenizer expected not nil")
	}
	if tokenizer.freq == nil {
		t.Error("freq in tokenizer expected not nil")
	}
}

func TestAddDelWordConcurrent(t *testing.T) {
	TestBefore(t)
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	N := 100
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			tokenizer.AddWord("结巴", 1)
			wg.Done()
		}()
	}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			tokenizer.DelWord("结巴")
			wg.Done()
		}()
	}
	wg.Wait()
}
