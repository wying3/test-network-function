package tnf

import (
    "fmt"
    "strconv"
    "github.com/redhat-nfvpe/test-network-function/internal/reel"
)

type Test interface {
    Args(arg interface{}) ([]string, error)
    Result() int
    ReelFirst(arg interface{}) *reel.Step
}
type part struct {
    id int
    runner *reel.Reel
    chain *reel.Chain
    cstep chan *reel.Step
}
type Executor struct {
    logfile string
    chain   *reel.Chain
    cnotify chan *reel.Notification
    running map[int]*part
    lastId int
    activeId int
}
type Context struct {
    executor *Executor
    id int
    err error
}

type NoContextError struct {
    test interface{}
}

func (e *NoContextError) Error() string {
    return fmt.Sprintf("Test %T can only be run by Executor", e.test)
}

type ProgrammingError struct {
    reason string
}

func (e *ProgrammingError) Error() string {
    return "Programming error: " + e.reason
}

func (e *Executor) Start(test Test) (int, error) {
    e.lastId += 1
    id := e.lastId
    logfile := e.logfile
    if logfile != "" {
        logfile += "." + strconv.Itoa(id)
    }
    context := &Context{e, id, nil}
    args, err := test.Args(context)
    if err == nil {
        err = context.err
    }
    if err == nil {
        runner, err := reel.NewReel(logfile, args)
        if err == nil {
            handlers := []reel.Handler{test.(reel.Handler)}
            chain := e.chain.Extended(handlers)
            cstep := make(chan *reel.Step, 1)
            p := &part{id: id, runner: runner, chain: chain, cstep: cstep}
            e.running[p.id] = p
            go p.runner.Runch(p.cstep, e.cnotify, p.id)
            return p.id, nil
        }
    }
    return 0, err
}

func (e *Executor) Exec(test Test) (int, error) {
    id, err := e.Start(test)
    if err == nil {
        e.activeId = id
        var context Context = Context{e, e.activeId, nil}
        e.running[e.activeId].cstep <- test.ReelFirst(&context)
        for notification := range e.cnotify {
            if notification.Error != nil {
                err = notification.Error
                break
            } else if notification.Event != nil {
                handler := e.running[e.activeId].chain
                step := reel.Dispatch(notification.Event, handler, &context)
                if step == nil || context.err != nil {
                    err = context.err
                    break
                }
                e.activeId = context.id
                e.running[e.activeId].cstep <- step
            }
        }
    }
    for id := range e.running {
        close(e.running[id].cstep)
        notification := <-e.cnotify
        if notification.Token.(int) != 0 {
            delete(e.running, id)
        }
    }
    if err != nil {
        return ERROR, err
    } else {
        return test.Result(), nil
    }
}

func NewExecutor(logfile string, handlers []reel.Handler) *Executor {
    return &Executor{
        logfile: logfile,
        chain: reel.NewChain(handlers),
        cnotify: make(chan *reel.Notification, 1),
        running: make(map[int]*part),
    }
}
