#!/usr/bin/env bash

kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=podinfo  -n podinfo --timeout=300s
kubectl get all -n podinfo
