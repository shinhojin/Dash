#! /bin/bash

echo "stop all container"
docker stop $(docker ps -a -q)
echo "remove all container"
docker rm $(docker ps -a -q)