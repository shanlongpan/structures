package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//最大层级32层，层级参数0.25，减少高层级数据数量
const (
	MAX_LEVEL   = 32
	PROBABILITY = 0.25
)
//node具体节点
type Node struct {
	index     uint64
	value     []int64
	nextNodes []*Node
}

// 创建新的节点
func newNode(index uint64, value []int64, level int) *Node {
	return &Node{
		index:     index,
		value:     value,
		nextNodes: make([]*Node, level, level),
	}
}

// 返回index
func (n *Node) Index() uint64 {
	return n.index
}

// 返回value
func (n *Node) Value() interface{} {
	return n.value
}

//跳跃表结构体
type skipList struct {
	level  int
	length int32
	head   *Node
	tail   *Node
	mutex  sync.RWMutex
}


//创建跳跃表，合适的层级经验值   L(N) = log(1/PROBABILITY)(N).  PROBABILITY=0.25 ，100万的数据， L(N) ≈ 12

func newSkipList(level int) *skipList {
	if level>MAX_LEVEL{
		level=MAX_LEVEL
	}
	head := newNode(0, nil, level)
	var tail *Node
	for i := 0; i < len(head.nextNodes); i++ {
		head.nextNodes[i] = tail
	}

	return &skipList{
		level:  level,
		length: 0,
		head:   head,
		tail:   tail,
	}
}


//插入的时候使用，返回需要更新的 previous nodes ， The second return value represents the value with given index or the closet value whose index is larger than given index.
func (s *skipList) searchWithPreviousNodes(index uint64) ([]*Node, *Node) {
	// Store all previous value whose index is less than index and whose next value's index is larger than index.
	previousNodes := make([]*Node, s.level)

	// fmt.Printf("start doSearch:%v\n", index)
	currentNode := s.head

	// Iterate from top level to bottom level.
	for l := s.level - 1; l >= 0; l-- {
		// Iterate value util value's index is >= given index.
		// The max iterate count is skip list's length. So the worst O(n) is N.
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}

		// When next value's index is >= given index, add current value whose index < given index.
		previousNodes[l] = currentNode
	}

	// Avoid point to tail which will occur panic in Insert and Delete function.
	// When the next value is tail.
	// The index is larger than the maximum index in the skip list or skip list's length is 0. Don't point to tail.
	// When the next value isn't tail.
	// Next value's index must >= given index. Point to it.
	if currentNode.nextNodes[0] != s.tail {
		currentNode = currentNode.nextNodes[0]
	}

	return previousNodes, currentNode
}

//查询的时候使用，仅返回index对应的node index 不存在，返回nil

func (s *skipList) searchWithoutPreviousNodes(index uint64) *Node {
	currentNode := s.head

	// Read lock and unlock.
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Iterate from top level to bottom level.
	for l := s.level - 1; l >= 0; l-- {
		// Iterate value util value's index is >= given index.
		// The max iterate count is skip list's length. So the worst O(n) is N.
		for currentNode.nextNodes[l] != s.tail && currentNode.nextNodes[l].index < index {
			currentNode = currentNode.nextNodes[l]
		}
	}

	currentNode = currentNode.nextNodes[0]
	if currentNode == s.tail || currentNode.index > index {
		return nil
	} else if currentNode.index == index {
		return currentNode
	} else {
		return nil
	}
}

//插入或者更新跳跃表，不存在就插入，存在就覆盖更新

func (s *skipList) insert(index uint64, value []int64) {
	// Write lock and unlock.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	previousNodes, currentNode := s.searchWithPreviousNodes(index)

	if currentNode != s.head && currentNode.index == index {
		currentNode.value = value
		return
	}

	// Make a new value.
	newNode := newNode(index, value, s.randomLevel())

	// Adjust pointer. Similar to update linked list.
	for i := len(newNode.nextNodes) - 1; i >= 0; i-- {
		// Firstly, new value point to next value.
		newNode.nextNodes[i] = previousNodes[i].nextNodes[i]

		// Secondly, previous nodes point to new value.
		previousNodes[i].nextNodes[i] = newNode

		// Finally, in order to release the slice, point to nil. 降低内存
		previousNodes[i] = nil
	}

	atomic.AddInt32(&s.length, 1)

	// 降低内存
	for i := len(newNode.nextNodes); i < len(previousNodes); i++ {
		previousNodes[i] = nil
	}
}

//删除index，存在就删除，更新跳跃表长度，不存在不处理
func (s *skipList) delete(index uint64) {
	// Write lock and unlock.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	previousNodes, currentNode := s.searchWithPreviousNodes(index)

	// If skip list length is 0 or could not find value with the given index.
	if currentNode != s.head && currentNode.index == index {
		// Adjust pointer. Similar to update linked list.
		for i := 0; i < len(currentNode.nextNodes); i++ {
			previousNodes[i].nextNodes[i] = currentNode.nextNodes[i]
			currentNode.nextNodes[i] = nil
			// 降低内存
			previousNodes[i] = nil
		}

		atomic.AddInt32(&s.length, -1)
	}
    // 降低内存
	for i := len(currentNode.nextNodes); i < len(previousNodes); i++ {
		previousNodes[i] = nil
	}
}

//创建一个跳跃表快照，返回node的切片

func (s *skipList) snapshot() []*Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*Node, s.length)
	i := 0

	currentNode := s.head.nextNodes[0]
	for currentNode != s.tail {
		node := &Node{
			index:     currentNode.index,
			value:     currentNode.value,
			nextNodes: nil,
		}

		result[i] = node
		currentNode = currentNode.nextNodes[0]
		i++
	}

	return result
}

// 返回跳跃表长度
func (s *skipList) getLength() int32 {
	return atomic.LoadInt32(&s.length)
}

//生成一个随机的层级

func (s *skipList) randomLevel() int {
	level := 1
	for rand.Float64() < PROBABILITY && level < s.level {
		level++
	}

	return level
}

//打印跳跃表
func (s *skipList) print() {

	for i := s.level - 1; i >= 0; i-- {
		current := s.head
		for current.nextNodes[i] != nil {
			fmt.Printf("%d ", current.nextNodes[i].index)
			current = current.nextNodes[i]
		}
		fmt.Printf("***************** Level %d \n", i+1)
	}
}

func main(){
	skipNew:=newSkipList(12)
	beginTimer:=time.Now()
	timer:=beginTimer.Unix()
	var fistTimer int64
	for i:=0;i<15000000;i++{
		rand.Seed(time.Now().UnixNano())
		randNumber:=rand.Int63n(8640000) // 十天内
		fistTimer=timer-randNumber
		skipNew.insert(uint64(i),[]int64{fistTimer,fistTimer})
	}

	fmt.Println(time.Since(beginTimer))
	var s =make([]int64,0,15000)
	t:=time.Now()

	week := t.AddDate(0, 0, -7).Unix()
	todaySince := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	for i:=0;i<15000000;i++{
		if i%1000==0{
			nod:=skipNew.searchWithoutPreviousNodes(uint64(i))

			if nod.value[1]>todaySince{
				s=append(s,int64(i))
			}else if nod.value[0]>week{
				s=append(s,int64(i))
			}else {
				continue
			}
			nod.value=append(nod.value,timer)
			nod.value=nod.value[1:]
		}
	}
	fmt.Println(len(s))
	fmt.Println(time.Since(t))
}