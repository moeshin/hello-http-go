# Hello HTTP

Build an HTTP server based on the standard library and respond to the request message.

## Usage

```text
Usage of ./hello-http:
  -h string
        Listen host.
        If 0.0.0.0 will only listen all IPv4.
        If [::] will only listen all IPv6.
        If :: will listen all IPv4 and IPv6.
         (default "*")
  -p int
        Listen port.
        If 0, random.
         (default 8080)
  -d string
        Disallowed methods.
         (format: <method>[,<method>...])
  -m string
        Allowed methods.
         (format: <method>[,<method>...])
```

Run

```shell
./hello-http
```

Hello HTTP output

```text
Listening tcp 127.0.0.1:8080
```

Test with cURL

```shell
curl http://127.0.0.1:8080
```

cURL output

```text
Hello HTTP

GET / HTTP/1.1
Host: 127.0.0.1:8080
Accept: */*
User-Agent: curl/7.74.0

```

Hello HTTP output

```text
GET / HTTP/1.1
```
