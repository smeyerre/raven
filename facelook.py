import sys
import os
import json

if __name__ == '__main__':

    fileName = sys.argv[1]
    with open(fileName, "r") as jsonFile:
        data = json.load(jsonFile)

        participants = data["participants"]
        messages = data["messages"]
        title = data["title"]
        is_still_participant = data["is_still_participant"]
        thread_type = data["thread_type"]
        thread_path = data["thread_path"]


