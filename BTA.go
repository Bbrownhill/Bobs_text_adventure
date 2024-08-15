package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

const menu_path = "files/menu.json"
const stories_path = "stories/"

var clear map[string]func() //create a map for storing clear funcs
var screen_functions = make(map[string]func())
var target_actions = make(map[string]func())
var game_state = make(map[string]string)
var stories = make(map[string]Story)

// var menu Story

type Choice struct {
	Id, Text, Target string
}

type Screen struct {
	Id, Function, Text string
	Choices            map[string]Choice
}

type Story struct {
	Title   string
	Screens map[string]Screen
}

func init() {
	// determine what clear screen your OS most likely uses.
	fetch_os_clear_method()
	screen_functions["Display Stories"] = Display_Stories
	screen_functions["Load Stories"] = Load_Stories
	screen_functions["Display Save Files"] = Display_Save_Files
	screen_functions["Exit Game"] = Exit_Game

	target_actions["Load Story"] = Load_Story //start a new story
	target_actions["Load Game"] = Load_Game   //load a saved game
	target_actions["Save Game"] = Save_Game   //save a current game
	target_actions["Exit Game"] = Exit_Game   //exit the game
	target_actions["Exit Story"] = Exit_Story //exit story, return to menu

	// load the menu
	var menu = load(menu_path)
	stories[menu.Title] = menu
	var item = make(map[string]string)

	// set the initial state
	item["current_story"] = menu.Title
	item["position"] = "1"
	updatestate(item, false)

	//load stories in the stories dir
	entries, err := os.ReadDir(stories_path)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		var story = load(stories_path + e.Name())
		stories[story.Title] = story
	}
	game_state["game_state"] = "Running"
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
		execute_target_action(val.Target)
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

		// check if the screen has any functions that need executed.
		if current_screen.Function != "" {
			screen_functions[current_screen.Function]()
		}
		render(current_screen)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

}

// primes := [6]int{2, 3, 5, 7, 11, 13}

func execute_target_action(input string) string {

	numeric_targets, _ := regexp.Compile("[0-9]")
	var matchers []regexp.Regexp
	matchers = append(matchers, *numeric_targets)

	for k, _ := range target_actions {

		new_match, _ := regexp.Compile(fmt.Sprintf("^%v", k)) //Compile a new regex to match each target action key
		matchers = append(matchers, *new_match)
	}

	for i, v := range matchers {
		splitter, _ := regexp.Compile(":")
		if v.Match([]byte(input)) {
			target_components := splitter.Split(input, 2)
			target_actions[target_components[0]](target_components[1])
		}
	}

	return ""
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

func Display_Stories() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	var index = len(current_screen.Choices) + 1 // I want the index to start one higher than the number of choices
	for title, _ := range stories {
		// special case, check if the story is main menu
		// we don't want to display an option to load the main menu when were already there
		if title == "Main Menu" {
			continue
		}
		var new_choice = Choice{
			Id:     strconv.Itoa(index),
			Text:   fmt.Sprintf("%v. %v", strconv.Itoa(index), title),
			Target: fmt.Sprintf("Load Story: %v", title),
		}
		current_screen.Choices[strconv.Itoa(index)] = new_choice
		index++
	}
}

func Load_Stories() {
	game_state["current_story"] = game_state["next_story"]
	game_state["position"] = "1"
}

func Display_Save_Files() {

}

func Load_Story() {

}

func Load_Game() {

}

func Save_Game() {

}

func Exit_Story() {

}

func Exit_Game() {

}
