package tnf

// Test result codes, defined as follows:
//
// `SUCCESS`: the test procedure executed without error and the condition for
// success was achieved;
//
// `FAILURE`: the test procedure executed without error but the condition for
// success was not achieved;
//
// `ERROR`: the test procedure could not execute due to an error; or the test
// procedure executed without error but a test-specific error condition was
// detected.
const (
    SUCCESS = iota
    FAILURE
    ERROR
)
