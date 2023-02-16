package utf

import (
	"fmt"
	"os"
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// Setup sets up the environment where the deployment will happen.
func Setup(t *testing.T, platform interface{}) {
	teststructure.RunTestStage(t, "SETUP", func() {
		switch platform := platform.(type) {
		case EC2Platform:
			setupEC2(t, platform)
		default:
			require.Fail(t, "Unknown platform type")
		}
	})
}

// Teardown the deployment environment
func Teardown(t *testing.T, platform interface{}) {
	teststructure.RunTestStage(t, "TEARDOWN", func() {
		switch platform := platform.(type) {
		case EC2Platform:
			teardownEC2(t, platform)
		default:
			require.Fail(t, "Unknown platform type")
		}
	})
}

// GetEnvVar gets an environment variable, returning an error if it isn't found.
func GetEnvVar(varName string) (string, error) {
	val, present := os.LookupEnv(varName)
	if !present {
		return "", fmt.Errorf("expected env var %v not set", varName)
	}

	return val, nil
}
