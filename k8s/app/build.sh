#! /bin/bash

IMAGE_NAME="docker.io/lvsoso/greeter-app:v0.1"

docker build . -t  "$IMAGE_NAME"

docker push "$IMAGE_NAME"
