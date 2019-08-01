#!/bin/sh
iname="todogo" # image name
cname="$iname.ctnr"  # container name
vname="$iname.home" # volume name

options="-v $vname:/home"
docker run --rm $options --name $cname -it $iname /home/admin/todogo/bin/todo $*



