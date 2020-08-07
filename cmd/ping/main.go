package main

import (
    "flag"
    "fmt"
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
    "github.com/redhat-nfvpe/test-network-function/pkg/tnf"
    "os"
)

func parseArgs() (string, *tnf.Ping) {
    logfile := flag.String("d", "", "Filename for expect dialogue")
    timeout := flag.Int("t", 2, "Timeout in seconds")
    count := flag.Int("c", 1, "Number of requests to send")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "usage: %s [-d logfile] [-t timeout] [-c count] host\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(tnf.ERROR)
    }
    flag.Parse()
    hosts := flag.Args()
    if len(hosts) != 1 {
        flag.Usage()
    }
    return *logfile, tnf.NewPing(*timeout, hosts[0], *count)
}

// Execute a ping test with exit code 0 on success, 1 on failure, 2 on error.
// Print interaction with the controlled subprocess which implements the test.
// Optionally log dialogue with the controlled subprocess to file.
func main() {
    arg := 0
    result := tnf.ERROR
    logfile, test := parseArgs()
    printer := reel.NewPrinter("")
    chain := reel.NewChain([]reel.Handler{printer, test})
    args, err := test.Args(arg)
    if err == nil {
        reel, err := reel.NewReel(logfile, args)
        if err == nil {
            err = reel.Run(test.ReelFirst(arg), chain, arg)
            if err == nil {
                result = test.Result()
            }
        }
    }
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    os.Exit(result)
}
