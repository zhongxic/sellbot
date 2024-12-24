package jieba

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testDict = "dict.txt"

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
	filename := filepath.Join("testdata", testDict)
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
	freq := map[string]int{
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
	total := 0
	for _, v := range freq {
		total += v
	}
	tokenizer, err := NewTokenizer(filepath.Join("testdata", testDict))
	if err != nil {
		t.Fatal(err)
	}
	if tokenizer == nil {
		t.Fatal("tokenizer expected not nil")
	}
	if tokenizer.freq == nil {
		t.Error("freq in tokenizer expected not nil")
	}
	for key, value := range tokenizer.freq {
		if value != freq[key] {
			t.Errorf("expected freq of word [%v] is [%v] actual [%v]", key, freq[key], value)
		}
	}
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

func TestLoadUserDict(t *testing.T) {
	userDict := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	dictFile := filepath.Join("testdata", "empty.txt")
	userDictFile := filepath.Join("testdata", "userdict.txt")
	err := os.MkdirAll(filepath.Dir(dictFile), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(dictFile, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}
	content := strings.Builder{}
	for word, freq := range userDict {
		content.WriteString(fmt.Sprintf("%v %v\n", word, freq))
	}
	err = os.WriteFile(userDictFile, []byte(content.String()), 0644)
	if err != nil {
		t.Fatal(err)
	}
	tokenizer, err := NewTokenizer(dictFile)
	if err != nil {
		t.Fatal(err)
	}
	err = tokenizer.LoadUserDict(userDictFile)
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range tokenizer.freq {
		if value != userDict[key] {
			t.Errorf("expected freq of word [%v] is [%v] actual [%v]", key, userDict[key], value)
		}
	}
}
