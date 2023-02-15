package test

import (
	"os"
	"testing"

	utf "github.com/defenseunicorns/unicorn-testing-framework/src/go"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	t.Parallel()
	setupScript, err := os.ReadFile("setupScript.sh")
	require.NoError(t, err)
	options := utf.NewEC2Options(t, "andy-and-kiran-funtimes", "us-east-1", "t2.micro", string(setupScript))
	defer utf.Teardown(t, options)
	utf.Setup(t, options)
	teststructure.RunTestStage(t, "TEST", func() {
		// My test code here
	})
}
