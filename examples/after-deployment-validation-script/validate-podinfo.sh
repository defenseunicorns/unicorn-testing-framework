#!/usr/bin/env bash

#eqivilent to: zarf tools wait-for pod app.kubernetes.io/name=podinfo ready -n podinfo
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=podinfo -n podinfo --timeout=300s
kubectl get all -n podinfo
