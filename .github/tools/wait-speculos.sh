#!/bin/bash

counter=0
while true; do
  docker logs --tail=5 speculos
  docker logs --tail=5 speculos 2>&1 | grep using\ SDK\ version\ 1.6
  if [ $? == 0 ]; then
    exit 0
  fi
  if [[ "$counter" -gt 60 ]]; then
    echo "Counter: $counter times reached; Exiting loop!"
    exit 1
  fi
  counter=$((counter+1))
  sleep 1
done
