#!/bin/python3

import os
import json

TEST_METADATA = 'flourish_test_metadata.json'

class DirectoryInfo:
    def __init__(self, directoryName, totalCount, individualCount, groupCount, zeroParticipantCount):
        self.directoryName = directoryName
        self.totalCount = totalCount
        self.individualCount = individualCount
        self.groupCount = groupCount
        self.zeroParticipantCount = zeroParticipantCount

    def __str__(self):
        return f"name: {self.directoryName}\ntotalCount: {self.totalCount}, individualCount: {self.individualCount}, groupCount: {self.groupCount}, zeroParticipantCount: {self.zeroParticipantCount}"


def testFlourish(username, convoDir, msgFileType):
    print(convoDir)

    (root, dirs, files) = os.walk(convoDir).__next__()

    totalConvoCount = len(dirs)
    individualChatCount = 0
    groupChatCount = 0
    zeroParticipantsChatCount = 0

    for d in dirs:
        dirPath = os.path.join(root, d)
        (_, _, files) = os.walk(dirPath).__next__()

        jsonFiles = (x for x in files if x.lower().endswith(msgFileType))
        fileList = []
        for f in jsonFiles:
            fileList.append(os.path.join(dirPath, f))

        # TODO: handle multiple message json files in same convo
        # if len(fileList) != 1:
        #     print('Error finding messageFile. Exiting.')
        #     return errorList
        messageFile = fileList[0]

        with open(messageFile) as f:
            try:
                d = json.load(f)
                participants = d['participants']
            except Exception as e:
                print('Error reading message file: ' + messageFile + '. Exception: ' + str(e) + '. Exiting.')
                return DirectoryInfo(None, None, None, None, None)

        if len(participants) == 1 or len(participants) == 2:
            individualChatCount += 1
        elif len(participants) > 2:
            groupChatCount += 1
        else:
            zeroParticipantsChatCount += 1

    return DirectoryInfo(convoDir, totalConvoCount, individualChatCount, groupChatCount, zeroParticipantsChatCount)

def main():

    # fetch metadata
    if os.path.exists(TEST_METADATA):
        with open(TEST_METADATA) as jsonFile:
            try:
                data = json.load(jsonFile)
                infoFile = data['infoFile']
                rootDir = data['rootDir']
                convoDirNames = data['convoDirNames']
                msgFileType = data['msgFileType']
                outFile = data['outFile']

                infoFilePath = os.path.join(rootDir, infoFile)
            except Exception as e:
                print('Error reading file: ' + TEST_METADATA + '. Exception: ' + str(e) + '. Exiting.')
                return
    else:
        print('File ' + TEST_METADATA + ' could not be found. Exiting.')
        return


    # get user's name
    if os.path.exists(infoFilePath):
        with open(infoFilePath) as info:
            try:
                data = json.load(info)
                username = data[list(data.keys())[0]]['FULL_NAME'][0]
            except Exception as e:
                print('Error finding username from ' + infoFilePath + '. Exception: ' + str(e) + '. Exiting.')
                return
    else:
        print('File ' + infoFilePath + ' could not be found. Exiting.')
        return

    # iterate through all conversations
    for convoDir in convoDirNames:

        with open(outFile, 'a') as f:
            try:
                absoluteConvoDir = os.path.join(rootDir, convoDir)
                directoryInfo = testFlourish(username, absoluteConvoDir, msgFileType)
                print(str(directoryInfo))

                # output results
                f.write('%s\n\n' % str(directoryInfo))
            except Exception as e:
                print('Error writing to file: ' + outFile + '. Exception: ' + str(e) + '. Exiting.')
                return


if __name__ == '__main__':
    main()
