#!/bin/sh

cfgrootdir=$(todo config -i | grep "Configuration root directory:" | cut -d":" -f2)
vi $cfgrootdir/config.json
