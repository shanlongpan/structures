package btree
//
//import (
//	"fmt"
//)
//
//type BTreeNode struct {
//	n        int          // number of keys
//	t        int          // minimum degree(minimum children)
//	leaf     bool         // true for leaf node
//	keys     []int        // key slice
//	children []*BTreeNode // child node slice
//}
//
//func NewBTreeNode(t int, leaf bool) *BTreeNode {
//	return &BTreeNode{
//		n:        0,
//		t:        t,
//		leaf:     leaf,
//		keys:     make([]int, 2*t-1),
//		children: make([]*BTreeNode, 2*t),
//	}
//}
//
//// 遍历所有以本节点为根的所有结点，打印它们的key
//func (s *BTreeNode) Traverse() {
//	var i int
//	for i = 0; i < s.n; i++ {
//		if !s.leaf {
//			s.children[i].Traverse()
//		}
//		fmt.Printf("%d ", s.keys[i])
//	}
//	if s.leaf {
//		return
//	}
//	// 遍历最后一个子节点
//	s.children[i].Traverse()
//}
//
//// 在以本节点为根的结点中搜索包含指定关键字的结点以及关键字所在位置
//func (s *BTreeNode) Search(k int) (*BTreeNode, int) {
//	var i int
//	// 找到第一个不大于key的关键字所在位置
//	// 结束条件，i = s.n || s.keys[i] >= k
//	for i = 0; i < s.n && s.keys[i] < k; i++ {
//	}
//	if i < s.n && s.keys[i] == k {
//		return s, i
//	}
//	if s.leaf {
//		return nil, 0
//	}
//	// 在孩子中查找
//	// 包含 i == s.n 或 s.keys[i] > k
//	return s.children[i].Search(k)
//}
//func (s *BTreeNode) insertNonFull(k int) {
//	// 找到要插入的位置
//	var i int
//	for i = 0; i < s.n && s.keys[i] < k; i++ {
//	}
//	// k 已经在本节点里存在
//	if i < s.n && s.keys[i] == k {
//		return
//	}
//	// 如果是叶节点则插入
//	if s.leaf {
//		// 将i之后的key右移
//		copy(s.keys[i+1:], s.keys[i:s.n])
//		// k放到位置i
//		s.keys[i] = k
//		s.n++
//		return
//	}
//	var c = s.children[i]
//	// 如果子节点已满，则分裂子节点
//	if c.isFull() {
//		s.splitChild(c, i)
//		// 如果提升上来的key等于k则不插入
//		if s.keys[i] == k {
//			return
//		}
//		// 如果提升上来的新的key小于k则插入到它的右孩子
//		if s.keys[i] < k {
//			c = s.children[i+1]
//		}
//	}
//	c.insertNonFull(k)
//}
//
//// c 为要分裂的子节点
//// i 为父节点保存要提升的子节点key的位置
//// 关键步骤：
//// 	1. 生成新的结点将分裂结点后t个key和child复制过去
////  2. 父节点key和child数组右移
//// 	3. 新节点和上升的key复制到父节点
//func (s *BTreeNode) splitChild(c *BTreeNode, i int) {
//	// 生成新的结点
//	z := NewBTreeNode(s.t, c.leaf)
//	z.n, c.n = c.t-1, c.t-1
//	// 将c的后t-1个key复制到新的结点
//	// c中keys索引: 0..t-2, t-1, t..2t-2
//	for j := c.t; j < 2*c.t-1; j++ {
//		z.keys[j-c.t] = c.keys[j]
//		c.keys[j] = 0
//	}
//	// 将s的keys里i位之后的元素右移
//	for j := s.n - 1; j >= i; j-- {
//		s.keys[j+1] = s.keys[j]
//	}
//	// 将中间的key复制到s的keys列表的位置i
//	s.keys[i] = c.keys[c.t-1]
//	c.keys[c.t-1] = 0
//	// 将c的后t个孩子指针复制到新的结点
//	// c中children索引: 0..t-1, t..2t-1
//	if !c.leaf {
//		for j := c.t; j < 2*c.t; j++ {
//			z.children[j-c.t] = c.children[j]
//			c.children[j] = nil
//		}
//	}
//	// 将s的i+1之后的孩子指针在列表中后移
//	for j := s.n; j >= i+1; j-- {
//		s.children[j+1] = s.children[j]
//	}
//	// 将新的结点指针放到s的孩子列表位置i
//	s.children[i+1] = z
//	s.n++
//}
//
//// 找到第一个大于key的关键字的位置
//func (s *BTreeNode) findKey(k int) (idx int) {
//	for idx = 0; idx < s.n && s.keys[idx] < k; idx++ {
//	}
//	return
//}
//
//// 1. 如果key在节点中，并且该结点是leaf，则直接删除
//// 2. 如果key在节点中，并且该结点是内结点，分如下情况
//// 	2.1 若key所在的位置的左子树是满足最少t个关键字，则将其子树下最接近key的关键字key1复制到key然后在左子树递归删掉key1，否则
//// 	2.2 若key所在的位置的右子树是满足最少t个关键字，则将其子树下最接近key的关键字key1复制到key然后在左子树递归删掉key1，否则
//// 	2.3 若key的左右子树都少于t个关键字，则合并两子树合并并删掉key
//// 3. 若key不在结点中，找到key最有可能在的子结点，并删除，删除之前判断如果子结点的key个数为t-1个
//// 	3.1 如果其左兄弟有t个key则借一个key，如果右兄弟有t个key则从其借一个key
////  3.2 如果左或右兄弟都只有t-1个key，则从合并左或右兄弟
//func (s *BTreeNode) delete(k int) {
//	idx := s.findKey(k)
//	// key 在本结点内
//	if idx < s.n && s.keys[idx] == k {
//		switch {
//		case s.leaf: // 情况1
//			s.deleteFromLeaf(idx)
//		case s.children[idx].n >= s.t: // 情况2.1
//			pre := s.getPre(idx)
//			s.keys[idx] = pre
//			s.children[idx].delete(pre)
//		case s.children[idx+1].n >= s.t: // 情况2.2
//			succ := s.getSucc(idx)
//			s.keys[idx] = succ
//			s.children[idx+1].delete(succ)
//		default: // 情况2.3
//			s.merge(idx)
//			s.children[idx].delete(k)
//		}
//		return
//	}
//	// 未在本节点并是leaf
//	if s.leaf {
//		return
//	}
//	// key在最后一个子节点子树中
//	flag := idx == s.n
//	// 情况3
//	if s.children[idx].n < s.t {
//		s.fill(idx)
//	}
//	// fill时如果有merge操作，关键词个数减少
//	// 如果关键字在最后一个子节点，但其被merge
//	// 则需要在前一个子节点删除
//	if flag && idx > s.n {
//		s.children[idx-1].delete(k)
//	} else {
//		s.children[idx].delete(k)
//	}
//}
//
//// 将第idx的子树从其兄弟中借一个key
//func (s *BTreeNode) fill(idx int) {
//	switch {
//	// 有左子树且其关键字个数>=t
//	case idx != 0 && s.children[idx-1].n >= s.t: // 3.1
//		s.borrowFromPre(idx)
//		// 有右子树且其关键字个数>=t
//	case idx != s.n && s.children[idx+1].n >= s.t: // 3.1
//		s.borrowFromSucc(idx)
//		// idx在最后一个key位置
//	case idx != s.n:
//		s.merge(idx)
//		// idx为key的最后一个
//	default:
//		s.merge(idx - 1)
//	}
//}
//
//// 从左兄弟借一个key
//func (s *BTreeNode) borrowFromPre(idx int) {
//	cur := s.children[idx]
//	left := s.children[idx-1]
//	// 右移cur关键字为借的key留出空间
//	copy(cur.keys[1:], cur.keys[0:cur.n])
//	// s的关键字下移
//	cur.keys[0] = s.keys[idx-1]
//	// 左兄弟的最后一个key上升
//	s.keys[idx-1] = left.keys[left.n-1]
//	// 右移cur的子指针，为借的指针留出空间
//	if !cur.leaf {
//		copy(cur.children[1:], cur.children[0:cur.n+1])
//		// 左兄弟最后一个child移到cur
//		cur.children[0] = left.children[left.n]
//	}
//	left.n--
//	cur.n++
//}
//
//// 从右兄弟借一个key
//func (s *BTreeNode) borrowFromSucc(idx int) {
//	cur := s.children[idx]
//	right := s.children[idx+1]
//	// s的关键字下移到cur的关键字尾部
//	cur.keys[cur.n] = s.keys[idx]
//	// 右兄弟的第一个关键词上移
//	s.keys[idx] = right.keys[0]
//	// 右兄弟的关键字列表左移
//	copy(right.keys[0:], right.keys[1:right.n])
//	// 右兄弟的第一个子指针借到cur的最后一个child
//	if !cur.leaf {
//		cur.children[cur.n+1] = right.children[0]
//		// 有兄弟的子指针列表左移
//		copy(right.children[:], right.children[1:right.n+1])
//	}
//	right.n--
//	cur.n++
//}
//
//// 合并key的左右子树
//func (s *BTreeNode) merge(idx int) {
//	left := s.children[idx]
//	right := s.children[idx+1]
//	// s的关键字下降
//	left.keys[left.t-1] = s.keys[idx]
//	// 复制右子t-1个key
//	for i := 0; i < right.t-1; i++ {
//		left.keys[left.t+i] = right.keys[i]
//	}
//	if !left.leaf {
//		// 复制右子t个child
//		for i := 0; i < right.t; i++ {
//			left.children[left.n+i+1] = right.children[i]
//		}
//	}
//	left.n += right.t
//	// s的key复制
//	for i := idx; i < s.n-1; i++ {
//		s.keys[i] = s.keys[i+1]
//	}
//	// s的child复制
//	for i := idx + 1; i < s.n; i++ {
//		s.children[i] = s.children[i+1]
//	}
//	s.n--
//}
//
//// 从左子树取最大值
//func (s *BTreeNode) getPre(idx int) int {
//	cur := s.children[idx]
//	for !cur.leaf {
//		cur = cur.children[cur.n]
//	}
//	return cur.keys[cur.n-1]
//}
//
//// 从右子树取最小值
//func (s *BTreeNode) getSucc(idx int) int {
//	cur := s.children[idx+1]
//	for !cur.leaf {
//		cur = cur.children[0]
//	}
//	return cur.keys[0]
//}
//func (s *BTreeNode) deleteFromLeaf(idx int) {
//	copy(s.keys[idx:], s.keys[idx+1:s.n])
//	s.n--
//}
//func (s *BTreeNode) isFull() bool {
//	return s.n == 2*s.t-1
//}
//
//type BTree struct {
//	root *BTreeNode
//	t    int
//}
//
//func NewBTree(t int) *BTree {
//	return &BTree{t: t}
//}
//
//// 单程下行方式遍历树插入关键字。
//// 关键方法是"主动分裂", 即，在遍历一个子节点前，如果子节点已满则对其进行分裂。
//// 相反的"被动分裂"则是在要插入的时候遇满才分裂，会出现重复遍历的情况。
//// 比如，从根节点到叶节点都是满的，当到达叶节点要发现其已满需进行分裂，
//// 提升一个关键字到父节点，发现父节点也已满需要分裂父节点，重复下去一直到根节点。
//// 这样就出现了从根节点到叶节点的重复遍历。而"主动分裂"则不出有出现这种情况，因为
//// 在分裂一个子节点时候父节点是已经有足够空间容纳要提升的新的key了
////
//// B数的增高依赖的是root结点分裂
//// B数新增的关键字都增加到了叶节点上
//func (t *BTree) Insert(k int) {
//	if t.root == nil {
//		t.root = NewBTreeNode(t.t, true)
//	}
//	var r = t.root
//	// root已满则分裂root
//	if r.isFull() {
//		// 为新的root生成结点
//		s := NewBTreeNode(t.t, false)
//		// 老的root称为新结点的孩子
//		s.children[0] = r
//		// 分裂老的root，并将一个key提升到新的root
//		s.splitChild(r, 0)
//		// 新的root有两个孩子，决定将key插入到哪个孩子
//		var i int
//		if k > s.keys[0] {
//			i++
//		}
//		s.children[i].insertNonFull(k)
//		// 更新root
//		t.root = s
//		return
//	}
//	t.root.insertNonFull(k)
//}
//func (t *BTree) Delete(k int) {
//	if t.root == nil {
//		return
//	}
//	t.root.delete(k)
//}
//func (t *BTree) Traverse() {
//	if t.root == nil {
//		return
//	}
//	t.root.Traverse()
//	println("\n")
//}
