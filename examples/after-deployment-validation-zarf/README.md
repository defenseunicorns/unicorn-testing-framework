# After Deployment Validation

This example shows how Zarf component actions can be used to validate a deployment against the cluster using Kubectl.
We use the following `actions` block to run a script that validates `Podinfo` has been
successfully deployed.

``` yaml
  actions:
    onDeploy:
      after:
      - cmd: ./validate-podinfo.sh
```