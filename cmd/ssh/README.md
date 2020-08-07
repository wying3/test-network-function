# ssh

By default, the `ssh` tool simply establishes a SSH session to the target host,
then closes it. The controlled subprocess is the native `ssh` client tool, with
its command line options and args passed through. The session is closed when the
supplied prompt string (regex) is matched.

```
$ env TERM=vt220 go run cmd/ssh/main.go 'user@hhh:\S+\$ ' hhh -o 'PreferredAuthentications=publickey'
Last login: Thu Jul 23 18:10:34 2020 from 10.3.0.109
user@hhh:~$ logout
Connection to hhh.ddd closed.
```

## Interactive session

Specifying `-f lines` allows interactive execution of lines of text from stdin,
mimicking the native client tool.

In this example, two commands are supplied from a here document.

```
$ env TERM=vt220 go run cmd/ssh/main.go -d ssh.log -f lines 'user@hhh:\S+\$ ' hhh -o 'PreferredAuthentications=publickey' <<EOF
> echo foobar
> date
> EOF
Last login: Tue Jul 14 15:58:14 2020 from 10.3.0.109
user@hhh:~$ echo foobar
foobar
user@hhh:~$ date
Tue 14 Jul 15:59:08 BST 2020
user@hhh:~$ logout

$ echo $?
0

$ cat ssh.log
{"expect":["Are you sure you want to continue connecting \\(yes/no\\)\\?","Please type 'yes' or 'no': ","user@hhh:\\S+\\$ "],"timeout":2}
Last login: Tue Jul 14 15:58:14 2020 from 10.3.0.109
user@hhh:~$ {"event":"match","idx":2,"pattern":"user@hhh:\\S+\\$ ","before":"Last login: Tue Jul 14 15:58:14 2020 from 10.3.0.109\r\r\n","match":"user@hhh:~$ "}
{"execute":"echo foobar","expect":["user@hhh:\\S+\\$ "],"timeout":2}
echo foobar
foobar
user@hhh:~$ {"event":"match","idx":0,"pattern":"user@hhh:\\S+\\$ ","before":"echo foobar\r\nfoobar\r\n","match":"user@hhh:~$ "}
{"execute":"date","expect":["user@hhh:\\S+\\$ "],"timeout":2}
date
Tue 14 Jul 15:59:08 BST 2020
user@hhh:~$ {"event":"match","idx":0,"pattern":"user@hhh:\\S+\\$ ","before":"date\r\nTue 14 Jul 15:59:08 BST 2020\r\n","match":"user@hhh:~$ "}
{"execute":"\u0004","expect":["Connection to .+ closed\\..*$"],"timeout":2}
logout
{"event":"match","idx":0,"pattern":"Connection to .+ closed\\..*$","before":"logout\r\n","match":"Connection to hhh.ddd closed.\r"}
```

## Test execution

Specifying `-f tests` allows execution of test configurations; each line read
from stdin must contain a complete JSON object defining a valid configuration.

