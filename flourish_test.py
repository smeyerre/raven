#!/bin/python3

import os
import json

TEST_FILE = 'flourish_test_metadata.json'
OUT_FILE = 'test_input'

def testFlourish(username, convoDir, msgFile):
    (root, dirs, files) = os.walk(convoDir).__next__()

    totalConvoCount = len(dirs)
    individualChatCount = 0
    groupChatCount = 0
    zeroParticipantsChatCount = 0
    errorList = [-1, -1, -1, -1]

    for d in dirs:
        dirPath = os.path.join(root, d)
        (_, _, files) = os.walk(dirPath).__next__()
        jsonFiles = (x for x in files if x == msgFile)
        fileList = []
        for f in jsonFiles:
            fileList.append(os.path.join(dirPath, f))

        if len(fileList) != 1:
            print('Error finding messageFile. Exiting.')
            return errorList
        messageFile = fileList[0]

        with open(messageFile) as f:
            try:
                d = json.load(f)
                participants = d['participants']
            except Exception as e:
                print('Error reading message file: ' + messageFile + '. Exception: ' + str(e) + '. Exiting.')
                return errorList

        if len(participants) == 1 or len(participants) == 2:
            individualChatCount += 1
        elif len(participants) > 2:
            groupChatCount += 1
        else:
            zeroParticipantsChatCount += 1

    return [totalConvoCount, individualChatCount, groupChatCount, zeroParticipantsChatCount]

def main():
    if os.path.exists(TEST_FILE):
        with open(TEST_FILE) as jsonFile:
            try:
                data = json.load(jsonFile)
                username = data['username']
                convoDir = data['convoDir']
                msgFile = data['msgFile']
            except Exception as e:
                print('Error reading file: ' + TEST_FILE + '. Exception: ' + str(e) + '. Exiting.')
                return
    else:
        print('File ' + TEST_FILE + ' could not be found. Exiting.')
        return

    values = testFlourish(username, convoDir, msgFile)
    print(values)

    with open(OUT_FILE, 'w') as f:
        try:
            for x in values:
                f.write('%s\n' % x)
        except Exception as e:
            print('Error writing to file: ' + OUT_FILE + '. Exception: ' + str(e) + '. Exiting.')
            return

if __name__ == '__main__':
    main()
