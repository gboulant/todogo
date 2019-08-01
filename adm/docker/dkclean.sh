#!/bin/sh

iname="todogo"
vname="$iname.home"
cname="$iname.ctnr"

docker container rm $cname
docker volume rm $vname
#docker image rm $iname
