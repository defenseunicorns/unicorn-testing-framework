package main

import (
	"bufio"
	"net/http"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// need kubeconfig
func TestGamesAreRunning(t *testing.T) {
	t.Logf("Running test: %s", t.Name())
	// Start the process in the background
	cmd := exec.Command("zarf", "connect", "games", "--cli-only")

	zarfConnectGamesStdOut, err := cmd.StdoutPipe()
	assert.NoError(t, err, "Error creating stdout pipe")

	err = cmd.Start()
	assert.NoError(t, err, "Error starting command")

	// kill the process when the test exits
	defer func() {
		err = cmd.Process.Kill()
		assert.NoError(t, err, "Error waiting for command")
	}()

	// read the console output to get the port-forwarded address
	reader := bufio.NewReader(zarfConnectGamesStdOut)
	content := make([]byte, 100)
	count, err := reader.Read(content)
	// line, err := reader.ReadString('\n')
	assert.NoError(t, err, "Error reading from stdout")
		
	url := string(content[:count])

	// log the address
	t.Logf("The address we received is: %s\n", url)

	// test the address
	resp, err := http.Get(url)
	assert.NoError(t, err, "There shouldn't be an error")
	assert.Equal(t, resp.StatusCode, 200)
}
