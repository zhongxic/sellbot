package jieba

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestCutAll(t *testing.T) {
	TestBefore(t)
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	sentence := "我在学习结巴分词"
	expected := []string{"我", "在", "学习", "结巴", "结巴分词", "分词"}
	cuts := tokenizer.CutAll(sentence)
	if !reflect.DeepEqual(cuts, expected) {
		t.Errorf("expected cut [%v] actual [%v]", expected, cuts)
	}
}

func TestCutDAGNoHMM(t *testing.T) {
	TestBefore(t)
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	sentence := "我在学习结巴分词"
	expected := []string{"我", "在", "学习", "结巴", "分词"}
	cuts := tokenizer.CutDAGNoHMM(sentence)
	if !reflect.DeepEqual(cuts, expected) {
		t.Errorf("expected cuts [%v] actual [%v]", expected, cuts)
	}
}
