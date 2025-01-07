package main

import (
	"fmt"
	"sort"
)
// Node structure 
type Node struct {
	keys     []int     
	children []*Node   
	isLeaf   bool      
	next     *Node     
}
// B+ Tree structure
type BPlusTree struct {
	root  *Node 
	order int   
}

// Create a new node
func newNode(isLeaf bool) *Node {
	return &Node{
		keys:     []int{},
		children: []*Node{},
		isLeaf:   isLeaf,
		next:     nil,
	}
}

// Create a new B+ Tree
func newBPlusTree(order int) *BPlusTree {
	return &BPlusTree{
		root:  nil,
		order: order,
	}
}

// Insert a key into the B+ Tree
func (tree *BPlusTree) Insert(key int) {
	if tree.root == nil {
		tree.root = newNode(true)
		tree.root.keys = append(tree.root.keys, key)
		return
	}
	parent, leaf := tree.findLeafNode(key)
	tree.insertInLeaf(leaf, key)

	if len(leaf.keys) >= tree.order {
		tree.splitNode(parent, leaf)
	}
}

// Find the leaf node where a key should be inserted
func (tree *BPlusTree) findLeafNode(key int) (*Node, *Node) {
	var parent *Node
	current := tree.root

	// Traverse the tree to find the leaf node
	for !current.isLeaf {
		parent = current
		idx := 0
		for idx < len(current.keys) && key >= current.keys[idx] {
			idx++
		}
		current = current.children[idx]
	}
	return parent, current
}

// Insert a key into a leaf node in sorted order
func (tree *BPlusTree) insertInLeaf(leaf *Node, key int) {
	idx := sort.SearchInts(leaf.keys, key)
	leaf.keys = append(leaf.keys[:idx], append([]int{key}, leaf.keys[idx:]...)...)
}



func (tree *BPlusTree) splitNode(parent, node *Node) {
	mid := len(node.keys) / 2
	newNode := newNode(node.isLeaf)

	// Move keys to the new node
	newNode.keys = append(newNode.keys, node.keys[mid:]...)
	node.keys = node.keys[:mid]

	// Handle leaf node split
	if node.isLeaf {
		newNode.next = node.next
		node.next = newNode
	} else { // Handle internal node split
		newNode.children = append(newNode.children, node.children[mid:]...)
		node.children = node.children[:mid+1]
	}

	// If there's no parent, create a new root
	if parent == nil {
		newRoot := newNode(false)
		newRoot.keys = append(newRoot.keys, newNode.keys[0])
		newRoot.children = append(newRoot.children, node, newNode)
		tree.root = newRoot
	} else {
		tree.insertInParent(parent, newNode, newNode.keys[0])
	}
}

// Insert key into parent after splitting
func (tree *BPlusTree) insertInParent(parent, newNode *Node, key int) {
	idx := sort.SearchInts(parent.keys, key)
	parent.keys = append(parent.keys[:idx], append([]int{key}, parent.keys[idx:]...)...)
	parent.children = append(parent.children[:idx+1], append([]*Node{newNode}, parent.children[idx+1:]...)...)

	// Handle parent overflow
	if len(parent.keys) >= tree.order {
		grandParent, _ := tree.findLeafNode(parent.keys[0])
		tree.splitNode(grandParent, parent)
	}
}


// Display keys in sorted order using leaf traversal
func (tree *BPlusTree) Display() {
	current := tree.root

	// Move to the first leaf node
	for !current.isLeaf {
		current = current.children[0]
	}

	// Print keys in sequential order
	for current != nil {
		fmt.Print(current.keys, " ")
		current = current.next
	}
	fmt.Println()
}

// Search for a key in the B+ Tree
func (tree *BPlusTree) Search(key int) bool {
	current := tree.root

	for current != nil {
		if current.isLeaf {
			for _, k := range current.keys {
				if k == key {
					return true
				}
			}
			return false
		}

		// Traverse to the appropriate child
		idx := sort.SearchInts(current.keys, key)
		current = current.children[idx]
	}
	return false
}


func main() {
	var order, n, key int

	// Get tree order from the user
	fmt.Print("Enter the order of the B+ Tree: ")
	fmt.Scan(&order)

	// Create a new B+ Tree
	bpt := newBPlusTree(order)

	// Insert keys into the tree
	fmt.Print("Enter the number of keys to insert: ")
	fmt.Scan(&n)

	fmt.Println("Enter the keys:")
	for i := 0; i < n; i++ {
		fmt.Scan(&key)
		bpt.Insert(key)
	}

	// Display keys
	fmt.Println("B+ Tree keys in sorted order:")
	bpt.Display()

	// Search keys in the tree
	var searchKey int
	fmt.Print("Enter a key to search: ")
	fmt.Scan(&searchKey)
	if bpt.Search(searchKey) {
		fmt.Printf("Key %d found in the tree.\n", searchKey)
	} else {
		fmt.Printf("Key %d not found in the tree.\n", searchKey)
	}
}
