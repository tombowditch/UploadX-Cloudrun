#!/bin/bash

SCREENCAPTURE="/usr/sbin/screencapture"
UPLOADX_URL="https://example.com"
UPLOADX_KEY="xxxxxx"
IMG_NAME="/Users/tom/code/screeny/screenshot.png"

echo "*"
echo "* screeny"
echo "*"
echo ""

echo "* clearing clipboard..."
echo "<img still uploading>" | pbcopy

echo "* removing any old screenshot"
rm "$IMG_NAME"

echo "* running screencapture"
$SCREENCAPTURE -o -i "$IMG_NAME"

echo "* screencapture ran, checking if $IMG_NAME exists"

if [ -e "$IMG_NAME" ]; then
    echo "* file exists, uploading..."
    jsonresult=$(curl -s -F "img=@${IMG_NAME}" -F "key=${UPLOADX_KEY}" ${UPLOADX_URL}/upload)
    result=$(echo "$jsonresult" | /usr/local/bin/jq -r ".Name")
    jres=$(echo '$jsonresult' | /usr/local/bin/jq -r ".Name")

    if [ $result = "null" ]; then
	osascript -e 'display notification "File not uploaded :-(" with title "screeny"'
	echo "* error: $jsonresult"
    else
	echo "* successfully uploaded: $jsonresult"
	echo "$jsonresult" >> ~/.screeny.log
	echo "${UPLOADX_URL}/${result}" | pbcopy
	osascript -e 'display notification "File uploaded!" with title "screeny"'
    fi

    echo "* removing $IMG_NAME"
    rm "$IMG_NAME"
else 
    echo "* file does not exist, exiting..."
    osascript -e 'display notification "File not found :-(" with title "screeny"'
fi 

