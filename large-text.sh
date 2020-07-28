#!/bin/bash

while true
do
  name=$(uuidgen)
  dd if=/dev/zero of="${name}.txt" count=1000 bs=1000000
done
