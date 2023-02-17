package utf

import (
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
)

type EC2Platform struct {
	InstanceName string
	Region       string
	InstanceType string
	TestFolder   string
	SetupScript  string
}

type EKSPlatform struct {
	InstanceName string
	Region       string
	TestFolder   string
}

func NewEC2Platform(t *testing.T, instanceName string, region string, instanceType string, setupScript string) EC2Platform {
	platform := EC2Platform{
		InstanceName: instanceName,
		Region:       region,
		InstanceType: instanceType,
		SetupScript:  setupScript,
	}
	platform.TestFolder = teststructure.CopyTerraformFolderToTemp(t, "../..", "src/tf/public-ec2-instance")

	return platform
}

func NewEKSPlatform(t *testing.T, instanceName, region string) EKSPlatform {
	platform := EKSPlatform{
		InstanceName: instanceName,
		Region:       region,
	}
	platform.TestFolder = teststructure.CopyTerraformFolderToTemp(t, "../..", "src/tf/eks")

	return platform
}
