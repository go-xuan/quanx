package tirex

import (
	"strings"
)

var Tire *TrieTree

func Init() {
	if Tire == nil {
		Tire = NewTrieTree()
	}
}

// Trie树
type TrieTree struct {
	root      *TrieNode      // 前缀树
	replace   string         // 替换符
	shieldMap map[string]int // 铭感词以及屏蔽方式
}

// Trie树节点
type TrieNode struct {
	children map[string]*TrieNode // 子节点
	end      bool                 // 单词词尾
	terminal bool                 // 树枝末端
}

// 创建前缀树
func NewTrieTree() *TrieTree {
	return &TrieTree{
		root:      NewTrieNode(),
		shieldMap: make(map[string]int),
		replace:   "*",
	}
}

// 创建前缀树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		end:      false,
		terminal: true,
	}
}

// 添加敏感词
func (t *TrieTree) AddWords(words map[string]int) {
	if len(words) > 0 {
		for k, v := range words {
			t.AddWord(k, v)
		}
	}
}

// 添加敏感词
func (t *TrieTree) AddWord(word string, shield int) {
	t.shieldMap[word] = shield
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

// 更新屏蔽方式
func (t *TrieTree) UpdateShieldMethod(word string, shield int) {
	t.shieldMap[word] = shield
}

// 过滤敏感词
func (t *TrieTree) Filter(text string) string {
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

// 判断时候含有敏感词
func (t *TrieTree) HasSensitive(text string) (has bool, shield int) {
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
			var max = t.shieldMap[sb.String()]
			if max > shield {
				shield = max
			}
			start, node = i+1, t.root
		}
	}
	return
}
