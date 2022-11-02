#!/bin/sh

cfgrootdir=$(todo config -i | grep "Configuration root directory:" | cut -d":" -f2)
gititems=$(ls $cfgrootdir)

echo "=============================================="
echo "todo-sync: pull"
git -C $cfgrootdir pull

echo "=============================================="
echo "todo-sync: add and commit"
git -C $cfgrootdir add $gititems
git -C $cfgrootdir commit -m "todo database commit from $(whoami)@$(hostname)"

echo "=============================================="
echo "todo-sync: push"
git -C $cfgrootdir push 
