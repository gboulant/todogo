#!/bin/sh

iname="todogo" # image name
cname="$iname.ctnr"  # container name
vname="$iname.home"

options="-v $vname:/home"
docker run --rm $options --name $cname -it $iname

