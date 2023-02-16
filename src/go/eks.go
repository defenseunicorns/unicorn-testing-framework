package utf

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func setupEKS(t *testing.T, platform EKSPlatform) {
	t.Parallel()

	clusterName := fmt.Sprintf("%s-%s", platform.InstanceName, random.UniqueId())
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: platform.TestFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"instance_name": clusterName,
			"region":        platform.Region,
		},
	})
	terraform.InitAndApply(t, terraformOptions)
}

func teardownEKS(t *testing.T, platform EKSPlatform) {
	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)
	terraform.Destroy(t, terraformOptions)
}
