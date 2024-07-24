package treex

import (
	"strings"
)

var tireTree *TrieTree

func Trie() *TrieTree {
	if tireTree == nil {
		tireTree = NewTrieTree()
	}
	return tireTree
}

// TrieTree Trie树
type TrieTree struct {
	root    *TrieNode      // 前缀树
	replace string         // 替换符
	level   map[string]int // 敏感度
}

// TrieNode Trie树节点
type TrieNode struct {
	children map[string]*TrieNode // 子节点
	end      bool                 // 单词词尾
	terminal bool                 // 树枝末端
}

// NewTrieTree 创建前缀树
func NewTrieTree() *TrieTree {
	return &TrieTree{
		root:    NewTrieNode(),
		level:   make(map[string]int),
		replace: "*",
	}
}

// NewTrieNode 创建前缀树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		end:      false,
		terminal: true,
	}
}

// AddWordMap 添加敏感词
func (t *TrieTree) AddWordMap(words map[string]int) {
	if len(words) > 0 {
		for word, level := range words {
			t.AddWord(word, level)
		}
	}
}

// AddWords 添加敏感词
func (t *TrieTree) AddWords(words []string) {
	if len(words) > 0 {
		for _, word := range words {
			t.AddWord(word)
		}
	}
}

// AddWord 添加敏感词
func (t *TrieTree) AddWord(word string, level ...int) {
	if len(level) > 0 {
		t.level[word] = level[0]
	}
	texts := strings.Split(word, "")
	node := t.root
	for _, item := range texts {
		if _, ok := node.children[item]; !ok {
			node.children[item] = NewTrieNode()
		}
		node.terminal = false
		node = node.children[item]
	}
	node.end = true
}

// UpdateWordLevel 更新敏感度
func (t *TrieTree) UpdateWordLevel(word string, level int) {
	t.level[word] = level
}

// Desensitize 敏感词脱敏
func (t *TrieTree) Desensitize(text string) string {
	texts := strings.Split(text, "")
	var sb strings.Builder
	var node, start = t.root, 0
	var temp *TrieNode
	for i, item := range texts {
		if child, ok := node.children[item]; ok {
			node = child
		} else if temp != nil {
			if child, ok = temp.children[item]; ok {
				node = child
				temp = nil
			} else {
				sb.WriteString(strings.Join(texts[start:i+1], ""))
				start, node = i+1, t.root
			}
		} else {
			sb.WriteString(strings.Join(texts[start:i+1], ""))
			start, node = i+1, t.root
		}
		if node.end {
			if !node.terminal {
				temp = node
			}
			for j := start; j <= i; j++ {
				sb.WriteString(t.replace)
			}
			start, node = i+1, t.root
		}
	}
	sb.WriteString(strings.Join(texts[start:], ""))
	return sb.String()
}

// HasSensitive 判断时候含有敏感词
func (t *TrieTree) HasSensitive(text string) (has bool, level int) {
	texts := strings.Split(text, "")
	var node, start = t.root, 0
	var sb strings.Builder
	var temp *TrieNode
	for i, item := range texts {
		if child, ok := node.children[item]; ok {
			node = child
		} else if temp != nil {
			if child, ok = temp.children[item]; ok {
				node = child
				temp = nil
			} else {
				sb.Reset()
				start, node = i+1, t.root
			}
		} else {
			sb.Reset()
			start, node = i+1, t.root
		}
		if node.end {
			if !node.terminal {
				temp = node
			}
			has = true
			sb.WriteString(strings.Join(texts[start:i+1], ""))
			if max, ok := t.level[sb.String()]; ok && max > level {
				level = max
			}
			start, node = i+1, t.root
		}
	}
	return
}
