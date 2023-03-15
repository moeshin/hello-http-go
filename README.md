# Hello HTTP

## Usage

```shell
# Build
go build .

# Get help
./hello-http --help
```

```text
Usage of ./hello-http:
  -4    Listen all IPv4.
  -6    Listen all IPv6.
  -d string
        Disallowed methods.
  -h string
        Listen host. (default "127.0.0.1")
  -m string
        Allowed methods.
  -p int
        Listen port. If 0, random. (default 8080)
```

--------

Start

```shell
./hello-http
```

```text
Listening tcp 127.0.0.1:8080
```

Test by `curl`

```shell
curl http://127.0.0.1:8080
```

```text
Hello Http
GET / HTTP/1.1
Host: 127.0.0.1:8080
Accept: */*
User-Agent: curl/7.74.0
```

Log out

```text
GET / HTTP/1.1
```
