
events = {}



class Event():
    def __init__(self):
        pass

    
def Subscribe( event_class, handler):

    if event_class not in events:
        events[event_class] = [handler]
    else:
        events[event_class].append(handler)

async def notify(event: Event):
    for handler in events.get(event.__class__, ()):
        handler()



class Handler():
    def __call__(self):
        pass
