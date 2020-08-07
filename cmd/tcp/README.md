# tcp

The `tcp` tool demonstrates a complex test implementation: running a TCP server
on an arbitrary port, connecting a TCP client to that server, generating and
sending a random magic string from client to server, finally ensuring that the
same string is received at the server.

```
$ go run cmd/tcp/main.go -d tcp.log localhost localhost
listening on [127.0.0.1] 44887 ...Wu1UFA==
connect to [127.0.0.1] from localhost [127.0.0.1] 47680
Wu1UFA==(eof)
(eof)
```

The dialogues with the server and client subprocesses are captured separately.

The server dialogue:

```
$ cat tcp.log.1
running: sh -c {nc -s localhost -l -v 2>&1}
{"expect":["listening on \\[.*\\] (\\d+) \\.\\.\\."],"timeout":2}
listening on [127.0.0.1] 44887 ...{"event":"match","idx":0,"pattern":"listening on \\[.*\\] (\\d+) \\.\\.\\.","before":"","match":"listening on [127.0.0.1] 44887 ..."}
{"expect":["Wu1UFA=="],"timeout":2}

connect to [127.0.0.1] from localhost [127.0.0.1] 47680
Wu1UFA==
{"event":"match","idx":0,"pattern":"Wu1UFA==","before":"\r\nconnect to [127.0.0.1] from localhost [127.0.0.1] 47680\r\n","match":"Wu1UFA=="}
{"expect":["\u0004"],"timeout":2}
{"event":"eof"}
```

The client dialogue:

```
$ cat tcp.log.2
running: nc -q 0 localhost 44887
{"execute":"Wu1UFA==","expect":["Wu1UFA=="]}
Wu1UFA==
{"event":"match","idx":0,"pattern":"Wu1UFA==","before":"","match":"Wu1UFA=="}
{"execute":"\u0004","expect":["\u0004"],"timeout":2}

{"event":"eof"}
```
