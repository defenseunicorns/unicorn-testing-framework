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
		TerraformDir: platform.TestFolder,
		Vars: map[string]interface{}{
			"instance_name": clusterName,
			"region":        platform.Region,
		},
	})
	teststructure.SaveTerraformOptions(t, terraformOptions.TerraformDir, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
}

func teardownEKS(t *testing.T, platform EKSPlatform) {
	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)
	terraform.Destroy(t, terraformOptions)
}
