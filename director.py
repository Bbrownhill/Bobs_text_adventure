import os
import json

from loader import Loader
from resource_manager import ResourceManager, Resource
from event_manager import events, Event

from pprint import pprint


class Director():

    gamestate = "Initalising"
    save_dir = "./saves"
    resource_manager = ResourceManager()
    functions = {}

    def __init__(self):
        self.functions = {
            "Exit": self.exit,
            "Save": self.save_game,
            "Save and exit": self.save_and_exit,
            "New": self.new_game,
            "Load": "placeholder",
            "Conf": "placeholder"
        }

    def game_loop(self):
        print(self.gamestate)
        self.resource_manager.current_script = self.load_menu("Main Menu")
        print("Main menu loaded")
        self.gamestate = "Running"
        print(self.gamestate)
        # LOOOOPS!
        while self.gamestate is not "Exiting":
            self.execute_game_script()
        else:
            self.exit()

    def load_menu(self, target):
        print("Loading Main Menu")
        return self.resource_manager.fetch_script(target)



    def execute_game_script(self):
        os.system('clear')
        screen = self.fetch_screen(self.resource_manager.current_script.Current_Screen)
        for line in screen["Screen_lines"]:
            print(line)
        for option in screen["Options"]:
            print(option["Text"])
        user_input = input()

        selection = self.select_screen(screen, user_input)
        if selection != None:
            self.resource_manager.current_script.Previous_Screen = self.resource_manager.current_script.Current_Screen
            self.resource_manager.current_script.Current_Screen = selection
        self.parse_selection(selection)


    def select_screen(self, screen, input):
        new_target = None
        for option in screen["Options"]:
            if option["Option"] == input:
                return option["Links To"]

    def parse_selection(self, value):
        for k, v in self.functions.items():
            if k == value:
                self.functions[value]()
        return value


    def fetch_screen(self, target):
        for screen in self.resource_manager.current_script.Screens:
            if screen["Screen_id"] == target:
                return screen

    def new_game(self):
        game_script_path = "files/scripts/Test_story_1/"
        self.resource_manager.current_script = self.resource_manager.fetch_script("Collisions")

    def save_game(self):
        print("Enter a save file name:")
        save_file_name = input()
        save_file_content = {
                             "script" : self.resource_manager.current_script.Resource_title,
                             "screen" : self.resource_manager.current_script.Previous_Screen
        }
        with open(f'{self.save_dir}/{save_file_name}.json', "w") as save_file:
            json.dump(save_file_content, save_file)



    def save_and_exit(self):
        self.save_game()
        self.exit()

    def load_game(self):
        pass

    def exit(self):
        self.gamestate = "Exiting"
        print('Exiting\nThank you for playing!')
