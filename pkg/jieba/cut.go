package jieba

func (t *Tokenizer) getDAG(sentence string) map[int][]int {
	DAG := map[int][]int{}
	runes := []rune(sentence)
	N := len(runes)
	for k := 0; k < N; k++ {
		tmplist := make([]int, 0)
		i := k
		frag := string(runes[k])
		freq, ok := t.freq[frag]
		for i < N && ok {
			if freq > 0 {
				tmplist = append(tmplist, i)
			}
			i += 1
			if i == N {
				break
			}
			frag = string(runes[k : i+1])
			freq, ok = t.freq[frag]
		}
		if len(tmplist) == 0 {
			tmplist = append(tmplist, k)
		}
		DAG[k] = tmplist
	}
	return DAG
}
