#!/bin/bash

case "$1" in
  create | destroy)
    k3d cluster "$1" --config ./k3d.yaml
    ;;
  *)
    echo "usage: k3d cluster (create||destroy) --config ./k3d.yaml"
    exit 1
    ;;
esac