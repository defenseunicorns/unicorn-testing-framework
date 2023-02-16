package utf

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func setupEC2(t *testing.T, platform EC2Platform) {
	randomInstanceName := fmt.Sprintf("%s-%s", platform.InstanceName, random.UniqueId())
	keyPair, err := aws.CreateAndImportEC2KeyPairE(t, platform.Region, randomInstanceName)
	require.NoError(t, err)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: platform.TestFolder,
		Vars: map[string]interface{}{
			"aws_region":    platform.Region,
			"name":          randomInstanceName,
			"key_pair_name": randomInstanceName,
			"instance_type": platform.InstanceType,
		},
	})
	teststructure.SaveTerraformOptions(t, terraformOptions.TerraformDir, terraformOptions)
	teststructure.SaveEc2KeyPair(t, terraformOptions.TerraformDir, keyPair)
	terraform.InitAndApply(t, terraformOptions)
	err = waitForInstanceReady(t, platform, 5*time.Second, 15) //nolint:gomnd
	require.NoError(t, err)
	err = copyFileToInstanceE(t, platform, 0777, "/home/ubuntu/setupScript.sh", platform.SetupScript)
	require.NoError(t, err)
	output, err := platform.RunSSHCommand(t, `/home/ubuntu/setupScript.sh`, true)
	require.NoError(t, err, output)
}

func teardownEC2(t *testing.T, platform EC2Platform) {
	keyPair := teststructure.LoadEc2KeyPair(t, platform.TestFolder)
	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)
	terraform.Destroy(t, terraformOptions)
	aws.DeleteEC2KeyPair(t, keyPair)
}

// RunSSHCommand runs the given command on the given host via SSH, and return the stdout and stderr as a string.
func (platform EC2Platform) RunSSHCommand(t *testing.T, command string, asSudo bool) (string, error) {
	precommand := "bash -c"
	if asSudo {
		precommand = fmt.Sprintf(`sudo %v`, precommand)
	}
	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)
	keyPair := teststructure.LoadEc2KeyPair(t, platform.TestFolder)
	host := ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "public_instance_ip"),
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}
	var output string
	var err error
	count := 0
	done := false
	// Try up to 3 times to do the command, to avoid "i/o timeout" errors which are transient
	for !done && count < 3 {
		count++
		output, err = ssh.CheckSshCommandE(t, host, fmt.Sprintf(`%v '%v'`, precommand, command))
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				// There was an error, but it was an i/o timeout, so wait a few seconds and try again
				logger.Default.Logf(t, "i/o timeout error, trying again")
				logger.Default.Logf(t, output)
				time.Sleep(3 * time.Second)
				continue
			} else {
				logger.Default.Logf(t, output)
				return "nil", fmt.Errorf("ssh command failed: %w", err)
			}
		}
		done = true
	}
	logger.Default.Logf(t, output)
	return output, nil
}

// waitForInstanceReady tries/retries a simple SSH command until it works successfully, meaning the server is ready to accept connections.
func waitForInstanceReady(t *testing.T, platform EC2Platform, timeBetweenRetries time.Duration, maxRetries int) error {
	t.Helper()
	_, err := retry.DoWithRetryE(t, "Wait for the instance to be ready", maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := platform.RunSSHCommand(t, "whoami", true)
		if err != nil {
			return "", fmt.Errorf("unknown error: %w", err)
		}

		return "", nil
	})
	if err != nil {
		return fmt.Errorf("error while waiting for instance to be ready: %w", err)
	}

	// Wait another 5 seconds because race conditions suck
	time.Sleep(5 * time.Second) //nolint:gomnd

	return nil
}

func copyFileToInstanceE(t *testing.T, platform EC2Platform, mode os.FileMode, remotePath string, contents string) error {
	terraformOptions := teststructure.LoadTerraformOptions(t, platform.TestFolder)
	keyPair := teststructure.LoadEc2KeyPair(t, platform.TestFolder)
	host := ssh.Host{
		Hostname:    terraform.Output(t, terraformOptions, "public_instance_ip"),
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}
	return ssh.ScpFileToE(t, host, mode, remotePath, contents)
}
