# podinfo Package Testing Example

This is an example zarf package that will deploy the podinfo helm chart and test deployment's ingress after.

## Quick Start

From this directory
1. (optional) deploy a k3d cluster: `../../utils/k3d/k3d_cluster.sh create`
2. init zarf: `zarf init ./zarf-init-<architecture>-<version matching whatever zarf version>.tar.zst --confirm`
3. If not using the k3d cluster config under k3d.yaml, confirm environment variables under [terratest/terratest.env](terratest/terratest.env)
4. zarf package create: `zarf package create --confirm`
5. zarf package deploy: `zarf package deploy zarf-package-podinfo-<architecture>.tar.zst --components podinfo-test --confirm`

On successful completion you should see similar output:

```sh
📦 PODINFO-TEST COMPONENT


  ✔  Copying 2 files
     podinfo-testing
     === RUN   TestKubernetesIngressCheck
     === PAUSE TestKubernetesIngressCheck
     === CONT  TestKubernetesIngressCheck
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 retry.go:91: Wait for service podinfo-testing to be provisioned.
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 client.go:42: Configuring Kubernetes client using config file /Users/user/.kube/config with context
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 node.go:33: Getting list of nodes from Kubernetes
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 client.go:42: Configuring Kubernetes client using config file /Users/user/.kube/config with context
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 service.go:86: Service is now available
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 retry.go:91: HTTP GET to URL http://localhost:8080
     TestKubernetesIngressCheck 2023-02-16T12:08:40-08:00 http_helper.go:59: Making an HTTP GET call to URL http://localhost:8080
     --- PASS: TestKubernetesIngressCheck (0.02s)
     PASS
  ✔  Completed command "set -a; source ./terratest.env; ...."
  ✔  Zarf deployment complete
```
