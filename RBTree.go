package main

import (
	"fmt"
)

type RBTree struct {
	root     *node
	numCount int64
}

//node 节点为一个结构体
type node struct {
	color      bool
	leftNode   *node
	rightNode  *node
	fatherNode *node
	value      int
}

//红黑通过bool类型来存储，并设置常量
const (
	RBTRed   = false
	RBTBlack = true
)

func NewRBtree() *RBTree {
	return &RBTree{}
}

//find 查找
func (t *RBTree) find(v int) *node {
	n := t.root
	for n != nil {
		if v < n.value {
			//小于当前节点的话，往左节点找
			n = n.leftNode
		} else if v > n.value {
			//大于当前节点的话，往右节点找
			n = n.rightNode
		} else {
			//等于的话表示找到，返回
			return n
		}
	}
	//循环结束没找到，返回
	return nil
}

//          |                                  |
//          X                                  Y
//         / \         left rotate            / \
//        α  Y       ------------->          X   γ
//          / \                             / \
//         β  γ                            α  β
//     y 替换X

func (t *RBTree) leftRotate(n *node) {
	rn := n.rightNode
	//first give n's father to rn's father
	rn.fatherNode = n.fatherNode
	if n.fatherNode != nil {
		if n.fatherNode.leftNode == n {
			n.fatherNode.leftNode = rn
		} else {
			n.fatherNode.rightNode = rn
		}
	} else {
		t.root = rn
	}
	n.fatherNode = rn
	n.rightNode = rn.leftNode
	if n.rightNode != nil {
		n.rightNode.fatherNode = n
	}
	rn.leftNode = n
}

//          |                                  |
//          X                                  Y
//         / \         right rotate           / \
//        Y   γ      ------------->          α  X
//       / \                                   / \
//      α  β                                  β  γ

//     Y替换X

func (t *RBTree) rightRotate(n *node) {

	ln := n.leftNode
	ln.fatherNode = n.fatherNode
	if n.fatherNode != nil {
		if n.fatherNode.leftNode == n {
			n.fatherNode.leftNode = ln
		} else {
			n.fatherNode.rightNode = ln
		}
	} else {
		t.root = ln
	}
	n.fatherNode = ln

	n.leftNode = ln.rightNode
	if n.leftNode != nil {

		n.leftNode.fatherNode = n
	}
	ln.rightNode = n

}

//判断是否为黑，空为黑
func isBlack(n *node) bool {
	if n == nil {
		return true
	} else {
		return n.color == RBTBlack
	}
}

//寻找兄弟节点
func findBroNode(n *node) (bro *node) {
	if n.fatherNode == nil {
		return nil
	}

	if n.fatherNode.leftNode == n {
		bro = n.fatherNode.rightNode
	} else {
		bro = n.fatherNode.leftNode
	}
	return bro
}

func (t *RBTree) insert(v int) {
	//如果根节点为nil，则先插入新的根节点。 第一种情况最简单 case1
	if t.root == nil {
		t.numCount++
		t.root = &node{value: v, color: RBTBlack}
		return
	}

	n := t.root
	//新插入的节点为红色
	insertNode := &node{value: v, color: RBTRed}

	//标记父节点，找到父节点追加到父节点下面
	var nf *node

	//以下代码找到插入位置的父节点
	for n != nil {
		nf = n
		if v < n.value {
			n = n.leftNode
		} else if v > n.value {
			n = n.rightNode
		} else {
			//已经存在，返回
			return
		}
	}
	t.numCount++
	//设置新插入节点的父节点
	insertNode.fatherNode = nf
	//将新的节点挂到父节点上
	if v < nf.value {
		nf.leftNode = insertNode
	} else {
		nf.rightNode = insertNode
	}
	t.insertFixUp(insertNode)
}

//父节点为红色，分为两大类，1.兄弟节点是否是红色，红色的话父节点和叔叔节点变黑，祖父节点变红,然后往上遍历
//2.兄弟节点是黑色，分别为两个小类

func (t *RBTree) insertFixUp(n *node) {
	//case 2 父节点是黑色不用修复，
	for !isBlack(n.fatherNode) {
		uncleNode := findBroNode(n.fatherNode)
		if !isBlack(uncleNode) {
			//case3 父节点是红，叔叔节点红
			n.fatherNode.color = RBTBlack
			uncleNode.color = RBTBlack
			uncleNode.fatherNode.color = RBTRed
			n = n.fatherNode.fatherNode //祖父节点最为新的插入节点
		} else if n.fatherNode == n.fatherNode.fatherNode.leftNode {
			if n == n.fatherNode.leftNode {
				//case 4 在一个方向 ，变色后旋转祖父节点
				n.fatherNode.fatherNode.color = RBTRed
				n.fatherNode.color = RBTBlack
				n = n.fatherNode.fatherNode
				t.rightRotate(n)
				//旋转后n是的父节点必定是黑，n现在是原来n的祖父节点，n的父节点是原来的叔叔节点
			} else {
				//case 5 不在一个方向，直接旋转父节点
				n = n.fatherNode
				t.leftRotate(n)
				//旋转后n变成了原来的父节点，他的父节点还是红色，变成case4
			}
		} else {
			if n == n.fatherNode.rightNode {
				n.fatherNode.fatherNode.color = RBTRed
				n.fatherNode.color = RBTBlack
				n = n.fatherNode.fatherNode
				t.leftRotate(n)

			} else {
				n = n.fatherNode
				t.rightRotate(n)
			}
		}
		t.root.color = RBTBlack
	}
}

