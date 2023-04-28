#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

case "$1" in
  create | delete)
    k3d cluster "$1" --config "${SCRIPT_DIR}/k3d.yaml"
    ;;
  *)
    echo "usage: k3d cluster (create||delete) --config ./k3d.yaml"
    exit 1
    ;;
esac