```
$ cat tests.json
{"test": "https://tnf.redhat.com/ping/one", "host": "ggg"}
{"test": "https://tnf.redhat.com/ping/flexi", "count": 77, "host": "10.3.0.99"}

$ env TERM=vt220 go run cmd/ssh/main.go -d ssh.log -f tests 'user@hhh:\S+\$ ' hhh -o 'PreferredAuthentications=publickey' <tests.json
Last login: Tue Jul 14 15:59:07 2020 from 10.3.0.109
user@hhh:~$ ping -c 1 ggg
PING ggg.ddd (10.5.0.5) 56(84) bytes of data.
64 bytes from ggg.ddd (10.5.0.5): icmp_seq=1 ttl=64 time=0.344 ms

--- ggg.ddd ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.344/0.344/0.344/0.000 ms
user@hhh:~$ (timeout)
ping -c 77 10.3.0.99
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
From 10.5.0.1 icmp_seq=1 Destination Host Unreachable
From 10.5.0.1 icmp_seq=2 Destination Host Unreachable
From 10.5.0.1 icmp_seq=3 Destination Host Unreachable
From 10.5.0.1 icmp_seq=4 Destination Host Unreachable
From 10.5.0.1 icmp_seq=5 Destination Host Unreachable
From 10.5.0.1 icmp_seq=6 Destination Host Unreachable
From 10.5.0.1 icmp_seq=7 Destination Host Unreachable
From 10.5.0.1 icmp_seq=8 Destination Host Unreachable
^C

--- 10.3.0.99 ping statistics ---
11 packets transmitted, 0 received, +8 errors, 100% packet loss, time 10183ms
pipe 4
user@hhh:~$ 
user@hhh:~$ logout

$ cat ssh.log
{"expect":["Are you sure you want to continue connecting \\(yes/no\\)\\?","Please type 'yes' or 'no': ","user@hhh:\\S+\\$ "],"timeout":2}
Last login: Tue Jul 14 15:59:07 2020 from 10.3.0.109
user@hhh:~$ {"event":"match","idx":2,"pattern":"user@hhh:\\S+\\$ ","before":"Last login: Tue Jul 14 15:59:07 2020 from 10.3.0.109\r\r\n","match":"user@hhh:~$ "}
{"execute":"ping -c 1 ggg","expect":["user@hhh:\\S+\\$ "],"timeout":2}
ping -c 1 ggg
PING ggg.ddd (10.5.0.5) 56(84) bytes of data.
64 bytes from ggg.ddd (10.5.0.5): icmp_seq=1 ttl=64 time=0.344 ms

--- ggg.ddd ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.344/0.344/0.344/0.000 ms
user@hhh:~$ {"event":"match","idx":0,"pattern":"user@hhh:\\S+\\$ ","before":"ping -c 1 ggg\r\nPING ggg.ddd (10.5.0.5) 56(84) bytes of data.\r\n64 bytes from ggg.ddd (10.5.0.5): icmp_seq=1 ttl=64 time=0.344 ms\r\n\r\n--- ggg.ddd ping statistics ---\r\n1 packets transmitted, 1 received, 0% packet loss, time 0ms\r\nrtt min/avg/max/mdev = 0.344/0.344/0.344/0.000 ms\r\n","match":"user@hhh:~$ "}
{"execute":"ping -c 77 10.3.0.99","expect":["user@hhh:\\S+\\$ "],"timeout":10}
ping -c 77 10.3.0.99
PING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.
From 10.5.0.1 icmp_seq=1 Destination Host Unreachable
From 10.5.0.1 icmp_seq=2 Destination Host Unreachable
From 10.5.0.1 icmp_seq=3 Destination Host Unreachable
From 10.5.0.1 icmp_seq=4 Destination Host Unreachable
From 10.5.0.1 icmp_seq=5 Destination Host Unreachable
From 10.5.0.1 icmp_seq=6 Destination Host Unreachable
From 10.5.0.1 icmp_seq=7 Destination Host Unreachable
From 10.5.0.1 icmp_seq=8 Destination Host Unreachable
{"event":"timeout"}
{"execute":"\u0003","expect":["user@hhh:\\S+\\$ "],"timeout":2}
^C

--- 10.3.0.99 ping statistics ---
11 packets transmitted, 0 received, +8 errors, 100% packet loss, time 10183ms
pipe 4
user@hhh:~$ 
user@hhh:~$ {"event":"match","idx":0,"pattern":"user@hhh:\\S+\\$ ","before":"ping -c 77 10.3.0.99\r\nPING 10.3.0.99 (10.3.0.99) 56(84) bytes of data.\r\nFrom 10.5.0.1 icmp_seq=1 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=2 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=3 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=4 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=5 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=6 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=7 Destination Host Unreachable\r\nFrom 10.5.0.1 icmp_seq=8 Destination Host Unreachable\r\n^C\r\n\r\n--- 10.3.0.99 ping statistics ---\r\n11 packets transmitted, 0 received, +8 errors, 100% packet loss, time 10183ms\r\npipe 4\r\nuser@hhh:~$ \r\n","match":"user@hhh:~$ "}
{"execute":"\u0004","expect":["Connection to .+ closed\\..*$"],"timeout":2}
logout
Connection to hhh.ddd closed.
{"event":"match","idx":0,"pattern":"Connection to .+ closed\\..*$","before":"logout\r\n","match":"Connection to hhh.ddd closed.\r"}
```
