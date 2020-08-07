package reel

type Chain struct {
    handlers []Handler
}

type StepFunc func(Handler) *Step

func (c *Chain) dispatch(fp StepFunc) *Step {
    for _, handler := range c.handlers {
        step := fp(handler)
        if step != nil {
            return step
        }
    }
    return nil
}

func (c *Chain) ReelMatch(pattern string, before string, match string, arg interface{}) *Step {
    fp := func(handler Handler) *Step {
        return handler.ReelMatch(pattern, before, match, arg)
    }
    return c.dispatch(fp)
}

func (c *Chain) ReelTimeout(arg interface{}) *Step {
    fp := func(handler Handler) *Step {
        return handler.ReelTimeout(arg)
    }
    return c.dispatch(fp)
}

func (c *Chain) ReelEof(arg interface{}) *Step {
    fp := func(handler Handler) *Step {
        return handler.ReelEof(arg)
    }
    return c.dispatch(fp)
}

func (c *Chain) Extended(handlers []Handler) *Chain {
    cath := make([]Handler, len(c.handlers))
    copy(cath, c.handlers)
    return &Chain{handlers: append(cath, handlers...)}
}

func NewChain(handlers []Handler) *Chain {
    return &Chain{handlers: handlers}
}
