
import click
import time
import shutil
import sys

from director import Director

def game():
    director = Director()
    director.game_loop()

if __name__ == '__main__':
    game()
