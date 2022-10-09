#! /bin/bash

if [ ! -d scripts/ ]
then
  printf "Please run script from project root directory."
fi

docker build -t crypto-app .