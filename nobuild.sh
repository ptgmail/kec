#!/bin/bash

go build
docker stop apptwo
docker rm apptwo
docker rmi secondapp:latest
docker build . -f Dockerfile -t secondapp:latest
