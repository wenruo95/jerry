/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : dirty_words.go
*   coder: zemanzeng
*   date : 2021-06-21 20:07:36
*   desc : 脏词检测
*
================================================================*/

package dirtywords

import (
	"sync"

	"github.com/wenruo95/jerry/trie"
)

type DirtyWords struct {
	tree    *trie.Tree
	lock    sync.RWMutex
	ignores map[rune]struct{}
}

func NewDirtyWords() *DirtyWords {
	return &DirtyWords{
		tree:    trie.NewTree(),
		ignores: make(map[rune]struct{}),
	}
}

func NewDirtyWordsWithIgnores(ignoreRunes []rune) *DirtyWords {
	dw := &DirtyWords{
		tree:    trie.NewTree(),
		ignores: make(map[rune]struct{}),
	}
	for _, ig := range ignoreRunes {
		dw.ignores[ig] = struct{}{}
	}
	return dw
}

// DebugLogFunc debug日志输出
type DebugLogFunc func(format string, i ...interface{})

func (dw *DirtyWords) ReBuild(list []string, logFunc DebugLogFunc) {
	tree := trie.NewTree()
	for index, words := range list {
		if succ, list := tree.Insert([]rune(words)); len(list) > 0 && logFunc != nil {
			if succ {
				logFunc("use words[%v]:%v replace %+v", index, words, list)
				continue
			}
			logFunc("insert words[%v]:%v failed because exist %+v", index, words, list)
		}
	}

	dw.lock.Lock()
	dw.tree = tree
	dw.lock.Unlock()
}

func (dw *DirtyWords) Len() int {
	dw.lock.RLock()
	defer dw.lock.RUnlock()
	return dw.tree.Len()
}

func (dw *DirtyWords) List() []string {
	dw.lock.RLock()
	defer dw.lock.RUnlock()
	return dw.tree.Root.List()
}

func (dw *DirtyWords) Match(s string, notRepeat bool) []string {
	dw.lock.RLock()
	matched := dw.tree.Match([]rune(s), dw.ignores)
	dw.lock.RUnlock()

	if !notRepeat || len(matched) == 0 {
		return matched
	}

	// 去重
	m := make(map[string]struct{})
	list := make([]string, 0)
	for _, word := range matched {
		if _, exist := m[word]; exist {
			continue
		}
		m[word] = struct{}{}

		list = append(list, word)
	}
	return list
}

// FilterFunc 过滤函数 默认SingleCustomizeChar
type FilterFunc func(s string, match string) string

// Filter 过滤匹配项
func (dw *DirtyWords) Filter(s string, fn FilterFunc) (string, []string) {
	words := []rune(s)
	indexs := dw.tree.Indexs(words, dw.ignores)
	if len(indexs) == 0 {
		return s, nil
	}

	if fn == nil {
		fn = dw.starReplace
	}

	pos := 0
	match := make([]string, 0, len(indexs))
	result := make([]rune, 0, len(words))
	for index := 0; index < len(words); index++ {
		if pos == len(indexs) {
			result = append(result, []rune(words[index:])...)
			break
		}

		from, length := indexs[pos].Index, indexs[pos].Length
		// 如果没有匹配
		if from != index {
			result = append(result, words[index])
			continue
		}

		// 如果匹配了
		pos = pos + 1
		index = from + length - 1
		match = append(match, string(words[from:from+length]))

		repl := fn(s, string(words[from:from+length]))
		result = append(result, []rune(repl)...)
	}

	return string(result), match
}

func (dw *DirtyWords) starReplace(s string, match string) string {
	return "*"
}

func (dw *DirtyWords) printLine(logFunc DebugLogFunc, format string, i ...interface{}) {
	if logFunc == nil {
		return
	}
	logFunc(format, i...)
}
