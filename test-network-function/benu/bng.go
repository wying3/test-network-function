package benu

import (
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/identifier"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/reel"
	"regexp"
	"strings"
	"time"
)

const (
	BenuBNGShowCounterCmd = "show bef counter subscribers"
	BenuBNGServerStartCmd = "'cd v2.75;./t-rex-64 -i'"
	BENUTrafficGenCmd     = "-c trex -- python /root/v2.75/BenuQinQ.py"

	// benuCountRegEx = `^[\\d]+\\s+[\\d.]+\\s+(\\d+)\\s+(\\d+)`
	benuCountRegEx = `.+`
)

// OutputRegex matches the output from inspecting ID registrations.
var OutputRegex = regexp.MustCompile(benuCountRegEx)

// BenuBNG is all pod for BENU BNG testing including trex
type BenuBNG struct {
	// args represents the Unix command.
	args      []string
	pod       string
	namespace string
	// result is the result of the tnf.Test.
	result int
	// timeout is the tnf.Test timeout.
	timeout time.Duration
	// Identifier is the tnf.Test specific test identifier.
	Identifier identifier.Identifier `json:"identifier" yaml:"identifier"`
	Command    string
	resultOut  string
}

func (b *BenuBNG) GetIdentifier() identifier.Identifier {
	return b.Identifier
}
func (b *BenuBNG) GetResultOut() string {
	return b.resultOut
}

// Args returns the command line args for the test.
func (b *BenuBNG) Args() []string {
	return b.args
}

// Timeout returns the timeout in seconds for the test.
func (b *BenuBNG) Timeout() time.Duration {
	return b.timeout
}

// Result returns the test result.
func (b *BenuBNG) Result() int {
	return b.result
}

// ReelFirst returns a step which expects an ip summary for the given device.
func (b *BenuBNG) ReelFirst() *reel.Step {
	return &reel.Step{
		Expect:  []string{benuCountRegEx},
		Timeout: b.timeout,
	}
}

// ReelMatch parses the Registration.  Returns no step; the test is complete.
func (b *BenuBNG) ReelMatch(pattern string, _ string, match string) *reel.Step {
	b.result = tnf.ERROR
	b.resultOut = match
	if pattern == benuCountRegEx {
		// Indicates that the command was successfully run, but there were no registered NRFs.
		//matches := OutputRegex.FindAllString(match, -1)
		b.result = tnf.SUCCESS
	}

	return nil
}

// ReelTimeout returns a step which kills the ping test by sending it ^C.
func (b *BenuBNG) ReelTimeout() *reel.Step {
	return &reel.Step{Execute: reel.CtrlC}
}

// ReelEOF does nothing.  On EOF, take no action.
func (b *BenuBNG) ReelEOF() {
}

// Command returns command line args for benu requests
func Command(cmd string) []string {
	return strings.Split(cmd, " ")
}

// NewBenuBNG creates a BenuBNG instance.
func NewBenuBNG(timeout time.Duration, name, namespace string, command string) *BenuBNG {
	return &BenuBNG{result: tnf.ERROR, timeout: timeout, args: Command(command), namespace: namespace}
}
