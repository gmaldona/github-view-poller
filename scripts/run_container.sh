#! /bin/bash

docker stop github-view-app
docker rm github-view-app
docker run -it --name github-view-app github-view-app