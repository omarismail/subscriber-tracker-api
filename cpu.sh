#!/bin/bash

fulload() { 
  dd if=/dev/zero of=/dev/null | dd if=/dev/zero of=/dev/null | dd if=/dev/zero of=/dev/null | dd if=/dev/zero of=/dev/null & 
}

while true
do
  echo "hello world"
  fulload
done

