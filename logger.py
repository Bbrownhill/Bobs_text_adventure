import os
from datetime import datetime


def log(message):
    now = datetime.now()
    date_time = now.strftime("%Y-%m-%dT%H:%M:%S")
    with open('Logs/log.txt', 'a') as log_file:
        log_file.write(f'{date_time}, {message}')
