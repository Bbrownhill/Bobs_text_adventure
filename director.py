from gamestate import GameState
from loader import Loader
from resource_manager import ResourceManager, Resource

from pprint import pprint


class Director():

    gamestate = GameState()
    resource_manager = ResourceManager()

    def init(self):
        pass

    def game_loop(self):
        if self.gamestate.state is  "Initalising":
            self.initialise_game()


        while self.gamestate.state is not "Exiting":
            self.execute_game_script(self.resource_manager.current_script)
        else:
            self.exit()

    def initialise_game(self):
        pprint('Initialising game')


    def execute_game_script(self, script):
        pprint('executing a script')
        pprint(script)


    def exit(self):
        pprint('exiting')
