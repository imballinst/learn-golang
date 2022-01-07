#!/bin/bash

action=$1
# This second argument is replacable.
# By default, this refers to Hades' uuid.
uuid="${2:-4c0ba5d1-8139-4506-9334-08a8c3314c0d}"

if [ $action = "" ]; then
  echo "Action not provided. Available ones: 'list', 'create', 'get"
  exit 1
fi

case $action in
  "list")
    curl localhost:8080/characters
    ;;
  "create")
    curl http://localhost:8080/characters \
      --include \
      --header "Content-Type: application/json" \
      --request "POST" \
      --data '{"name": "Themis", "role": "Elidibus", "level": 99}'
    ;;
  "get")
    curl localhost:8080/characters/$uuid
    ;;
  "update")
    curl http://localhost:8080/characters/$uuid \
      --include \
      --header "Content-Type: application/json" \
      --request "PUT" \
      --data '{"id": "4c0ba5d1-8139-4506-9334-08a8c3314c0d", "name": "Haydes", "role": "Emet-Selch", "level": 99},'
    ;;
esac
