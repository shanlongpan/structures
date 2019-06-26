package main

import (
	"fmt"
)

type Node struct {
	data  string
	left  *Node
	right *Node
}
type seqStack struct {
	data [4]*Node
	tag [4]int // 后续遍历准备
	top int // 数组下标
}
func main() {
	nodeG := Node{data: "g", left: nil, right: nil}
	nodeF := Node{data: "f", left: &nodeG, right: nil}
	nodeE := Node{data: "e", left: nil, right: nil}
	nodeD := Node{data: "d", left: &nodeE, right: nil}
	nodeP := Node{data: "p", left: nil, right: nil}
	nodeH := Node{data: "h", left: nil, right: nil}
	nodeC := Node{data: "c", left: &nodeP, right: &nodeH}
	nodeB := Node{data: "b", left: &nodeD, right: &nodeF}
	nodeA := Node{data: "a", left: &nodeB, right: &nodeC}

	//fmt.Println(breadthFirstSearch(nodeA))
	fmt.Println(preOrderLoop(&nodeA))
}
//广度优先遍历
func breadthFirstSearch(node Node) []string {
	var result []string
	var nodes  = []Node{node}

	for len(nodes) > 0 {
		node := nodes[0]
		nodes = nodes[1:]
		result = append(result, node.data)
		if node.left != nil {
			nodes = append(nodes, *node.left)
		}
		if node.right != nil {
			nodes = append(nodes, *node.right)
		}
	}
	return result
}
//递归版 前 中 后序遍历
func preOrderRecursive(node Node) {
	//前序 前序遍历按照“根结点-左孩子-右孩子”的顺序进行访问。
	//fmt.Println(node.data)
	if node.left != nil {
		preOrderRecursive(*node.left)
	}
	// 在这里输出就是中序   中序遍历按照“左孩子-根结点-右孩子”的顺序进行访问。
	fmt.Println(node.data)
	if node.right != nil {
		preOrderRecursive(*node.right)
	}
	// 在这里输出是后序   后序遍历按照“左孩子-右孩子-根结点”的顺序进行访问。
}

//非递归版前序遍历

func preOrderLoop(node *Node) (result []string) {
	var s seqStack
	s.top = -1 // 空
	if node == nil {
		panic("no data here")
	}else {
		for node != nil || s.top != -1 {
			//if node !=nil{
			//	fmt.Println(s.top,node.data)
			//}else{
			//	fmt.Println(s.top)
			//}

			for node != nil {
				result = append(result, node.data)
				s.top++
				s.data[s.top] = node
				node = node.left
			}
			s.top--
			node = s.data[s.top + 1]
			node = node.right
		}
	}
	return
}

//非递归中序遍历

func midOrderLoop(node *Node) (result []string) {
	var s seqStack
	s.top = -1
	if node == nil {
		panic("no data here")
	}else {
		for node != nil || s.top != -1 {
			for node != nil {
				s.top++
				s.data[s.top] = node
				node = node.left
			}
			s.top--
			node = s.data[s.top + 1]
			result = append(result, node.data)
			node = node.right
		}
	}
	return
}

//非递归后序遍历

func postOrderLoop(node *Node) (result []string)  {
	var s seqStack
	s.top = -1

	if node == nil {
		panic("no data here")
	}else {
		for node != nil || s.top != -1 {
			for node != nil {
				s.top++
				s.data[s.top] = node
				s.tag[s.top] = 0
				node = node.left
			}

			if s.tag[s.top] == 0 {
				node = s.data[s.top]
				s.tag[s.top] = 1
				node = node.right
			}else {
				for s.tag[s.top] == 1 {
					s.top--
					node = s.data[s.top + 1]
					result = append(result, node.data)
					if s.top < 0 {
						break
					}
				}
				node = nil
			}
		}
	}
	return
}