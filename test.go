package main

func main() {
	nodeG := Node{data: "g", left: nil, right: nil}
	nodeF := Node{data: "f", left: &nodeG, right: nil}
	nodeE := Node{data: "e", left: nil, right: nil}
	nodeD := Node{data: "d", left: &nodeE, right: nil}
	nodeC := Node{data: "c", left: nil, right: nil}
	nodeB := Node{data: "b", left: &nodeD, right: &nodeF}
	nodeA := Node{data: "a", left: &nodeB, right: &nodeC}

}

