package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const menu_path = "files/menu.json"
const stories_path = "stories/"

var game_state = make(map[string]string)
var story = make(map[string]Screen)

type Screen struct {
	id, text string
	choices  []string
}

type Story struct {
	title   string
	screens map[string]Screen
}

func main() {
	// initialization steps
	// load the menu
	fmt.Println("Initializing")
	var menu Story = load(menu_path)
	fmt.Println("Loaded menu")
	fmt.Println(menu)
	// set the initial state
	fmt.Println("Setting up state")
	var item = make(map[string]string)
	item["current_story"] = "menu"
	item["position"] = "1"
	updatestate(item, false)
	fmt.Println("State updated")
	fmt.Println(game_state)
	// display the menu to the human
	var current_screen = menu.screens[game_state["position"]]
	render(current_screen)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
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
	println(screen.text)
	for c := 0; c < len(screen.choices); c++ {
		fmt.Println(screen.choices[c])
	}
}

func load(filename string) Story {
	file, _ := os.ReadFile(filename)
	var data Story
	err := json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Cannot unmarshal the json %d\n", err)
	}
	return data
}
