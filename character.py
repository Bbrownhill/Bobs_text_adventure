import random
import re


reg = '([dD])(\d+)'


class DieToLarge(Exception):
    pass


class TooManyDie(Exception):
    pass

class Character():


    roll_methods = {
        "PointBuy": "",
        "StandardArray": "",
        "Rolled": rolled_method,
    }

    attributes = {}

    def init(self):
        pass

    def assign_race():
        pass

    def assign_class():
        pass

    def assign_character_details():
        pass
    
    def roll_stats(self, method):
        stat_rolls = roll_methods[method]




    def rolled_method():
        stat_rolls = []
        for stat in range(6):
            initial_rolls = [random.randint(1, 6) for x in range(4)]
            initial_rolls.sort()
            initial_rolls.pop(0)
            stat_rolls.append(sum(initial_rolls))
        return stat_rolls
