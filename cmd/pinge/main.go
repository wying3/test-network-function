package main

import (
    "flag"
    "fmt"
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
    "github.com/redhat-nfvpe/test-network-function/pkg/tnf"
    "os"
)

type Pings struct {
    idx     int
    hosts   []string
    count   int
    timeout int
}

func (p *Pings) Next() bool {
    p.idx++
    return p.idx < len(p.hosts)
}
func (p *Pings) Eval(executor *tnf.Executor) (int, error) {
    ping := tnf.NewPing(p.timeout, p.hosts[p.idx], p.count)
    return executor.Exec(ping)
}

func parseArgs() (string, tnf.Expr) {
    logfile := flag.String("d", "", "Filename prefix for expect dialogues")
    timeout := flag.Int("t", 2, "Timeout in seconds")
    count := flag.Int("c", 1, "Number of requests to send")
    combiner := flag.String("e", "all", "Expression...")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "usage: %s [-d logprefix] [-t timeout] [-c count] [-e all|any|one] host ?host .. host?\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(tnf.ERROR)
    }
    flag.Parse()
    hosts := flag.Args()
    if len(hosts) == 0 {
        flag.Usage()
    }
    pings := &Pings{idx: -1, hosts: hosts, count: *count, timeout: *timeout}
    switch *combiner {
    case "any":
        return *logfile, tnf.Any(pings)
    case "one":
        return *logfile, tnf.One(pings)
    default:
        return *logfile, tnf.All(pings)
    }
}

// Execute a ping test with exit code 0 on success, 1 on failure, 2 on error.
// Print interaction with the controlled subprocess which implements the test.
// Optionally log dialogue with the controlled subprocess to file.
func main() {
    logfile, expr := parseArgs()
    printer := reel.NewPrinter("")
    executor := tnf.NewExecutor(logfile, []reel.Handler{printer})
    result, err := expr.Eval(executor)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    os.Exit(result)
}
