package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) printStory(depth int) {
	for i := 0; i < depth*2; i++ {
		fmt.Println("  ")
	}
	fmt.Println(node.text)
	fmt.Println()
	// handles nil to avoid crash
	if node.yesPath != nil {
		node.yesPath.printStory(depth + 1) // adds a blank space for each deapth in the story
	}
	if node.noPath != nil {
		node.noPath.printStory(depth + 1) // adds a blank space for each deapth in the story
	}
}

func (node *storyNode) play() {
	fmt.Println(node.text)

	// to break the program
	if node.yesPath != nil && node.noPath != nil {

		scanner := bufio.NewScanner(os.Stdin)

		for {
			scanner.Scan()
			answer := scanner.Text()

			if answer == "yes" {
				node.yesPath.play()
				break
			} else if answer == "no" {
				node.noPath.play()
				break
			} else {
				fmt.Println("That answer was not an option!")
			}
		}
	}
}

func main() {
	root := storyNode{"You saw a hot girl. Do you want to go talk to her?", nil, nil}
	winning := storyNode{"She's a feminist. You dodged a bullet there. You won", nil, nil}
	losing := storyNode{"She's a feminist, she accuse you of sexual harrassment. You lost", nil, nil}
	root.yesPath = &losing
	root.noPath = &winning

	root.printStory(0)
	root.play()
}
