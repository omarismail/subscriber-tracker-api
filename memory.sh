#!/bin/bash


A="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
for power in $(seq 8); do
  A="${A}${A}"
done
