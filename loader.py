import os
import json

class Loader():

    resource_content = {}

    def __init__(self):
        load(scan())

    def scan(self):
        files = os.listdir('files')
        return files



    def load(self,files):
        for file in files:
            with open(file) as content:
                data = json.load(content)
                resource_content['screen_title'] = data['screen_lines']
