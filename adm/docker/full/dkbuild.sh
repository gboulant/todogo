#!/bin/sh

iname="todogo"

root=$(dirname $0)
contextpath=$root/dockerfiles
SYSTEM=ubuntu
HTTP_PROXY=""

options="--build-arg SYSTEM=${SYSTEM}"
options="--build-arg HTTP_PROXY=${HTTP_PROXY} $options"

# Creation of the docker image
dockerfile=$contextpath/Dockerfile
docker build -f $dockerfile $options -t $iname $contextpath

# Creation of a volume for persistence of the home directory
vname=$iname.home
docker volume create $vname --label cfolder=/home --label contains="users homes"
