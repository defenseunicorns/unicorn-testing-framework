package games_package_test

import (
	"bytes"
	utf "github.com/defenseunicorns/unicorn-testing-framework/src/go"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"text/template"
)

func TestSimple(t *testing.T) {
	// ### SETUP CODE THAT MUST BE AT THE BEGINNING OF EVERY TEST ###
	t.Parallel()

	setupScript, err := os.ReadFile("setupScript.sh.gotmpl")
	require.NoError(t, err)

	// If you need variables in the setup script, you can use this to get them
	gitRepoURL, err := utf.GetEnvVar("REPO_URL")
	require.NoError(t, err)
	gitBranch, err := utf.GetEnvVar("GIT_BRANCH")
	require.NoError(t, err)
	setupScriptTemplate, err := template.New("setupScript").Parse(string(setupScript))
	require.NoError(t, err)
	var setupScriptBuffer bytes.Buffer
	err = setupScriptTemplate.Execute(&setupScriptBuffer, map[string]string{
		"gitRepoURL": gitRepoURL,
		"gitBranch":  gitBranch,
	})
	require.NoError(t, err)

	platform := utf.NewEC2Platform(t, "andy-and-kiran-funtimes", "us-east-1", "m5a.large", setupScriptBuffer.String())
	defer utf.Teardown(t, platform)
	utf.Setup(t, platform)
	// ### END SETUP CODE ###

	// We now have a running EC2 instance with the setup script run on it. We can now run tests to validate that our deployment that we ran in the setup script has resulted in the outcomes that we expect.
	teststructure.RunTestStage(t, "TEST", func() {
		// Nothing to do, The package deployment already did the test.
	})
}
