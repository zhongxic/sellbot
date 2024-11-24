package jieba

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetDAG(t *testing.T) {
	TestBefore(t)
	sentence := "我在学习结巴分词"
	DAG := map[int][]int{
		0: {0},
		1: {1},
		2: {2, 3},
		3: {3},
		4: {4, 5, 7},
		5: {5},
		6: {6, 7},
		7: {7},
	}
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	dag := tokenizer.getDAG(sentence)
	if !reflect.DeepEqual(dag, DAG) {
		t.Errorf("expected dag [%v] actual [%v]", DAG, dag)
	}
}

func TestCutAll(t *testing.T) {
	TestBefore(t)
	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	sentence := "我在学习结巴分词"
	expected := []string{"我", "在", "学习", "结巴", "结巴分词", "分词"}
	cuts := tokenizer.cutAll(sentence)
	if !reflect.DeepEqual(cuts, expected) {
		t.Errorf("expected cut [%v] actual [%v]", expected, cuts)
	}
}
