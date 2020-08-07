package tnf

import (
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
    "fmt"
    "regexp"
    "encoding/base64"
    "crypto/rand"
    "strings"
)

// A TCP test implemented using command line tool `nc`.
type Tcp struct {
    result     int
    // config
    timeout    int
    bindHost   string
    serverHost string
    // state
    state      int
    // state variables
    serverId   int
    clientId   int
    serverPort string
    magic      string
}

// Test states
const (
    invalid_state = iota
    wait_server_args
    wait_server_up
    wait_client_args
    wait_client_tx
    wait_server_rx
    wait_client_eof
    wait_server_eof
)

const listening string = `listening on \[.*\] (\d+) \.\.\.`

// Return the command line args for either the server or the client.
func (test *Tcp) Args(arg interface{}) ([]string, error) {
    context, ok := arg.(*Context)
    if !ok {
        return nil, &NoContextError{test}
    }
    switch test.state {
    case wait_server_args:
        test.serverId = context.id
        return ServerCmd(test.bindHost), nil
    case wait_client_args:
        return ClientCmd(test.serverHost, test.serverPort), nil
    default:
        reason := fmt.Sprintf("%T in bad state %d for Args", test, test.state)
        return nil, &ProgrammingError{reason}
    }
}

// Return the timeout in seconds for the test.
func (test *Tcp) Timeout() int {
    return test.timeout
}

// Return the test result.
func (test *Tcp) Result() int {
    return test.result
}

// Return step which expects the port number of the TCP server.
func (test *Tcp) ReelFirst(arg interface{}) *reel.Step {
    test.state = wait_server_up
    return &reel.Step{
        Expect:  []string{listening},
        Timeout: test.timeout,
    }
}

// On match, drive the FSM.
// When the port number of the TCP server is received, the TCP client is started
// and a step is returned to send a random magic string to the server.
// When the random magic string has been sent by the client, a step is returned
// to receive the same string at the server.
// When the random magic string has been received by the server, a step is
// returned to close the connection from the client.
func (test *Tcp) ReelMatch(pattern string, before string, match string, arg interface{}) *reel.Step {
    context := arg.(*Context)
    switch test.state {
    case wait_server_up:
        re := regexp.MustCompile(listening)
        matched := re.FindStringSubmatch(match)
        if matched != nil {
            test.serverPort = matched[1]
            test.state = wait_client_args
            id, err := context.executor.Start(test)
            if err == nil {
                test.state = wait_client_tx
                test.clientId = id
                context.id = test.clientId
                magic := make([]byte, 4)
                rand.Read(magic)
                test.magic = base64.StdEncoding.EncodeToString(magic)
                return &reel.Step{
                    Execute: test.magic,
                    Expect: []string{test.magic},
                }
            } else {
                context.err = err
            }
        }
    case wait_client_tx:
        if match == test.magic {
            test.state = wait_server_rx
            context.id = test.serverId
            return &reel.Step{
                Expect: []string{test.magic},
                Timeout: test.timeout,
            }
        }
    case wait_server_rx:
        if match == test.magic {
            test.state = wait_client_eof
            context.id = test.clientId
            return &reel.Step{
                Execute: reel.CTRL_D,
                Expect: []string{reel.CTRL_D},
                Timeout: test.timeout,
            }
        }
    }
    return nil
}

// On timeout, take no action.
func (test *Tcp) ReelTimeout(arg interface{}) *reel.Step {
    return nil
}

// On eof, drive the FSM.
// When the client has closed, close the server.
// When the server has closed, the test result is success.
func (test *Tcp) ReelEof(arg interface{}) *reel.Step {
    context := arg.(*Context)
    switch test.state {
    case wait_client_eof:
        test.state = wait_server_eof
        context.id = test.serverId
        return &reel.Step{
            Expect: []string{reel.CTRL_D},
            Timeout: test.timeout,
        }
    case wait_server_eof:
        test.result = SUCCESS
    }
    return nil
}

// Return command line args for running a TCP server bound to `host` on an
// arbitrary port number.
func ServerCmd(host string) []string {
    server := strings.Join([]string{"nc", "-s", host, "-l", "-v", "2>&1"}, " ") // does not work with OpenBSD netcat
    return []string{"sh", "-c", server}
}

// Return command line args for running a TCP client to connect to `host`:`port`
// and which exits when its stdin is closed.
func ClientCmd(host string, port string) []string {
    return []string{"nc", "-q", "0", host, port}
}

// Create a new `Tcp` test which:
// * runs a TCP server bound to `bindHost` and an arbitrary port number;
// * connects a TCP client to that server (to `serverHost`);
// * generates and send a random magic string from client to server; and,
// * ensures the random magic string is received by the server.
// This complex test can only be run under `Executor`: error otherwise.
func NewTcp(timeout int, bindHost string, serverHost string) *Tcp {
    return &Tcp{
        result: ERROR,
        timeout: timeout,
        bindHost: bindHost,
        serverHost: serverHost,
        state: wait_server_args,
    }
}
