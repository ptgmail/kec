#!/bin/bash

docker rm apptwo
docker run -d -p 8081:3000/tcp --add-host host.docker.internal:host-gateway --net kec3 --name apptwo secondapp:latest
sleep 2
curl http://localhost:8081
