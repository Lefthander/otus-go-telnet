## Lesson - 20
### Low level protocols TCP,UDP, DNS



Develop a basic telnet client.

Examples:

```shell
$ go-telnet --timeout=10s host port

$ go-telnet mysite.ru 8080

$ go-telnet --timeout=3s 1.1.1.1 123
```


The client must connect to the specified host ( ip address or domain name) & port via the TCP protocol

When the client is getting connected the STDIN of client must be directed to the socket, and all data received from the 
socket must be out to the STDOUT

Optionaly client can receive a timeout command argument (default value 10s, when the parameter omited)

Client must catch a Ctrl-D and does close the socket and exit.
In case the socket is closed by the server side, the client should exit.

In case of connection to the unreachable host, the client should be terminated after timeout.

