// Package parser contains all DOM traversal functions used by other packages.
package parser

import (
	"errors"
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// See https://pkg.go.dev/golang.org/x/net/html#example-Parse
func FindIngredientList(node *html.Node) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "wprm-recipe-ingredients" {
				return node, nil
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		list, _ := FindIngredientList(c)
		if list != nil {
			return list, nil
		}
	}
	return nil, errors.New("ingredient list does not exist")
}

func PrintIngredientList(list *html.Node) {
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		for child := li.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "span" {
				for _, a := range child.Attr {
					if a.Val == "wprm-recipe-ingredient-name" {
						// The first child should be a text node
						fmt.Println(child.FirstChild.Data)
					}
				}
			}
		}
	}
}

// Technically don't need to return an error because we can just check for nil
func FindRecipeCard(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "div" {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "wprm-recipe-container" {
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		card := FindRecipeCard(c)
		if card != nil {
			return card
		}
	}
	return nil
}

func FindInstructionsList(node *html.Node) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "wprm-recipe-instructions" {
				return node, nil
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		list, _ := FindInstructionsList(c)
		if list != nil {
			return list, nil
		}
	}
	return nil, errors.New("instructions list does not exist")
}

// Returns the first element underneath and including `node` that has the given class value (as given in the HTML)
func GetElementWithClass(node *html.Node, tagname, class string) *html.Node {
	if node.Type == html.ElementNode && node.Data == tagname {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == class {
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		element := GetElementWithClass(c, tagname, class)
		if element != nil {
			return element
		}
	}
	return nil
}

// Returns the first text node under the given node
func GetTextNode(node *html.Node) *html.Node {
	if node.Type == html.TextNode {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		textNode := GetTextNode(c)
		if textNode != nil {
			return textNode
		}
	}
	return nil
}

// TraverseNode collecting the nodes that match the given function
// See https://gist.github.com/Xeoncross/8bbb84bc4bf540bd907f79ee17c4e1fc
func TraverseNode(doc *html.Node, matcher func(node *html.Node) (bool, bool)) (nodes []*html.Node) {
	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nodes
}

func FindIngredientLists(node *html.Node) []*html.Node {
	matcher := func(node *html.Node) (keep, exit bool) {
		if node.Type == html.ElementNode && node.DataAtom == atom.Ul {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == "wprm-recipe-ingredients" {
					keep = true
					exit = true // No nested ingredients lists
				}
			}
		}
		return
	}
	return TraverseNode(node, matcher)
}

func PrintNode(node *html.Node) {
	fmt.Print("Node Type: ")
	switch node.Type {
	case html.ElementNode:
		fmt.Println("Element")
	case html.TextNode:
		fmt.Println("Text")
	default:
		fmt.Println("Other")
	}

	fmt.Println("Node Data:", node.Data)

	fmt.Println("Node Attributes")
	for _, a := range node.Attr {
		fmt.Println(a.Key, a.Val)
	}
}
