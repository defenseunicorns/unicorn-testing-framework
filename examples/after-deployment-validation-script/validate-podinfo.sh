while [[ $(kubectl get pods -n podinfo -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do
  echo "waiting for pod" && sleep 1
done
kubectl get all -n podinfo
