/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : tree.go
*   coder: zemanzeng
*   date : 2021-06-17 17:32:39
*   desc : trie树
*
================================================================*/

package trie

import (
	"encoding/json"
	"sync/atomic"
)

// Node 树结点
type Node struct {
	Words  []rune         // 叶子结点专属
	Childs map[rune]*Node // 非叶子结点
}

func (node *Node) Leaf() bool {
	return node != nil && len(node.Words) > 0
}

func (node *Node) List() []string {
	if node.Leaf() {
		return []string{string(node.Words)}
	}
	list := make([]string, 0)
	for _, child := range node.Childs {
		list = append(list, child.List()...)
	}
	return list
}

// Tree 前缀树
type Tree struct {
	Root      *Node
	NodeCount int32
}

func NewTree() *Tree {
	return &Tree{
		Root: &Node{
			Childs: make(map[rune]*Node),
		},
	}
}

func (tree *Tree) Len() int {
	return int(atomic.LoadInt32(&tree.NodeCount))
}

func (tree *Tree) Insert(words []rune) (bool, []string) {
	if len(words) == 0 {
		return false, nil
	}

	node := tree.Root
	for index, word := range words {
		child, exist := node.Childs[word]
		if !exist {
			// 不存在的添加
			child = &Node{}
			node.Childs[word] = child

			if index == len(words)-1 {
				child.Words = words
				atomic.AddInt32(&tree.NodeCount, 1)
				return true, nil
			}

			child.Childs = make(map[rune]*Node)
			node = child
			continue
		}

		// 存在的首先判断是否叶子 是叶子直接返回叶子树的数据
		if child.Leaf() {
			return false, child.List()
		}

		// 存在 & 非叶子结点 & 当前是最后一个
		if index == len(words)-1 {
			// 最短原则 abc 被 ab 替换
			replaceList := make([]string, 0)
			for _, tmpNode := range child.Childs {
				replaceList = append(replaceList, tmpNode.List()...)
			}
			atomic.AddInt32(&tree.NodeCount, int32(1-len(replaceList)))

			child.Words = []rune(words)
			child.Childs = nil
			return true, replaceList
		}

		node = child
	}

	return false, nil
}

func (tree *Tree) Delete(words []rune) bool {
	if len(words) == 0 {
		return false
	}

	if !tree.Exist(words) {
		return false
	}
	atomic.AddInt32(&tree.NodeCount, -1)

	for cnt := 0; cnt < len(words); cnt++ {
		if !tree.delete(words[0 : len(words)-cnt]) {
			break
		}
	}
	return true
}

// 删除叶子结点 或者 非叶子但是没有子树的结点
func (tree *Tree) delete(words []rune) bool {
	if len(words) == 0 {
		return false
	}

	node := tree.Root
	var parrent *Node
	for index, word := range words {
		child, exist := node.Childs[word]
		if !exist {
			return false
		}
		parrent = node
		node = child

		if node.Leaf() {
			if (index == len(words)-1) && (node.Words[index] == words[index]) {
				delete(parrent.Childs, word)
				return true
			}
			return false
		}
		if len(node.Childs) == 0 && index == len(words)-1 {
			delete(parrent.Childs, word)
			return true
		}
	}
	return false
}

func (tree *Tree) Exist(words []rune) bool {
	if len(words) == 0 {
		return false
	}

	node := tree.Root
	for index, word := range words {
		child, exist := node.Childs[word]
		if !exist {
			return false
		}
		node = child

		if node.Leaf() {
			// 最后一个 & 长度相等 & 最后一个字符相等
			return (index == len(words)-1) && (node.Words[index] == words[index])
		}
	}
	return false
}

func (tree *Tree) String() string {
	if buff, err := json.Marshal(tree); err == nil {
		return string(buff)
	}
	return ""
}

// Match 返回匹配字符串
func (tree *Tree) Match(words []rune, ignores map[rune]struct{}) []string {
	list := make([]string, 0)
	for index, word := range words {
		node, exist := tree.Root.Childs[word]
		if !exist {
			continue
		}

		// 单个字符匹配
		if node.Leaf() {
			list = append(list, string(node.Words))
			continue
		}

		for cursor := index + 1; cursor < len(words); cursor++ {
			child, exist := node.Childs[words[cursor]]
			if !exist {
				if ignoreFn(ignores, words[cursor]) {
					continue
				}
				break
			}
			// 寻找最短匹配项
			if child.Leaf() {
				list = append(list, string(child.Words))
				index = cursor // 找到后移动到当前匹配项最后一个单词
				break
			}
			node = child
		}
	}

	return list
}

// RuneIndex 字符串匹配index数据
type RuneIndex struct {
	Index  int
	Length int
	Words  string
}

// Indexs 返回匹配项字符串起始下标和匹配字符串长度数组
func (tree *Tree) Indexs(words []rune, ignores map[rune]struct{}) []*RuneIndex {
	indexs := make([]*RuneIndex, 0)
	for index, word := range words {
		node, exist := tree.Root.Childs[word]
		if !exist {
			continue
		}

		// 单个字符匹配
		if node.Leaf() {
			indexs = append(indexs, &RuneIndex{
				Index:  index,
				Length: len(node.Words),
				Words:  string(node.Words),
			})
			continue
		}

		// 字符串匹配
		for cursor := index + 1; cursor < len(words); cursor++ {
			child, exist := node.Childs[words[cursor]]
			if !exist {
				if ignoreFn(ignores, words[cursor]) {
					continue
				}
				break
			}
			// 寻找最短匹配项
			if child.Leaf() {
				indexs = append(indexs, &RuneIndex{
					Index:  index,
					Length: cursor - index + 1,
					Words:  string(child.Words),
				})
				index = cursor // 找到后移动到当前匹配项最后一个单词
				break
			}
			node = child
		}
	}
	return indexs
}

func ignoreFn(ignores map[rune]struct{}, word rune) bool {
	if ignores == nil {
		return false
	}
	if _, exist := ignores[word]; exist {
		return true
	}
	return false
}
