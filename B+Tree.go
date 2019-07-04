package main

import "sort"

const (
	// 叶子节点最大元素存储数目
	MaxKV = 255
	// 中间节点最大元素存储数目
	MaxKC = 511
)

// Value 定义
type kc struct {
	key   int
	child node
}

// 预留一个空槽， 数组
type kcs [MaxKC + 1]kc

func (a *kcs) Len() int { return len(a) }

func (a *kcs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a *kcs) Less(i, j int) bool {
	// 处理数组中预留的空槽，value key 初始值是 0；应对 Search
	if a[i].key == 0 {
		return false
	}
	// 处理数组中预留的空槽，value key 初始值是 0；应对 Sort
	if a[j].key == 0 {
		return true
	}

	return a[i].key < a[j].key
}

// 中间节点数据结构定义
type interiorNode struct {
	kcs   kcs           // 存储元素
	count int           // 实际存储元素数目
	p     *interiorNode // 父亲节点
}

func newInteriorNode(p *interiorNode, largestChild node) *interiorNode {
	i := &interiorNode{
		p:     p,
		count: 1,
	}

	if largestChild != nil {
		i.kcs[0].child = largestChild
	}
	return i
}

// 从该中间节点找到 key 元素应该存储的位置
func (in *interiorNode) find(key int) (int, bool) {
	// 定义查询方法，这里只需要 ">"
	c := func(i int) bool { return in.kcs[i].key > key }
	// 查询
	i := sort.Search(in.count-1, c)

	return i, true
}

// 插入中间节点
func (in *interiorNode) insert(key int, child node) (int, *interiorNode, bool) {
	// 确定 key 在该中间节点应该存储的位置
	i, _ := in.find(key)
	// 中间节点没有达到数量限制
	if !in.full() {
		// 将元素插入中间节点
		copy(in.kcs[i+1:], in.kcs[i:in.count])
		// 设置子节点分裂后产生的元素 为当前位置 i 的key
		in.kcs[i].key = key
		// 设置子节点以及子节点设置父亲节点
		in.kcs[i].child = child
		child.setParent(in)
		// 元素计数加一
		in.count++
		return 0, nil, false
	}

	// 达到数量限制，则在最右侧保留的空槽追加该节点
	in.kcs[MaxKC].key = key
	in.kcs[MaxKC].child = child
	// 子节点设置父亲节点
	child.setParent(in)
	// 中间节点分裂
	next, midKey := in.split()

	return midKey, next, true
}

func (in *interiorNode) split() (*interiorNode, int) {
	// 节点排序，把新插入的节点防到正确的位置
	sort.Sort(&in.kcs)

	// 获取中间元素的位置，并设置 Value 的 子节点和 key
	midIndex := MaxKC / 2
	midChild := in.kcs[midIndex].child
	midKey := in.kcs[midIndex].key

	// 创建一个新没有父亲节点(第一个 nil)的中间节点
	next := newInteriorNode(nil, nil)
	// 将中间元素的右侧数组拷贝到新的分裂节点
	copy(next.kcs[0:], in.kcs[midIndex+1:])
	// 初始化原始节点的右半部分
	in.initArray(midIndex + 1)
	// 设置分裂节点的 Count
	next.count = MaxKC - midIndex
	// 更新分裂节点中所有元素子节点的父亲节点
	for i := 0; i < next.count; i++ {
		next.kcs[i].child.setParent(next)
	}

	// 更新原始节点的参数，将中间元素放进原始节点
	in.count = midIndex + 1
	in.kcs[in.count-1].key = 0
	in.kcs[in.count-1].child = midChild
	midChild.setParent(in)
	// 返回分裂后产生的中间节点和中间元素的 key，供父亲节点插入
	return next, midKey
}

// 判断是否达到中间节点的最大元素数目限制 MaxKC
func (in *interiorNode) full() bool { return in.count == MaxKC }

// 返回中间节点的父亲节点
func (in *interiorNode) parent() *interiorNode { return in.p }

// 设置中间节点的父亲节点
func (in *interiorNode) setParent(p *interiorNode) { in.p = p }

