#!/bin/bash

while true; do name=$(openssl rand -hex 16); dd if=/dev/zero of="${name}.txt" count=1000 bs=1000000; done;
