package test

import (
	"fmt"
	"testing"

	utf "github.com/defenseunicorns/unicorn-testing-framework/src/go"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestEKS(t *testing.T) {
	t.Parallel()
	platform := utf.NewEKSPlatform(t, "kiran", "us-east-1")
	defer utf.Teardown(t, platform)
	utf.Setup(t, platform)

	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)

	fmt.Println(terraform.Output(t, terraformOptions, "cluster_endpoint"))

}
