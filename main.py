
import click
import time
import shutil
import sys

import termios
import struct
import fcntl
import ipdb
from gamestate import GameState
from loader import Loader
from pprint import pprint

GAMESTATE = GameState

def game_loop():

    content = load()
    render(content)
    while True:
        try:
            game_tick()
        except KeyboardInterrupt:
            sys.exit()

def game_tick():
    pass

def render(content):
    for line in content['splash'].readlines():
            time.sleep(0.25)
            print(line.replace('\n',''))


def load():

    content = {
        "splash": open('files/menu.txt','r'),
    }
    return content

if __name__ == '__main__':
    game_loop()
