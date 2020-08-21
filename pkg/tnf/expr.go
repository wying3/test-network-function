package tnf

type Expr interface {
    Next()(bool)
    Eval(executor *Executor)(int, error)
}

type all struct {
    exprs Expr
}

func (expr *all) Next() bool {
    return expr.exprs.Next()
}
func (expr *all) Eval(executor *Executor) (int, error) {
    for expr.Next() {
        result, err := expr.exprs.Eval(executor)
        if result > SUCCESS || err != nil {
            return result, err
        }
    }
    return SUCCESS, nil
}

func All(exprs Expr) Expr {
    return &all{exprs: exprs}
}

type any struct {
    exprs Expr
}

func (expr *any) Next() bool {
    return expr.exprs.Next()
}
func (expr *any) Eval(executor *Executor) (int, error) {
    var err error = nil
    result := SUCCESS
    for expr.Next() {
        result, err = expr.exprs.Eval(executor)
        if result != FAILURE || err != nil {
            break
        }
    }
    return result, err
}

func Any(exprs Expr) Expr {
    return &any{exprs: exprs}
}

type one struct {
    exprs Expr
}

func (expr *one) Next() bool {
    return expr.exprs.Next()
}
func (expr *one) Eval(executor *Executor) (int, error) {
    var succeeded int
    for expr.Next() {
        result, err := expr.exprs.Eval(executor)
        if result == ERROR || err != nil {
            return result, err
        }
        if result == SUCCESS {
            succeeded += 1
            if succeeded > 1 {
                break
            }
        }
    }
    if succeeded == 1 {
        return SUCCESS, nil
    } else {
        return FAILURE, nil
    }
}

func One(exprs Expr) Expr {
    return &one{exprs: exprs}
}
