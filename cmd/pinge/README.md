# pinge

The `pinge` tool complements the [`ping`](../ping/) tool. It allows multiple
hosts to be pinged with the test result determined by a boolean expression.
(The trailing 'e' stands for 'expression'.) The `-e` option allows a choice of
expression. Expressions operate on test success and failure results: errors are
always reported immediately.

## all

Specifying the `-e` option as `all` employs a boolean AND expression. This is
the default. Each host will be pinged until the first failure. If all hosts were
successfully pinged, the test reports success.

```
$ go run cmd/pinge/main.go 10.5.0.3 10.5.0.1
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=4.42 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 4.429/4.429/4.429/0.000 ms
PING 10.5.0.1 (10.5.0.1) 56(84) bytes of data.
64 bytes from 10.5.0.1: icmp_seq=1 ttl=63 time=2.27 ms

--- 10.5.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 2.272/2.272/2.272/0.000 ms
```

```
$ go run cmd/pinge/main.go -e all 10.5.0.3 10.3.0.99 10.5.0.1
PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=4.29 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 4.292/4.292/4.292/0.000 ms
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

exit status 1
```

## any

Specifying the `-e` option as `any` employs a boolean OR expression. Each host
will be pinged until the first success. If any host was successfully pinged, the
test reports success.

```
$ go run cmd/pinge/main.go -e any 10.3.0.99 10.3.0.88
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

(timeout)
PING 10.3.0.88 (10.3.0.88) 56(84) bytes of data.
^C

--- 10.3.0.88 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

exit status 1
```

```
$ go run cmd/pinge/main.go -e any 10.3.0.99 10.5.0.3 10.3.0.88
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=4.27 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 4.277/4.277/4.277/0.000 ms
```

## one

Specifying the `-e` option as `one` employs a boolean XOR expression. All hosts
will be pinged unless more than one host was successfully pinged. If exactly one
host was successfully pinged, the test reports success.

```
$ go run cmd/pinge/main.go -e one 10.3.0.99 10.3.0.88
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

(timeout)
PING 10.3.0.88 (10.3.0.88) 56(84) bytes of data.
^C

--- 10.3.0.88 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

exit status 1
```

```
$ go run cmd/pinge/main.go -e one 10.3.0.99 10.5.0.3 10.3.0.88
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=5.65 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 5.654/5.654/5.654/0.000 ms
(timeout)
PING 10.3.0.88 (10.3.0.88) 56(84) bytes of data.
^C

--- 10.3.0.88 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms
```

```
$ go run cmd/pinge/main.go -e one 10.3.0.99 10.5.0.3 10.3.0.88 10.5.0.1 10.3.0.77
(timeout)
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
^C

--- 10.3.0.99 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

PING 10.5.0.3 (10.5.0.3) 56(84) bytes of data.
64 bytes from 10.5.0.3: icmp_seq=1 ttl=63 time=6.55 ms

--- 10.5.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 6.552/6.552/6.552/0.000 ms
(timeout)
PING 10.3.0.88 (10.3.0.88) 56(84) bytes of data.
^C

--- 10.3.0.88 ping statistics ---
1 packets transmitted, 0 received, 100% packet loss, time 0ms

PING 10.5.0.1 (10.5.0.1) 56(84) bytes of data.
64 bytes from 10.5.0.1: icmp_seq=1 ttl=63 time=8.86 ms

--- 10.5.0.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 8.869/8.869/8.869/0.000 ms
exit status 1
```
