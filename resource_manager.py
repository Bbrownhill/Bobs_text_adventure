import os
import json

from loader import Loader

class ResourceManager():
    current_script
    resources_dir = "./files"
    resource_files = []
    resources = {}

    def __init__(self):
        scan()
        for resource in self.resource_files():
            load_resource(resource_path)
        current_script = "Main Menu"

    def scan(self):
        self.resource_files = os.listdir(self.resources_dir)

    def load_resource(self, resource_path):
        resource = Resource()
        with open(resource_path) as resource_file:
            data = json.load(resource_file)
            for k,v in data.items():
                setattr(resource, k, v)
        self.resources[resource.resource_title] = resource






class Resource():


    def __init__(self):
        pass
