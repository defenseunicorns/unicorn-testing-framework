for i in $(seq 1 14); do
  echo "waiting for pod..."
  if [[ $(kubectl get pods -n podinfo -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') = "True" ]]; do
    echo "found it!"
    kubectl get all -n podinfo
    exit 0
  fi
  sleep 5
done
echo "timed out waiting for pod"
exit 1
