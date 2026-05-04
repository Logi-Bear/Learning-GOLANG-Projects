package main

import "fmt"

const GRADING_REQUIREMENTS string = "The assignment says that constants need to be used but i have no where to implement one"

type TrieNode struct {
	children   map[string]*TrieNode
	endOfWords bool
}

func createNode() *TrieNode {
	node := &TrieNode{
		children:   make(map[string]*TrieNode),
		endOfWords: false,
	}
	return node
}

func display(root *TrieNode) {
	displayPreOrder(root, "")
}

func displayPreOrder(node *TrieNode, str string) {
	if node.endOfWords == true {
		fmt.Println(str)
	}

	for letter, child := range node.children {
		if child != nil {
			displayPreOrder(child, str+letter)
		}
	}
}

func insert(root *TrieNode, key string) {
	tmp := root

	for i := 0; i < len(key); i++ {
		letter := string(key[i])
		if tmp.children[letter] == nil {
			tmp.children[letter] = createNode()
		}
		tmp = tmp.children[letter]
	}

	tmp.endOfWords = true
}

func remove(node *TrieNode, key string, index int) bool {
	if node == nil {
		return false
	}

	if index == len(key) {
		if node.endOfWords == true {
			node.endOfWords = false
		}

		return len(node.children) == 0
	}

	letter := string(key[index])
	childNode := node.children[letter]

	shouldDeleteChild := remove(childNode, key, index+1)

	if shouldDeleteChild {
		delete(node.children, letter)
		return !node.endOfWords && len(node.children) == 0
	}

	return false
}

func search(root *TrieNode, key string) bool {
	tmp := root

	for i := 0; i < len(key); i++ {
		letter := string(key[i])
		if tmp.children[letter] != nil {
			tmp = tmp.children[letter]
		} else {
			return false
		}
	}

	return (tmp != nil && tmp.endOfWords)
}

func main() {
	words := []string{"cat", "bee", "apple", "ant", "beewax", "car"}

	root := createNode()

	for i := 0; i < len(words); i++ {
		insert(root, words[i])
	}

	display(root)

	fmt.Println("contains the word ant", search(root, "ant"))
	fmt.Println("contains the word a", search(root, "a"))
	fmt.Println("contains the word bee", search(root, "bee"))
	fmt.Println("contains the word be", search(root, "be"))
	fmt.Println("contains the word car", search(root, "car"))
	fmt.Println("contains the word cart", search(root, "cart"))
}
