# Test Network Function

This repository contains a set of network function test cases.

## Install

0. Install packages for dependencies: `expect` and `tcllib`.
1. Ensure [`bin/reel.exp`](bin/reel.exp) is on your PATH.

## Command Line Tools

A set of command line tools is provided in [`cmd/`](cmd/). The intention is to
provide tools which are useful in their own right; and to serve as a set of
reference implementations for programatically composing and executing tests.

When implementing a new command line tool, the following conventions must be
adhered to.

### Subprocess dialog

Where a command line tool runs subprocesses using [`bin/reel.exp`](bin/reel.exp)
(see [`internal/reel/reel.go`](internal/reel/reel.go)) then it must provide a
`-d` option to allow subprocess interaction dialog to be captured. This option
takes a string, which the tool may interpret as either:

* a filename for the dialog when exactly one subprocess is always run; or,
* a filename prefix for the dialogs with a variable number of subprocesses.

The dialog captures lines of JSON sent by the controlling process to an instance
of [`bin/reel.exp`](bin/reel.exp); the lines of JSON sent in return; and the
text output by the controlled subprocess.

### Exit code

The exit code must use [test result codes](pkg/tnf/result.go) where:

* `0` (zero) indicates success;
* `1` indicates failure;
* `2` indicates error.

## Hacking

[`bin/reel.exp`](bin/reel.exp) was written to allow control logic to be written
in any language, not just [Tcl/Tk](https://www.tcl.tk/). The following snippets
demonstrate experimenting from the command line:

```
$ reel.exp -l hacking.log bash -c 'while read LINE; do echo foo$LINE; done' < <(
> echo '{"execute": "bar", "expect": ["foobar"]}'
> sleep 1
> echo '{"execute": "baz", "expect": ["foobaz"]}'
> sleep 1
> echo '{"execute": "\u0004", "expect": ["\u0004"]}'
> )
{"event":"match","idx":0,"pattern":"foobar","before":"bar\r\n","match":"foobar"}
{"event":"match","idx":0,"pattern":"foobaz","before":"\r\nbaz\r\n","match":"foobaz"}
{"event":"eof"}
```

```
running: bash -c {while read LINE; do echo foo$LINE; done}
{"execute": "bar", "expect": ["foobar"]}
bar
foobar
{"event":"match","idx":0,"pattern":"foobar","before":"bar\r\n","match":"foobar"}
{"execute": "baz", "expect": ["foobaz"]}
baz
foobaz
{"event":"match","idx":0,"pattern":"foobaz","before":"\r\nbaz\r\n","match":"foobaz"}
{"execute": "\u0004", "expect": ["\u0004"]}

{"event":"eof"}
```

## Running Tests

In order to run the CNF tests, issue the following command:

```shell script
make cnftests
```

A JUnit report containing results is created at `test-network-functions/test-network-function_junit.xml`
