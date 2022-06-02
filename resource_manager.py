import os
import json
from glob import glob

from loader import Loader
from event_manager import events, Event


class Initialisation(Event):
    pass

class ResourceManager():
    current_script = ""
    resources_dir = "./files"
    resource_files = []
    resources = {}

    def __init__(self):
        self.scan()
        for resource in self.resource_files:
            self.load_resource(resource)
        current_script = "Main Menu"

    def scan(self):
        for root, dirs, files in os.walk(self.resources_dir):
            for file in files:
                self.resource_files.append(os.path.join(root, file))


    def load_resource(self, resource_path):
        resource = Resource()
        with open(resource_path, 'r') as resource_file:
            data = json.load(resource_file)
            for k,v in data.items():
                setattr(resource, k, v)
        self.resources[resource.resource_title] = resource

    def fetch_script(self, target):

        return self.resources.get(target, None)




class Resource():


    def __init__(self):
        pass
