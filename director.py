
from loader import Loader
from resource_manager import ResourceManager, Resource
from event_manager import events, Event

from pprint import pprint


class Director():

    gamestate = "Initalising"
    resource_manager = ResourceManager()


    def __init__(self):
        pass

    def game_loop(self):
        self.resource_manager.current_script = self.load_menu("Main Menu")
        gamestate = "Running"
        # LOOOOPS!
        while self.gamestate is not "Exiting":

            self.execute_game_script(self.resource_manager.current_script)
        else:
            self.exit()

    def load_menu(self, target):
        return self.resource_manager.fetch_script(target)


    def execute_game_script(self, script):
        for line in script.screen_lines:
            print(line)
        selection = input()
        print(selection)





    def exit(self):
        pprint('exiting')
