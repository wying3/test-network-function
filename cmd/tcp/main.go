package main

import (
    "flag"
    "fmt"
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
    "github.com/redhat-nfvpe/test-network-function/pkg/tnf"
    "os"
)

type TcpExpr struct {
    idx     int
    bhost   string
    shost   string
    timeout int
}

func (t *TcpExpr) Next() bool {
    t.idx++
    return t.idx < 1
}
func (t *TcpExpr) Eval(executor *tnf.Executor) (int, error) {
    tcp := tnf.NewTcp(t.timeout, t.bhost, t.shost)
    return executor.Exec(tcp)
}

func parseArgs() (string, tnf.Expr) {
    logfile := flag.String("d", "", "Filename prefix for expect dialogues")
    timeout := flag.Int("t", 2, "Timeout in seconds")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "usage: %s [-d logprefix] [-t timeout] bind-host server-host\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(tnf.ERROR)
    }
    flag.Parse()
    hosts := flag.Args()
    if len(hosts) != 2 {
        flag.Usage()
    }
    return *logfile, &TcpExpr{
        idx: -1,
        bhost: hosts[0],
        shost: hosts[1],
        timeout: *timeout,
    }
}

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
