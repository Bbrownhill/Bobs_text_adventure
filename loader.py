import os
import json

class Loader():

    game_content = {}

    def init():
        load(scan())

    def scan():
        files = os.listdir('files')
        return files



    def load(files):
        for file in files:
            with open(file) as content:
                data = json.load(content)
                game_content['screen_title'] = data['screen_lines']
