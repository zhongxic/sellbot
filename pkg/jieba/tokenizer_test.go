package jieba

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	tokenizer, err := NewTokenizer("dict.txt")
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

func TestMarshalJSON(t *testing.T) {
	tokenizer, err := NewDefaultTokenizer()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = json.Marshal(tokenizer); err != nil {
		t.Fatal(err)
	}
}

func TestUnMarshalJSON(t *testing.T) {
	tokenizer, err := NewDefaultTokenizer()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(tokenizer)
	if err != nil {
		t.Fatal(err)
	}

	seg := &Tokenizer{}
	if err := json.Unmarshal(data, seg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tokenizer, seg) {
		t.Fatal("tokenizer not deep equal after unmarshal")
	}
}
