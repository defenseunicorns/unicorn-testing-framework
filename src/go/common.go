package utf

import (
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// Setup sets up the environment where the deployment will happen.
func Setup(t *testing.T, options interface{}) {
	teststructure.RunTestStage(t, "SETUP", func() {
		switch options := options.(type) {
		case EC2Options:
			setupEC2(t, options)
		default:
			require.Fail(t, "Unknown options type")
		}
	})
}

// Teardown the deployment environment
func Teardown(t *testing.T, options interface{}) {
	teststructure.RunTestStage(t, "TEARDOWN", func() {
		switch options := options.(type) {
		case EC2Options:
			teardownEC2(t, options)
		default:
			require.Fail(t, "Unknown options type")
		}
	})
}
