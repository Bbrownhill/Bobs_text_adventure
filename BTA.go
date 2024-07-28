package main

import (
	"bufio"
	"fmt"
	"os"
)

var game_state = make(map[string]string)
var story = make(map[string]Screen)

type Screen struct {
	id, text string
	choices  [8]string
}

func main() {

	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text() // Get the current line of text
		if text == "exit" {
			break // Exit loop if an empty line is entered
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

}

func updatestate(item map[string]string, remove bool) {
	if remove {
		for k, _ := range item {
			delete(game_state, k)
		}
	} else {
		for k, v := range item {
			game_state[k] = v
		}
	}
}

func render(screen Screen) {

}
