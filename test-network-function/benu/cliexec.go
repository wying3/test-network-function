package benu

import (
	expect "github.com/google/goexpect"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	"io"
	"time"
)

// SpawnShell creates an interactive shell subprocess based on the value of $SHELL, spawning the appropriate underlying
// PTY.
func SpawnCLIExec(spawner *interactive.Spawner, timeout time.Duration, in *io.WriteCloser, out *io.Reader, opts ...expect.Option) (*interactive.Context, error) {
	command := "/opt/benu/bin/cliexec"
	command = "hostname"
	var args []string
	return (*spawner).Spawn(command, args, timeout, in, out, opts...)
}
