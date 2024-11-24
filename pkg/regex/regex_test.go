package regex

import (
	"regexp"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	s := "这是一个伸手不见五指的黑夜。我叫孙悟空，我爱北京，我爱Python和C++。"
	re := regexp.MustCompile(`([\p{Han}a-zA-Z0-9+#&._%\-]+)`)
	splits := Split(s, re)
	joined := strings.Join(splits, "")
	if joined != s {
		t.Errorf("expected s after joined [%v] actual [%v]", s, joined)
	}
}