// 获取中间结点存储的元素数目
func (in *interiorNode) countNum() int { return in.count }

// 初始化数组从 num 起的元素为空结构
func (in *interiorNode) initArray(num int) {
	for i := num; i < len(in.kcs); i++ {
		in.kcs[i] = kc{}
	}
}

// 接口设计
type node interface {
	// 确定元素在节点中的位置
	find(key int) (int, bool)
	// 获取父亲节点
	parent() *interiorNode
	// 设置父亲节点
	setParent(*interiorNode)
	// 是否达到最大数目限制
	full() bool
	// 元素数目统计
	countNum() int
}

// Value 定义
type kv struct {
	key   int
	value string
}

// 叶子节点的存储数组，[M]value
type kvs [MaxKV]kv

func (a *kvs) Len() int           { return len(a) }
func (a *kvs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a *kvs) Less(i, j int) bool { return a[i].key < a[j].key }

// 叶子节点的数据结构
type leafNode struct {
	kvs   kvs           // 存储元素
	count int           // 实际存储元素数目
	next  *leafNode     // 右边第一个叶子节点（右指针）
	p     *interiorNode // 父亲节点，中间节点
}

// 创建新的叶子节点
func newLeafNode(p *interiorNode) *leafNode {
	return &leafNode{
		p: p,
	}
}

func (l *leafNode) find(key int) (int, bool) {
	c := func(i int) bool {
		return l.kvs[i].key >= key
	}
	// 查询
	i := sort.Search(l.count, c)
	// 判断是否 key 已经存在
	if i < l.count && l.kvs[i].key == key {
		return i, true
	}

	return i, false
}

// insert
func (l *leafNode) insert(key int, value string) (int, *leafNode, bool) {
	i, ok := l.find(key)

	if ok {
		l.kvs[i].value = value
		return 0, nil, false
	}
	// 判断叶子节点是否已经填满
	if !l.full() {
		copy(l.kvs[i+1:], l.kvs[i:l.count])
		l.kvs[i].key = key
		l.kvs[i].value = value
		l.count++
		return 0, nil, false
	}
	// 获取分裂出新的节点
	next := l.split()
	// 判断插入位置：新的节点或者旧节点
	if key < next.kvs[0].key {
		l.insert(key, value)
	} else {
		next.insert(key, value)
	}
	// 返回分裂节点的第一个key
	return next.kvs[0].key, next, true
}

// 叶子节点分裂过程
func (l *leafNode) split() *leafNode {
	// 申请一个右节点
	next := newLeafNode(nil)
	// 将原始节点的右半部分复制过去
	copy(next.kvs[0:], l.kvs[l.count/2+1:])
	// 初始化原始节点的右半部分
	l.initArray(l.count/2 + 1)
	// 设置右节点的参数
	next.count = MaxKV - l.count/2 - 1
	next.next = l.next
	// 重新设置原始节点的参数
	l.count = l.count/2 + 1
	l.next = next
	// // 设置右节点在 map 中的 key
	// next.setKey(next.kvs[0].key)
	// 返回右节点指针
	return next
}

// 判断是否达到 key 数目上限 MaxKV
func (l *leafNode) full() bool { return l.count == MaxKV }

// 返回父节点，中间节点
func (l *leafNode) parent() *interiorNode { return l.p }

// 设置父中间节点
func (l *leafNode) setParent(p *interiorNode) { l.p = p }

// 获取叶子结点存储的元素数目
func (l *leafNode) countNum() int { return l.count }

// 初始化数组从 num 起的元素为空结构
func (l *leafNode) initArray(num int) {
	for i := num; i < len(l.kvs); i++ {
		l.kvs[i] = kv{}
	}
}

type BTree map[int]node

// 创建自由一个父亲节点和叶子节点的 B+ 树
func NewBTree() *BTree {
	bt := BTree{}
	leaf := newLeafNode(nil)
	r := newInteriorNode(nil, leaf)
	leaf.p = r
	bt[-1] = r
	bt[0] = leaf
	return &bt
}

