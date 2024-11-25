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

func TestCutDAG(t *testing.T) {
	TestBefore(t)
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	sentence := "这是一个伸手不见五指的黑夜我在学习结巴分词"
	expected := []string{"这是", "一个", "伸手", "不见", "五指", "的", "黑夜", "我", "在", "学习", "结巴", "分词"}
	cuts := tokenizer.CutDAG(sentence)
	if !reflect.DeepEqual(cuts, expected) {
		t.Errorf("expected cuts [%v] actual [%v]", expected, cuts)
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
