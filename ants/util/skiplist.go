package util

import (
	"math/rand"
)

type IntKey struct {
	Key int
}

func NewIntKey(key int) *IntKey {
	return &IntKey{key}
}

func (this *IntKey) GetNum() int {
	return this.Key
}
func (this *IntKey) Compare(key Key) int8 {
	return (int8)(this.Key - key.GetNum())
}

/*
* skip list
**/

type Key interface {
	Compare(key Key) int8
	GetNum() int
}
type Node struct {
	Key   Key
	Value interface{}
	Next  *Node
}

func NewNode(Key Key, Value interface{}) *Node {
	return &Node{Key, Value, nil}
}

type Index struct {
	Node  *Node
	Right *Index
	Down  *Index
	Level int
}

func NewIndex(Node *Node, Down *Index, Level int) *Index {
	return &Index{Node, nil, Down, Level}
}

var BasicLevel int = 0
var SupporseRange int = 5

var randomSeed int = rand.Intn(100) | 0x0100

type SkipList struct {
	head   *Index
	length int
}

func NewSkipList() *SkipList {
	return &SkipList{nil, 0}
}

func (this *SkipList) Get(key Key) *Node {
	index := this.head
	for index != nil {
		compare := index.Node.Key.Compare(key)
		if compare == 0 {
			return index.Node
		} else if compare > 0 {
			return nil
		} else {
			now := index
			for now.Right == nil {
				now = now.Down
				if now == nil {
					return nil
				}
			}
			index = now.Right
		}
	}
	return nil
}

// head is nil
// key < head
// key>=head
func (this *SkipList) Insert(key Key, Value interface{}) *Node {
	if this.head == nil {
		level := this.randomLevel()
		node := NewNode(key, Value)
		var index *Index = nil
		for i := 0; i <= level; i++ {
			index = NewIndex(node, index, i)
		}
		this.head = index
		return node
	} else {
		if this.head.Node.Key.Compare(key) > 0 {
			level := this.randomLevel()
			if level < this.head.Level {
				level = this.head.Level
			}
			node := NewNode(key, Value)
			node.Next = this.head.Node
			var index *Index = nil
			for i := 0; i <= level; i++ {
				index = NewIndex(node, index, i)
			}
			right := this.head
			this.head = index
			for index != nil {
				if index.Level > right.Level {
					index = index.Down
					continue
				}
				index.Right = right
				right = right.Down
				index = index.Down
			}
			return node
		} else {
			oldNode := this.Get(key)
			if oldNode != nil {
				oldNode.Value = Value
				return oldNode
			}
			level := this.randomLevel()
			if level > this.head.Level {
				newHead := this.head
				for i := this.head.Level + 1; i <= level; i++ {
					newHead = NewIndex(newHead.Node, newHead, i)
				}
				this.head = newHead
			}
			index := this.head
			newNode := NewNode(key, Value)
			newIndexArray := make([]*Index, level+1, level+1)
			var newIndex *Index = nil
			for i := 0; i <= level; i++ {
				newIndex = NewIndex(newNode, newIndex, i)
				newIndexArray = append(newIndexArray, newIndex)
			}
			for index != nil {
				if index.Level > level {
					index = index.Down
					continue
				}
				if index.Right == nil {
					index.Right = newIndexArray[index.Level]
					index = index.Down
				} else {
					compare := index.Right.Node.Key.Compare(key)
					if compare < 0 {
						index = index.Right
					} else {
						newIndex = newIndexArray[index.Level]
						newIndex.Right = index.Right
						index.Right = newIndex
						index = index.Down
					}
				}
			}
			if newIndexArray[0].Right != nil {
				newNode.Next = newIndexArray[0].Right.Node
			}
			return newNode
		}
	}
}

// if not exist
// if head
// other
func (this *SkipList) Delete(key Key) *Node {
	oldNode := this.Get(key)
	if oldNode == nil {
		return nil
	}
	if this.head.Node.Key.Compare(key) == 0 {
		nextNode := this.head.Node.Next
		if nextNode == nil {
			oldNode := this.head.Node
			this.head = nil
			return oldNode
		}
		index := this.head
		for index != nil {
			if index.Node.Key.Compare(nextNode.Key) == 0 {
				break
			}
			index = index.Down
		}
		startLevel := index.Level + 1
		fixLength := this.head.Level - index.Level
		newIndexArray := make([]*Index, fixLength, fixLength)
		newIndex := index.Right
		for i := startLevel; i <= this.head.Level; i++ {
			newIndex := NewIndex(newIndex.Node, newIndex, i)
			newIndexArray = append(newIndexArray, newIndex)
		}
		replaceOldIndex := this.head
		for i := 0; i < fixLength; i++ {
			newIndex := newIndexArray[replaceOldIndex.Level-startLevel]
			newIndex.Right = replaceOldIndex.Right
			replaceOldIndex = replaceOldIndex.Down
		}
		this.head = newIndexArray[len(newIndexArray)-1]
	} else {
		last := this.head
		index := last.Right
		for last != nil {
			if index == nil {
				last = last.Down
				index = last.Right
			} else {
				compare := index.Node.Key.Compare(key)
				if compare == 0 {
					last.Right = index.Right
					last = last.Down
					index = last.Right
				} else if compare < 0 {
					last = index
					index = index.Right
				} else {
					last = last.Down
					index = last.Right
				}
			}
		}
	}
	return oldNode
}

func (this *SkipList) randomLevel() int {
	x := randomSeed
	x ^= x << 13
	x ^= x >> 16
	x ^= x << 5
	randomSeed = x
	if (x & 0x80000001) != 0 {
		return 0
	}
	level := 1
	x = x >> 1
	for (x & 1) != 0 {
		level += 1
		x = x >> 1
	}
	return level
}
