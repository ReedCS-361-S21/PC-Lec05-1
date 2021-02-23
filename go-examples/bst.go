package main

import (
	"fmt"
	"strconv"
)

type bstNode struct {
	key   int
	left  *bstNode
	right *bstNode
}

type bst struct {
	root *bstNode
}

func contains(t *bst, x int) bool {
	var n *bstNode = t.root
	for n != nil {
		if x == n.key {
			return true
		} else if x < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}
	return false
}

func insertHelper(x int, np **bstNode) {
	if (*np) == nil {
		(*np) = &bstNode{key: x, left: nil, right: nil}
	} else if x == (*np).key {
		return
	} else if x < (*np).key {
		insertHelper(x, &(*np).left)
	} else if x > (*np).key {
		insertHelper(x, &(*np).right)
	}
}

func insert(t *bst, x int) {
	insertHelper(x, &t.root)
}

func makeBST(xs []int) *bst {
	var t *bst = &bst{root: nil}
	for _, x := range xs {
		insert(t, x)
	}
	return t
}

func stringHelper(n *bstNode) string {
	if n == nil {
		return "Lf"
	} else {
		var ks string = strconv.Itoa(n.key)
		var ls string = stringHelper(n.left)
		var rs string = stringHelper(n.right)
		return "Br" + "(" + ks + "," + ls + "," + rs + ")"
	}
}

func String(t *bst) string {
	return stringHelper(t.root)
}

func main() {
	var xs []int = []int{5, 3, 1, 2, 6, 9, 8}
	var t *bst = makeBST(xs)
	fmt.Println(String(t))
	var i int = 0
	for i < 10 {
		if contains(t, i) {
			fmt.Printf("The tree above contains %d.\n", i)
		} else {
			fmt.Printf("The tree above does not contain %d.\n", i)
		}
		i++
	}
}
