package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const menu_path = "files/menu.json"
const stories_path = "stories/"

var clear map[string]func() //create a map for storing clear funcs
var special_functions = make(map[string]func())
var game_state = make(map[string]string)
var stories = make(map[string]Story)

// var menu Story

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

func init() {
	// determine what clear screen your OS most likely uses.
	fetch_os_clear_method()

	// load the menu
	var menu = load(menu_path)
	stories[menu.Title] = menu
	var item = make(map[string]string)

	// set the initial state
	item["current_story"] = menu.Title
	item["position"] = "1"
	updatestate(item, false)
}

func main() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	render(current_screen)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		val, ok := current_screen.Choices[input]
		if !ok {
			// If the player picks an invalid option simply skip the loop
			// this will result in the player staying on the same screen.
			fmt.Println("Please select a valid option")
			continue
		}

		//TO DO
		// upgrade this to accomodate multiple functions
		if val.Target == "Exit Game" {

			break // Exit loop if an empty line is entered
		}
		// prepare an update for the gamestate
		var next_screen = make(map[string]string)
		next_screen["position"] = val.Target
		updatestate(next_screen, false)

		current_screen = current_story.Screens[game_state["position"]]
		render(current_screen)
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
	CallClear()
	fmt.Println(screen.Text)

	// As Go maps are unordered this code will iterate by Choice ID
	// this is being done so options always render in the correct order
	var choice_count = len(screen.Choices)
	for choice := 1; choice <= choice_count; choice++ {
		sc := strconv.Itoa(choice) // convert the int to a string of the int before grabbing the text
		fmt.Println(screen.Choices[sc].Text)
	}
	// I do this because string(1) returns ascii character 1 and not "1"
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

func fetch_os_clear_method() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen")
	}
}
