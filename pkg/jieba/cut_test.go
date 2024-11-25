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

func TestCalc(t *testing.T) {
	TestBefore(t)
	sentence := "我在学习结巴分词"
	expected := map[int]edge{
		8: {0, 0},
		7: {-8.351846738828245, 7},
		6: {-2.2293539293138585, 7},
		5: {-10.581200668142102, 5},
		4: {-4.737656251110743, 5},
		3: {-13.089502989938989, 3},
		2: {-6.967010180424602, 3},
		1: {-9.863535803895145, 1},
		0: {-13.403198187350974, 0},
	}

	tokenizer, err := NewTokenizer(filepath.Join("testdata", "dict.txt"))
	if err != nil {
		t.Fatal(err)
	}
	DAG := tokenizer.getDAG(sentence)
	calc := tokenizer.calc(sentence, DAG)
	if !reflect.DeepEqual(calc, expected) {
		t.Errorf("expected calc [%v] actual [%v]", expected, calc)
	}
}
