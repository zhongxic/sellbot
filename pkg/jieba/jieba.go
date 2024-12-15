package jieba

import "regexp"

var (
	reHan  = regexp.MustCompile(`([\p{Han}a-zA-Z0-9+#&._%\-]+)`)
	reSkip = regexp.MustCompile(`(\r\n|\s)`)
)

// CutAll slices sentence into separated words.
//
// This method try to match all words contained in the dictionary,
// those parts who not appear in dictionary will be cut into single character.
func (t *Tokenizer) CutAll(sentence string) []string {
	return t.cut(sentence, reHan, reSkip, t.cutAll)
}

// CutDAG slices sentence into separated words with HMM.
func (t *Tokenizer) CutDAG(sentence string) []string {
	return t.cut(sentence, reHan, reSkip, t.cutDAG)
}

// CutDAGNoHMM slices sentence into separated words without HMM.
func (t *Tokenizer) CutDAGNoHMM(sentence string) []string {
	return t.cut(sentence, reHan, reSkip, t.cutDAGNoHHM)
}

// CutAllRegex slices sentence into separated words. use this regex to split string into blocks.
func (t *Tokenizer) CutAllRegex(sentence string, reHan, reSkip *regexp.Regexp) []string {
	return t.cut(sentence, reHan, reSkip, t.cutAll)
}

// CutDAGRegex slices sentence into separated words with HMM. use this regex to split string into blocks.
func (t *Tokenizer) CutDAGRegex(sentence string, reHan, reSkip *regexp.Regexp) []string {
	return t.cut(sentence, reHan, reSkip, t.cutDAG)
}

// CutDAGNoHMMRegex slices sentence into separated words without HMM. use this regex to split string into blocks.
func (t *Tokenizer) CutDAGNoHMMRegex(sentence string, reHan, reSkip *regexp.Regexp) []string {
	return t.cut(sentence, reHan, reSkip, t.cutDAGNoHHM)
}
