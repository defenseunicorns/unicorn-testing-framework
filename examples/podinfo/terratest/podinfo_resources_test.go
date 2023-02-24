package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

// An example of how to do more expanded verification of the Kubernetes resource config in examples/kubernetes-basic-example using Terratest.
func TestKubernetesIngressCheck(t *testing.T) {
	t.Parallel()

	//get environment variable for kubenetes namespace
	namespace := os.Getenv("K8S_NAMESPACE")
	servicename := os.Getenv("K8S_SERVICE_NAME")
	hostnameandport := os.Getenv("K8S_HOSTNAME_PORT")

	// - Current context of the kubectl config file
	options := k8s.NewKubectlOptions("", "", namespace)

	// This will wait up to 10 seconds for the service to become available, to ensure that we can access it.
	k8s.WaitUntilServiceAvailable(t, options, servicename, 10, 1*time.Second)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Test the ingress connection for up to 5 minutes. This will only fail if we timeout waiting for the service to return a 200
	// response.
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", hostnameandport),
		&tlsConfig,
		30,
		10*time.Second,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)
}
