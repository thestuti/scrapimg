#!/bin/bash

drive=$(echo $1 | cut -d ":" -f 1)

path=$(echo $1 | sed "s/.*://; s/\\\/\//g")

docker run -v ${PWD}/${drive}:${path} extract-images -l ${path}
