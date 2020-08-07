package tnf

import (
    "bufio"
    "fmt"
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
    "os"
    "strings"
)

type TestFeeder struct {
    timeout int
    prompt  string
    scanner *bufio.Scanner
    tester  Test
}

func (f *TestFeeder) ReelMatch(pattern string, before string, match string, arg interface{}) *reel.Step {
    if f.scanner != nil && f.scanner.Scan() {
        config, err := DecodeConfig(f.scanner.Bytes())
        if err == nil {
            // TODO: no such test => panic
            f.tester = Tests[config.Test](config)
            args, err := f.tester.Args(nil)
            if err == nil {
                return &reel.Step{
                    Execute: strings.Join(args, " "),
                    Expect:  []string{f.prompt},
                }
                // TODO: fold in result?
            }
        }
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            f.scanner = nil
        }
    }
    f.tester = nil
    return nil
}
func (f *TestFeeder) ReelTimeout(arg interface{}) *reel.Step {
    if f.tester != nil {
        return &reel.Step{
            Execute: reel.CTRL_C,
            Expect:  []string{f.prompt},
            Timeout: f.timeout,
        }
    }
    return nil
}
func (f *TestFeeder) ReelEof(arg interface{}) *reel.Step {
    return nil
}

func NewTestFeeder(timeout int, prompt string, scanner *bufio.Scanner) *TestFeeder {
    return &TestFeeder{
        timeout: timeout,
        prompt:  prompt,
        scanner: scanner,
    }
}
