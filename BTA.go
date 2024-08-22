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
	"strings"
	"time"
)

const menu_path = "files/menu.json"
const stories_path = "stories/"
const save_dir = "saves/"

var clear map[string]func() //create a map for storing clear funcs
var screen_functions = make(map[string]func())
var target_actions = make(map[string]func(string))
var game_state = make(map[string]string)
var stories = make(map[string]Story)
var save_files = make(map[int]string)

// var menu Story

type Choice struct {
	Id, Text string
	Target   map[string]string
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
	screen_functions["Display Save Files"] = Display_Save_Files

	target_actions["Next Screen"] = Next_Screen   //Move to the next screen
	target_actions["Change Story"] = Change_Story //start a new story
	target_actions["Load Game"] = Load_Game       //load a saved game
	target_actions["Save Game"] = Save_Game       //save a current game
	target_actions["Exit Game"] = Exit_Game       //exit the game

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
	saves, err := os.ReadDir(save_dir)
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range saves {

		save_files[i] = s.Name()
	}
	game_state["game_state"] = "Running"
}

func main() {
	//initial render for the main menu is required
	render(stories[game_state["current_story"]].Screens[game_state["position"]])
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		current_story := stories[game_state["current_story"]]
		current_screen := current_story.Screens[game_state["position"]]

		val, ok := current_screen.Choices[input]
		if !ok {
			// If the player picks an invalid option simply skip the loop
			// this will result in the player staying on the same screen.
			fmt.Println("Please select a valid option")
			continue
		}
		execute_target_action(val.Target)
		current_story = stories[game_state["current_story"]]
		current_screen = current_story.Screens[game_state["position"]]
		render(current_screen)

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func execute_target_action(input map[string]string) {
	// TODO consider ordering the actions
	for k, v := range input {
		target_actions[k](v)
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
	if screen.Function != "" {
		screen_functions[screen.Function]()
	}
	fmt.Println(screen.Text)
	// for _, c := range []byte(screen.Text) {
	// 	time.Sleep(30 * time.Millisecond)
	// 	fmt.Print(string(c))
	// }

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

// screen functions
func Display_Stories() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	var index = len(current_screen.Choices) + 1 // I want the index to start one higher than the number of choices
	for title := range stories {
		// special case, check if the story is main menu
		// we don't want to display an option to load the main menu when were already there
		if title == "Main Menu" {
			continue
		}
		var new_choice = Choice{
			Id:     strconv.Itoa(index),
			Text:   fmt.Sprintf("%v. %v", strconv.Itoa(index), title),
			Target: map[string]string{"Change Story": title},
		}
		current_screen.Choices[strconv.Itoa(index)] = new_choice
		index++
	}
}

func Display_Save_Files() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	var index = len(current_screen.Choices) + 1 // I want the index to start one higher than the number of choices
	for i, save := range save_files {
		// special case, check if the story is main menu
		// we don't want to display an option to load the main menu when were already there

		var new_choice = Choice{
			Id:     strconv.Itoa(index),
			Text:   fmt.Sprintf("%v. %v", strconv.Itoa(index), save),
			Target: map[string]string{"Load Game": save[i]},
		}
		current_screen.Choices[strconv.Itoa(index)] = new_choice
		index++
	}
}

// target functions
func Next_Screen(screen string) {
	var next_screen = make(map[string]string)
	next_screen["position"] = screen
	updatestate(next_screen, false)
}

func Change_Story(name string) {
	game_state["current_story"] = name
	game_state["position"] = "1"
}

func Load_Game(name string) {

}

func Save_Game(name string) {
	// using date time stamps for save game names for now
	t := time.Now().Format(time.StampMilli)
	t = strings.Replace(t, ":", "-", -1) // need to replace : with - so the name is valid
	content, err := json.Marshal(game_state)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf(save_dir+"%v.json", t), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Exit_Game(name string) {
	// I don't want to cause a panic...
	CallClear()
	fmt.Println("Thank you for playing! :)")
	os.Exit(0)
}
