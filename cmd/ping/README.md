# ping

By default, the `ping` tool sends a single ICMP Echo Request to a target host
and expects the test to complete within 2 seconds.

```
$ go run cmd/ping/main.go -d success.log 10.5.0.3
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=38.6 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 38.646/38.646/38.646/0.000 ms
```

This test reports success as at least one request was transmitted and the same
number of responses were received.

```
$ echo $?
0
```

For this example, `success.log` contains the dialog observed by `bin/reel.exp`.

```
$ cat success.log
{"expect":["\\D\\d+ packets transmitted.*\\r\\n(?:rtt )?.*$"],"timeout":2}
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=38.6 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 38.646/38.646/38.646/0.000 ms
{"event":"match","idx":0,"pattern":"\\D\\d+ packets transmitted.*\\r\\n(?:rtt )?.*$","before":"PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.\r\n64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=38.6 ms\r\n\r\n--- 10.5.0.3 ping statistics ---\r","match":"\n1 packets transmitted, 1 received, 0% packet loss, time 0ms\r\nrtt min/avg/max/mdev = 38.646/38.646/38.646/0.000 ms\r\n"}
```

## timeout

The following example shows the controlling process giving up on the test due
to a timeout and killing the controlled subprocess.

```
$ go run cmd/ping/main.go -d failure.log 10.3.0.99
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

exit status 1
```

This test reports failure as at least one request was transmitted but no
responses were received.

The dialog clearly shows the timeout and sending of ^C to kill the subprocess.

```
$ cat failure.log
{"expect":["\\D\\d+ packets transmitted.*\\r\\n(?:rtt )?.*$"],"timeout":2}
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
{"event":"timeout"}
{"execute":"\u0003","expect":["\\D\\d+ packets transmitted.*\\r\\n(?:rtt )?.*$"]}
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

{"event":"match","idx":0,"pattern":"\\D\\d+ packets transmitted.*\\r\\n(?:rtt )?.*$","before":"PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.\r\n^C\r\n\r\n--- 10.3.0.99 ping statistics ---\r","match":"\n1 packets transmitted, 0 received, 100% packet loss, time 0ms\r\n\r\n"}
```

## error

The following example shows a test-specific error in executing the test.

```
$ go run cmd/ping/main.go -c 8 -t 10 10.5.0.99
PING 10.5.0.99 (10.5.0.99) 56(84) bytes of data.
From 10.3.0.1 icmp_seq=1 Destination Host Unreachable
From 10.3.0.1 icmp_seq=2 Destination Host Unreachable
From 10.3.0.1 icmp_seq=3 Destination Host Unreachable
From 10.3.0.1 icmp_seq=4 Destination Host Unreachable
From 10.3.0.1 icmp_seq=5 Destination Host Unreachable
From 10.3.0.1 icmp_seq=6 Destination Host Unreachable
From 10.3.0.1 icmp_seq=7 Destination Host Unreachable
From 10.3.0.1 icmp_seq=8 Destination Host Unreachable

--- 10.5.0.99 ping statistics ---
8 packets transmitted, 0 received, +8 errors, 100% packet loss, time 7121ms
pipe 4
exit status 2
```

This test reports error due to ping reporting errors.

## ping with N requests

To send N requests, use the `-c` option.

```
$ go run cmd/ping/main.go -c 2 10.5.0.3
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=5.12 ms
64 bytes from 10.5.0.3: icmp_seq=2 ttl=63 time=8.92 ms

--- 10.5.0.3 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 5.128/7.024/8.920/1.896 ms
```

Note that the default test timeout still applies...

```
$ go run cmd/ping/main.go -c 10 10.5.0.3
(timeout)
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=4.27 ms
64 bytes from 10.5.0.3: icmp_seq=2 ttl=63 time=8.72 ms
^C

--- 10.5.0.3 ping statistics ---
3 packets transmitted, 2 received, 33% packet loss, time 2003ms
rtt min/avg/max/mdev = 4.271/6.498/8.726/2.228 ms
```

...so specify an appropriate test timeout with `-t`.

```
$ go run cmd/ping/main.go -c 10 -t 12 10.5.0.3
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=4.87 ms
64 bytes from 10.5.0.3: icmp_seq=2 ttl=63 time=9.73 ms
64 bytes from 10.5.0.3: icmp_seq=3 ttl=63 time=11.4 ms
64 bytes from 10.5.0.3: icmp_seq=4 ttl=63 time=11.1 ms
64 bytes from 10.5.0.3: icmp_seq=5 ttl=63 time=11.0 ms
64 bytes from 10.5.0.3: icmp_seq=6 ttl=63 time=4.00 ms
64 bytes from 10.5.0.3: icmp_seq=7 ttl=63 time=11.0 ms
64 bytes from 10.5.0.3: icmp_seq=8 ttl=63 time=4.70 ms
64 bytes from 10.5.0.3: icmp_seq=9 ttl=63 time=4.33 ms
64 bytes from 10.5.0.3: icmp_seq=10 ttl=63 time=12.2 ms

--- 10.5.0.3 ping statistics ---
10 packets transmitted, 10 received, 0% packet loss, time 9012ms
rtt min/avg/max/mdev = 4.001/8.454/12.290/3.306 ms
```

## ping for T seconds

To ping continuously for T seconds, specify a count of 0 (zero).

```
$ go run cmd/ping/main.go -c 0 -t 20 10.5.0.3
(timeout)
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=5.36 ms
64 bytes from 10.5.0.3: icmp_seq=2 ttl=63 time=25.8 ms
64 bytes from 10.5.0.3: icmp_seq=3 ttl=63 time=4.46 ms
64 bytes from 10.5.0.3: icmp_seq=4 ttl=63 time=5.71 ms
64 bytes from 10.5.0.3: icmp_seq=5 ttl=63 time=4.24 ms
64 bytes from 10.5.0.3: icmp_seq=6 ttl=63 time=11.0 ms
64 bytes from 10.5.0.3: icmp_seq=7 ttl=63 time=10.8 ms
64 bytes from 10.5.0.3: icmp_seq=8 ttl=63 time=4.63 ms
64 bytes from 10.5.0.3: icmp_seq=9 ttl=63 time=11.1 ms
64 bytes from 10.5.0.3: icmp_seq=10 ttl=63 time=6.88 ms
64 bytes from 10.5.0.3: icmp_seq=11 ttl=63 time=5.34 ms
64 bytes from 10.5.0.3: icmp_seq=12 ttl=63 time=4.75 ms
64 bytes from 10.5.0.3: icmp_seq=13 ttl=63 time=9.06 ms
64 bytes from 10.5.0.3: icmp_seq=14 ttl=63 time=8.92 ms
64 bytes from 10.5.0.3: icmp_seq=15 ttl=63 time=8.74 ms
64 bytes from 10.5.0.3: icmp_seq=16 ttl=63 time=8.86 ms
64 bytes from 10.5.0.3: icmp_seq=17 ttl=63 time=5.09 ms
64 bytes from 10.5.0.3: icmp_seq=18 ttl=63 time=10.3 ms
64 bytes from 10.5.0.3: icmp_seq=19 ttl=63 time=4.30 ms
64 bytes from 10.5.0.3: icmp_seq=20 ttl=63 time=4.66 ms
^C

--- 10.5.0.3 ping statistics ---
21 packets transmitted, 20 received, 4% packet loss, time 20031ms
rtt min/avg/max/mdev = 4.249/8.014/25.816/4.774 ms
```
