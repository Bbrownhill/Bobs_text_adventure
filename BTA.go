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

type Choice struct {
	Id, Text, Target string
}

type Screen struct {
	Id, Text string
	Choices  map[string]Choice
}

type Story struct {
	Title   string
	Screens map[string]Screen
}

func main() {
	// initialization steps
	// load the menu
	var menu Story = load(menu_path)

	// set the initial state
	var item = make(map[string]string)
	item["current_story"] = "menu"
	item["position"] = "1"
	updatestate(item, false)
	// display the menu to the human
	var current_screen = menu.Screens[game_state["position"]]
	render(current_screen)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if input == "exit" {
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
	fmt.Println(screen.Text)
	for _, v := range screen.Choices {
		fmt.Println(v.Text)
	}
}

func load(filename string) Story {
	file, _ := os.ReadFile(filename)
	var data Story // map[string]any
	// var story Story
	err := json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Cannot unmarshal the json %d\n", err)
	}
	return data
}
