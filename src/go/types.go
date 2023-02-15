package utf

import (
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
)

type EC2Options struct {
	InstanceName string
	Region       string
	InstanceType string
	TestFolder   string
	SetupScript  string
}

func NewEC2Options(t *testing.T, instanceName string, region string, instanceType string, setupScript string) EC2Options {
	options := EC2Options{
		InstanceName: instanceName,
		Region:       region,
		InstanceType: instanceType,
		SetupScript:  setupScript,
	}
	options.TestFolder = teststructure.CopyTerraformFolderToTemp(t, "..", "src/tf/public-ec2-instance")
	return options
}
