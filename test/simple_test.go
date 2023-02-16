package test

import (
	"os"
	"testing"

	utf "github.com/defenseunicorns/unicorn-testing-framework/src/go"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	// ### SETUP CODE THAT MUST BE AT THE BEGINNING OF EVERY TEST ###
	t.Parallel()
	setupScript, err := os.ReadFile("setupScript.sh")
	require.NoError(t, err)
	platform := utf.NewEC2Platform(t, "andy-and-kiran-funtimes", "us-east-1", "t2.micro", string(setupScript))
	defer utf.Teardown(t, platform)
	utf.Setup(t, platform)
	// ### END SETUP CODE ###

	// We now have a running EC2 instance with the setup script run on it. We can now run tests to validate that our deployment that we ran in the setup script has resulted in the outcomes that we expect.
	teststructure.RunTestStage(t, "TEST", func() {
		// Wait up to 30 seconds for google.com to be reachable (should happen immediately for google, but is a good example for testing an actual endpoint in a new cluster)
		output, err := platform.RunSSHCommandWithOptionalSudo(t, `timeout 30 bash -c "while ! curl -L -s --fail --show-error https://google.com > /dev/null; do sleep 5; done"`, false)
		require.NoError(t, err, output)
	})
}
