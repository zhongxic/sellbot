# 结巴中文分词器

简易中文分词器。参考[结巴中文分词](https://github.com/fxsjy/jieba)项目并使用 go 实现，词典与模型文件均来自于此项目。相关原理与源码解析参考[机器翻译教程-Chapter2-中文分词](https://github.com/BrightXiaoHan/MachineTranslationTutorial/blob/master/tutorials/Chapter2/ChineseTokenizer.md)。

初始化分词器：

```go
// 使用默认词典初始化分词器
func NewDefaultTokenizer() (tokenizer *Tokenizer, err error)
// 使用自定义词典初始化分词器
func NewTokenizer(dict string) (tokenizer *Tokenizer, err error)
```

分词器提供的分词方法：

```go
// 全模式
// 从句子中切分出所有在词典中出现过的词，对于未在词典中出现的部分会被切分成单个字
func (t *Tokenizer) CutAll(sentence string) []string

// 新词发现模式
// 找出句子的最大切分路径，对于未出现在词典中的词，采用 HMM 和 viterbi 进行新词识别
func (t *Tokenizer) CutDAG(sentence string) []string

// 精确模式
// 找出句子的最大切分路径
func (t *Tokenizer) CutDAGNoHMM(sentence string) []string
```