//后继节点，右边最小的数
func (t *RBTree) miniNum(n *node) *node {
	for n.leftNode != nil {
		n = n.leftNode
	}
	return n
}

//交换颜色
func (t *RBTree) changeColor(u, v *node) {
	u.color, v.color = v.color, u.color
}

//v替换节点u
func (t *RBTree) transplant(u, v *node) {

	if u.fatherNode == nil {
		t.root = v
		if v != nil {
			v.fatherNode = nil
		}
	} else if u == u.fatherNode.leftNode {
		u.fatherNode.leftNode = v
	} else {
		u.fatherNode.rightNode = v
	}
	if v != nil {
		v.fatherNode = u.fatherNode
	}
}

func (t *RBTree) delete(v int) {
	n := t.find(v)
	var child *node
	if n == nil {
		return
	}
	t.numCount--
	//两个子节点的时候找后继结点删除
	if n.leftNode != nil && n.rightNode != nil {
		successor := t.miniNum(n.rightNode)
		n.value = successor.value
		n = successor
	}
	var fixColor = n.color
	if n.leftNode == nil {
		child = n.rightNode
	} else {
		child = n.leftNode
	}

	t.transplant(n, child)
	//n 节点虽然被child节点替换了，由于child节点有可能是nil，所以child没有保存father节点，n还保存着father节点和child节点的地址
	//所以用 n 修复作为修复节点
	if fixColor == RBTBlack {
		t.deleteFixUp(n)
	}

}

func (t *RBTree) deleteFixUp(n *node) {

	if t.root == nil {
		return
	}
	//if t.root == n {
	//	t.root = nil
	//	return
	//}
	if n != nil {
		n.color = RBTBlack
		return
	}
	t.fixCase3(n)
}

// sibling 红
func (t *RBTree) fixCase3(node *node) {
	fmt.Println(node)
	bro := findBroNode(node)
	if bro.color == RBTRed {
		node.fatherNode.color = RBTRed
		bro.color = RBTBlack
		if bro == node.fatherNode.leftNode {
			t.leftRotate(node.fatherNode)
		} else {
			t.rightRotate(node.fatherNode)
		}
	} else {
		t.fixCase4(node)
	}
}

//sibling 黑 SL or SR 是红
func (t *RBTree) fixCase4(node *node) {
	bro := findBroNode(node)
	if bro == node.fatherNode.leftNode {
		if bro.rightNode != nil {
			//变色旋转
			t.changeColor(bro.leftNode, bro)
			t.rightRotate(bro)
			t.changeColor(node.fatherNode, bro)
			t.leftRotate(node.fatherNode)
			return
		} else if bro.leftNode != nil {
			t.changeColor(node.fatherNode, bro)
			t.leftRotate(node.fatherNode)
			return
		}
	} else {
		if bro.rightNode != nil {
			t.changeColor(node.fatherNode, bro)
			t.rightRotate(node.fatherNode)
			return
		} else if bro.leftNode != nil {
			//变色旋转
			t.changeColor(bro.leftNode, bro)
			t.leftRotate(bro)
			t.changeColor(node.fatherNode, bro)
			t.rightRotate(node.fatherNode)
			return
		}
	}
	t.fixCase5(node)
}

//p 红 p 和sibling 换颜色
func (t *RBTree) fixCase5(node *node) {
	if node.fatherNode.color == RBTRed {
		bro := findBroNode(node)
		t.changeColor(node.fatherNode, bro)
	} else {
		t.fixCase6(node)
	}
}

//p && sibling && SL &&SR 都是黑
func (t *RBTree) fixCase6(node *node) {
	bro := findBroNode(node)
	bro.color = RBTRed
	node = node.fatherNode
	t.fixCase3(node)
}
func main() {
	sliceInsert := []int{17, 8, 3, 6, 16, 10, 13, 14, 5, 18, 15, 4, 19, 7, 11, 0, 2, 9, 1, 12, 100, 99, 33, 20}
	rbTree := NewRBtree()
	for _, val := range sliceInsert {
		rbTree.insert(val)
	}
	preOrderRecursives(*rbTree.root)

	for _, val := range sliceInsert {
		rbTree.delete(val)
		fmt.Println("==========================节点个数", rbTree.numCount)
		if rbTree.numCount > 0 {
			preOrderRecursives(*rbTree.root)
		}
	}

}

func preOrderRecursives(nodes node) {
	//前序 前序遍历按照“根结点-左孩子-右孩子”的顺序进行访问。
	//fmt.Println(node.value)
	if nodes == (node{}) {
		return
	}
	if nodes.leftNode != nil {
		preOrderRecursives(*nodes.leftNode)
	}
	// 在这里输出就是中序   中序遍历按照“左孩子-根结点-右孩子”的顺序进行访问。
	fmt.Println(nodes.value)
	if nodes.rightNode != nil {
		preOrderRecursives(*nodes.rightNode)
	}
	// 在这里输出是后序   后序遍历按照“左孩子-右孩子-根结点”的顺序进行访问。
}
