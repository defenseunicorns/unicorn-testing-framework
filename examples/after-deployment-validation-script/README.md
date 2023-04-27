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

To include a script with the Zarf package we use the following `files` block, which bundles the shell script and unzips it to the working directory during the deploy.

``` yaml
  files:
  - source: ./validate-podinfo.sh
    target: ./validate-podinfo.sh
    executable: true
```