// 返回 B+ Tree 存储的元素数目
func (bt *BTree) Count() int {
	count := 0
	leaf := (*bt)[0].(*leafNode)
	for {
		count += leaf.countNum()
		if leaf.next == nil {
			break
		}
		leaf = leaf.next
	}
	return count
}

// 返回根结点
func (bt *BTree) Root() node {
	return (*bt)[-1]
}

// 返回 第一个叶子结点
func (bt *BTree) First() node {
	return (*bt)[0]
}

// 返回由叶子结点指针构成的数组，从最左侧开始依次追加
func (bt *BTree) Values() []*leafNode {
	nodes := make([]*leafNode, 0)
	leaf := (*bt)[0].(*leafNode)
	for {
		nodes = append(nodes, leaf)
		if leaf.next == nil {
			break
		}
		leaf = leaf.next
	}

	return nodes
}

// 在 B+ 树中，插入 key-value
func (bt *BTree) Insert(key int, value string) {
	// 确定插入的位置，是一个叶子节点
	_, oldIndex, leaf := search((*bt)[-1], key)
	// 获取叶子节点的父亲节点，中间节点
	p := leaf.parent()
	// 插入叶子节点，返回是否分裂
	mid, nextLeaf, bump := leaf.insert(key, value)
	// 未分裂，则直接返回
	if !bump {
		return
	}

	// 填充分裂的节点到 map
	(*bt)[mid] = nextLeaf

	// 分裂则继续插入中间节点
	var midNode node
	midNode = leaf
	// 设置父亲节点指向分裂出的子（叶子）节点
	p.kcs[oldIndex].child = leaf.next
	// 新分裂出的节点设置该中间节点为父亲节点
	leaf.next.setParent(p)
	// 赋值，获取该中间节点和其父节点
	interior, interiorP := p, p.parent()
	// 迭代向上判断父亲节点是否需要分裂
	for {
		var oldIndex int
		var newNode *interiorNode
		// 判断是否到达根节点
		isRoot := interiorP == nil
		// 未到达根节点，在父亲节点中查询元素的位置
		if !isRoot {
			oldIndex, _ = interiorP.find(key)
		}
		// 将叶子节点分裂后产生的中间元素同时传给父亲中间节点，并传入分裂后的原始叶子节点
		// 同时返回分裂后产生的中间节点和中间元素的 key
		mid, newNode, bump = interior.insert(mid, midNode)
		// 未分裂，直接返回
		if !bump {
			return
		}
		// 填充 map
		(*bt)[newNode.kcs[0].key] = newNode

		if !isRoot {
			// 未到达根节点，将元素插入父亲节点，基本过程同上
			interiorP.kcs[oldIndex].child = newNode
			newNode.setParent(interiorP)

			midNode = interior
		} else {
			// 更新 map 中的 root 节点
			(*bt)[interior.kcs[0].key] = (*bt)[-1]
			// 到达根节点，根节点上移，并插入原始中间节点
			(*bt)[-1] = newInteriorNode(nil, newNode)
			node := (*bt)[-1].(*interiorNode)
			node.insert(mid, interior)
			(*bt)[-1] = node
			newNode.setParent(node)

			return
		}
		// 赋值，获取该中间节点的父亲节点和其父亲的父节点
		interior, interiorP = interiorP, interior.parent()
	}
}

// 搜索： 找到，则返回 value ，否则返回空value
func (bt *BTree) Search(key int) (string, bool) {
	kv, _, _ := search((*bt)[-1], key)
	if kv == nil {
		return "", false
	}
	return kv.value, true
}

// 具体搜索过程
func search(n node, key int) (*kv, int, *leafNode) {
	curr := n
	oldIndex := -1

	for {
		switch t := curr.(type) {
		// 叶子节点，返回命中节点或者可插入位置
		case *leafNode:
			i, ok := t.find(key)
			if !ok {
				return nil, oldIndex, t
			}
			return &t.kvs[i], oldIndex, t
		// 中间节点迭代查询
		case *interiorNode:
			i, _ := t.find(key)
			curr = t.kcs[i].child
			oldIndex = i
		default:
			panic("")
		}
	}
}
