#!/bin/bash

docker rm apptwo
docker rmi secondapp:latest
docker build . -f Dockerfile-buildit -t secondapp:latest
