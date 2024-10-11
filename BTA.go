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
	"strings"
	"time"
)

const menu_path = "files/menu.json"
const stories_path = "stories/"
const save_dir = "saves/"

// regex's for parsing inline functions
var reg = regexp.MustCompile("{.*}")

var clear map[string]func() //create a map for storing clear funcs
var screen_functions = make(map[string]func())
var target_actions = make(map[string]func(string))
var game_state = make(map[string]string)
var stories = make(map[string]Story)
var save_slots = []string{"Save Slot 1", "Save Slot 2", "Save Slot 3", "Save Slot 4", "Save Slot 5"} //setting a limit of 5 save slots for now

// basic ANSI colours
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// var menu Story

type Choice struct {
	Id, Text string
	Target   map[string]string
}

type Screen struct {
	Id, Function string
	Text         []string
	Choices      map[string]Choice
}

type Story struct {
	Title   string
	Screens map[string]Screen
}

func init() {
	// determine what clear screen your OS most likely uses.
	fetch_os_clear_method()
	screen_functions["Display Stories"] = Display_Stories
	screen_functions["Display Saved Games"] = Display_Saved_Games
	screen_functions["Display Save Slots"] = Display_Save_Slots

	target_actions["Next Screen"] = Next_Screen   //Move to the next screen
	target_actions["Select Slot"] = Select_Slot   //choose a save slot for the game
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
	//load the save slots
	saves, err := os.ReadDir(save_dir)
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range saves {
		save_slots[i] = s.Name()
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
	for _, line := range screen.Text {

		fmt.Println(line)
	}
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
	var index = 2 // I want the index to start one higher than the number of choices
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

func Display_Saved_Games() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	var index = 2 // I want the index to start at two every time
	for _, save_file := range save_slots {
		game_save := make(map[string]string)
		file, _ := os.ReadFile(fmt.Sprintf("%v%v", save_dir, save_file))
		err := json.Unmarshal(file, &game_save)
		if err != nil {
			log.Printf("Cannot unmarshal the json %d\n", err)
		}
		_, ok := game_save["current_story"]
		if !ok {
			continue
		}
		sans_ext := strings.TrimSuffix(save_file, ".json")
		var new_choice = Choice{
			Id:     strconv.Itoa(index),
			Text:   fmt.Sprintf("%v. %v : %v - %v", strconv.Itoa(index), sans_ext, game_save["current_story"], game_save["save_time"]),
			Target: map[string]string{"Load Game": save_file},
		}
		current_screen.Choices[strconv.Itoa(index)] = new_choice
		index++
	}
}

func Display_Save_Slots() {
	var current_story = stories[game_state["current_story"]]
	var current_screen = current_story.Screens[game_state["position"]]
	var index = 2 // I want the index to start at two every time
	for _, save_file := range save_slots {
		game_save := make(map[string]string)
		file, _ := os.ReadFile(fmt.Sprintf("%v%v", save_dir, save_file))
		err := json.Unmarshal(file, &game_save)
		if err != nil {
			log.Printf("Cannot unmarshal the json %d\n", err)
		}
		sans_ext := strings.TrimSuffix(save_file, ".json")
		var new_choice = Choice{
			Id:     strconv.Itoa(index),
			Text:   fmt.Sprintf("%v. %v : %v - %v", strconv.Itoa(index), sans_ext, game_save["current_story"], game_save["save_time"]),
			Target: map[string]string{"Select Slot": save_file},
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
	file, _ := os.ReadFile(save_dir + name)
	err := json.Unmarshal(file, &game_state)
	if err != nil {
		log.Printf("Cannot unmarshal the json %d\n", err)
	}
}

func Select_Slot(slot string) {
	game_state["save_slot"] = slot //set the save slot in the game state
	Next_Screen("2")               // call the next screen for choosing a story

}

func Save_Game(name string) {
	// using date time stamps for save game names for now
	t := time.Now().Format(time.StampMilli)
	t = strings.Replace(t, ":", "-", -1) // need to replace : with - so the name is valid
	game_state["save_time"] = t
	content, err := json.Marshal(game_state)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf(save_dir+"%v", game_state["save_slot"]), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Exit_Game(name string) {

	CallClear()
	fmt.Println(Reset + "Thank you for playing! :)")
	// I don't want to cause a panic...
	os.Exit(0)
}
