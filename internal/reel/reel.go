// Run a target subprocess with programmatic control over interaction with it.
// Programmatic control uses a Read-Execute-Expect-Loop ("REEL") pattern.
// This pattern is implemented by `bin/reel.exp`, a general purpose expect script.
package reel

import (
    "bufio"
    "encoding/json"
    "io"
    "os/exec"
)

const CTRL_C string = "\003" // ^C
const CTRL_D string = "\004" // ^D

// A step is an instruction for a single REEL pass.
// To process a step, reel.exp first sends the `Execute` string to the target
// subprocess (if supplied); then it will block until the subprocess output to
// stdout matches one of the regular expressions in `Expect` (if any supplied).
// Supplying only the constant `CTRL_D` in `Expect` signals to reel.exp to block
// until EOF is detected on the subprocess output.
// A positive integer `Timeout` (seconds) prevents blocking forever.
// A step is sent to reel.exp as a JSON object in a single line of text.
type Step struct {
    Execute string   `json:"execute,omitempty"`
    Expect  []string `json:"expect,omitempty"`
    Timeout int      `json:"timeout,omitempty"`
}

// An event is a notification related to a single REEL pass.
// reel.exp may report an `Event` which is a "match" of a step `Expect` pattern;
// a "timeout", if there is no match within the specified period; or "eof",
// indicating that the subprocess has exited.
// reel.exp sends an event as a JSON object in a single line of text.
// Value constraints are documented in comments by field.
type Event struct {
    Event   string `json:"event"`             // only "match" or "timeout" or "eof"
    Idx     int    `json:"idx,omitempty"`     // only present when "event" is "match"
    Pattern string `json:"pattern,omitempty"` // only present when "event" is "match"
    Before  string `json:"before,omitempty"`  // only present when "event" is "match"
    Match   string `json:"match,omitempty"`   // only present when "event" is "match"
}

// TODO: document
type Notification struct {
    Token interface{}
    Event *Event
    Error error
}

// A Handler implements desired programmatic control:
// `ReelMatch` informs of a match event;
// `ReelTimeout` informs of a timeout event;
// `ReelEof` informs of the eof event.
// To provide a consistent interface, each function returns the next step to
// perform: returning nil indicates that there is no step to perform.
// `arg` is an arbitrary value supplied via the dispatcher of the event.
type Handler interface {
    ReelMatch(pattern string, before string, match string, arg interface{}) *Step
    ReelTimeout(arg interface{}) *Step
    ReelEof(arg interface{}) *Step
}

// A `Reel` instance allows interaction with a target subprocess.
type Reel struct {
    subp    *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
    scanner *bufio.Scanner
}

// Open the target subprocess for interaction.
func (reel *Reel) Open() error {
    return reel.subp.Start()
}

// Send `step` to be run.
func (reel *Reel) send(step *Step) error {
    msg, err := json.Marshal(step)
    if err == nil {
        reel.stdin.Write(msg)
        reel.stdin.Write([]byte("\n"))
    }
    return err
}

// Receive an event.
func (reel *Reel) recv() (*Event, error) {
    var event Event
    reel.scanner.Scan()
    err := json.Unmarshal(reel.scanner.Bytes(), &event)
    return &event, err
}

// Dispatch `event` via `handler` with `arg`, return the next step to perform.
func Dispatch(event *Event, handler Handler, arg interface{}) *Step {
    switch event.Event {
    case "match":
        return handler.ReelMatch(event.Pattern, event.Before, event.Match, arg)
    case "timeout":
        return handler.ReelTimeout(arg)
    case "eof":
        return handler.ReelEof(arg)
    default:
        return nil
    }
}

// Perform `step` then consequent steps in response to dispatching events via
// `handler` with `arg`.
// Return on first error, or when there is no next step to perform.
func (reel *Reel) Step(step *Step, handler Handler, arg interface{}) error {
    for step != nil {
        err := reel.send(step)
        if err != nil {
            return err
        }
        if len(step.Expect) == 0 {
            return nil
        }
        event, err := reel.recv()
        if err != nil {
            return err
        }
        step = Dispatch(event, handler, arg)
    }
    return nil
}

// Close the target subprocess; returns when the target subprocess has exited.
func (reel *Reel) Close() {
    reel.stdin.Close()
    reel.subp.Wait()
    reel.subp = nil
}

// Run the target subprocess to completion.
// The first step to take is `step`; consequent steps are determined by
// dispatching events via `handler` with `arg`.
// Return on first error, or when there is no next step to perform.
func (reel *Reel) Run(step *Step, handler Handler, arg interface{}) error {
    err := reel.Open()
    if err == nil {
        err = reel.Step(step, handler, arg)
        reel.Close()
    }
    return err
}

// Run the target subprocess to completion via channels.
// Steps to take are read from `cstep`; notifications are written to `cnotify`.
// Return on first error, or when there is no next step to perform.
// TODO: document notification protocol
// TODO: document completion `token`
func (reel *Reel) Runch(cstep chan *Step, cnotify chan *Notification, token interface{}) error {
    var event *Event
    err := reel.Open()
    if err == nil {
        for step := range cstep {
            err = reel.send(step)
            if err != nil {
                break
            }
            if len(step.Expect) > 0 {
                event, err = reel.recv()
                if err != nil {
                    break
                }
            } else {
                event = nil
            }
            cnotify <- &Notification{Event: event}
        }
        reel.Close()
    }
    if err != nil {
        cnotify <- &Notification{Error: err}
    }
    cnotify <- &Notification{Token: token}
    return err
}

func prependLogOption(args []string, logfile string) []string {
    args = append(args, "", "")
    copy(args[2:], args)
    args[0] = "-l"
    args[1] = logfile
    return args
}

// Create a new `Reel` instance for interacting with a target subprocess.
// The command line for the target is specified in `args`.
// Optionally log dialogue with the subprocess to `logfile`.
func NewReel(logfile string, args []string) (*Reel, error) {
    var err error

    if logfile != "" {
        args = prependLogOption(args, logfile)
    }

    subp := exec.Command("reel.exp", args[:]...)
    stdin, err := subp.StdinPipe()
    if err != nil {
        return nil, err
    }
    stdout, err := subp.StdoutPipe()
    if err != nil {
        stdin.Close()
        return nil, err
    }
    return &Reel{
        subp:    subp,
        stdin:   stdin,
        stdout:  stdout,
        scanner: bufio.NewScanner(stdout),
    }, err
}
