package tnf

import (
    "testing"
)

// Unit tests for expr.go

// A fact, a fixed (result, error, contributes) "tuple":
// `result` is a fixed result code to return when this fact is evaluated;
// `err` is a fixed error to return when this fact is evaluated;
// `contributes` is a fixed boolean indicating whether this fact must be used
// in evaluating an expression result;
// `touched` is a variable used to mark whether this fact was returned during
// evaluation.
type fact struct {
    result      int
    err         error
    contributes bool
    touched     bool
}
func factSuccess(contributes bool) *fact {
    return &fact{SUCCESS, nil, contributes, false}
}
func factFailure(contributes bool) *fact {
    return &fact{FAILURE, nil, contributes, false}
}
func factError(err error, contributes bool) *fact {
    return &fact{ERROR, err, contributes, false}
}

// A reductive error for testing.
type _error string
func (e _error) Error() string {
    return string(e)
}

// A sequence of facts for expression evaluation.
type seq struct {
    idx     int
    facts   []*fact
}
func (s *seq) Next() bool {
    s.idx++
    return s.idx < len(s.facts)
}
func (s *seq) Eval(e *Executor)(int, error) {
    s.facts[s.idx].touched = true
    return s.facts[s.idx].result, s.facts[s.idx].err
}

type Func func(exprs Expr) Expr

// A function for testing that expression evaluation of `facts` via `fn`
// returns `result`, `err` and only uses contributing facts in that evaluation.
func test(t *testing.T, fn Func, facts []*fact, result int, err error) {
    r, e := fn(&seq{-1, facts}).Eval(nil)
    testResult(t, result, r)
    testError(t, err, e)
    testFacts(t, facts)
}
func testResult(t *testing.T, exp int, got int) {
    if exp != got {
        t.Errorf("expected %d but got %d\n", exp, got)
    }
}
func testError(t *testing.T, exp error, got error) {
    if exp != got {
        t.Errorf("expected %s but got %s\n", exp, got)
    }
}
func testFacts(t *testing.T, facts []*fact) {
    for i := 0; i < len(facts); i++ {
        if facts[i].contributes == facts[i].touched {
            continue
        } else if facts[i].contributes {
            t.Errorf("fact %d did not contribute to evaluation\n", i)
        } else if facts[i].touched {
            t.Errorf("fact %d should not contribute to evaluation\n", i)
        }
    }
}

// Test `All` for SUCCESS.
func TestAllSuccess(t *testing.T) {
    facts := []*fact{
        factSuccess(true),
        factSuccess(true),
        factSuccess(true),
        factSuccess(true),
        factSuccess(true),
    }
    test(t, All, facts, SUCCESS, nil)
}

// Test `All` for FAILURE.
func TestAllFailure(t *testing.T) {
    facts := []*fact{
        factSuccess(true),
        factFailure(true),
        factSuccess(false),
        factFailure(false),
        factError(nil, false),
    }
    test(t, All, facts, FAILURE, nil)
}

// Test `All` for ERROR.
func TestAllError(t *testing.T) {
    var err _error = "(expected error)"
    var bad _error = "unexpected error!"
    // without err
    facts := []*fact{
        factSuccess(true),
        factError(nil, true),
        factSuccess(false),
        factFailure(false),
        factError(nil, false),
    }
    test(t, All, facts, ERROR, nil)
    // with err
    facts = []*fact{
        factSuccess(true),
        factError(err, true),
        factSuccess(false),
        factFailure(false),
        factError(bad, false),
    }
    test(t, All, facts, ERROR, err)
}

// Test `Any` for SUCCESS.
func TestAnySuccess(t *testing.T) {
    facts := []*fact{
        factFailure(true),
        factFailure(true),
        factSuccess(true),
        factSuccess(false),
        factFailure(false),
        factError(nil, false),
    }
    test(t, Any, facts, SUCCESS, nil)
}

// Test `Any` for FAILURE.
func TestAnyFailure(t *testing.T) {
    facts := []*fact{
        factFailure(true),
        factFailure(true),
        factFailure(true),
        factFailure(true),
        factFailure(true),
    }
    test(t, Any, facts, FAILURE, nil)
}

// Test `Any` for ERROR.
func TestAnyError(t *testing.T) {
    var err _error = "(expected error)"
    var bad _error = "unexpected error!"
    // without err
    facts := []*fact{
        factFailure(true),
        factError(nil, true),
        factSuccess(false),
        factFailure(false),
        factError(nil, false),
    }
    test(t, Any, facts, ERROR, nil)
    // with err
    facts = []*fact{
        factFailure(true),
        factError(err, true),
        factSuccess(false),
        factFailure(false),
        factError(bad, false),
    }
    test(t, Any, facts, ERROR, err)
}

// Test `One` for SUCCESS.
func TestOneSuccess(t *testing.T) {
    facts := []*fact{
        factFailure(true),
        factFailure(true),
        factSuccess(true),
        factFailure(true),
        factFailure(true),
    }
    test(t, One, facts, SUCCESS, nil)
}

// Test `One` for FAILURE.
func TestOneFailure(t *testing.T) {
    facts := []*fact{
        factFailure(true),
        factSuccess(true),
        factFailure(true),
        factSuccess(true),
        factFailure(false),
    }
    test(t, One, facts, FAILURE, nil)
}

// Test `One` for ERROR.
func TestOneError(t *testing.T) {
    var err _error = "(expected error)"
    var bad _error = "unexpected error!"
    // without err
    facts := []*fact{
        factSuccess(true),
        factError(nil, true),
        factSuccess(false),
        factFailure(false),
        factError(nil, false),
    }
    test(t, One, facts, ERROR, nil)
    // with err
    facts = []*fact{
        factSuccess(true),
        factError(err, true),
        factSuccess(false),
        factFailure(false),
        factError(bad, false),
    }
    test(t, One, facts, ERROR, err)
}
