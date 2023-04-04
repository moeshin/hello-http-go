package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
)

var allowedMethods, disallowedMethods map[string]bool

const formatMethods = "(format: <method>[,<method>...])"

type HttpResponseWriter struct {
	http.ResponseWriter
}

func (w *HttpResponseWriter) WriteNewLine() error {
	_, err := w.Write([]byte{'\r', '\n'})
	return err
}

func (w *HttpResponseWriter) WriteString(s string) error {
	_, err := w.Write([]byte(s))
	return err
}

func (w *HttpResponseWriter) WriteHeaderLine(name string, value string) error {
	err := w.WriteString(name)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{':', ' '})
	err = w.WriteString(value)
	if err != nil {
		return err
	}
	return w.WriteNewLine()
}

func handleRequest(w *HttpResponseWriter, r *http.Request) error {
	startLine := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)
	log.Println(startLine)
	if (disallowedMethods != nil && disallowedMethods[r.Method]) ||
		(allowedMethods != nil && !allowedMethods[r.Method]) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	if r.Method == "HEAD" {
		return nil
	}

	err := w.WriteString("Hello HTTP\n\n")
	if err != nil {
		return err
	}

	err = w.WriteString(startLine)
	if err != nil {
		return err
	}
	err = w.WriteNewLine()
	if err != nil {
		return err
	}

	if r.ProtoAtLeast(1, 1) {
		err = w.WriteHeaderLine("Host", r.Host)
		if err != nil {
			return err
		}
	}

	var names []string
	for name := range r.Header {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		values := r.Header[name]
		for _, value := range values {
			err = w.WriteHeaderLine(name, value)
			if err != nil {
				return err
			}
		}
	}

	err = w.WriteNewLine()
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r.Body)
	return err
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	err := handleRequest(&HttpResponseWriter{w}, r)
	if err != nil {
		log.Println("Error handle request:", err)
	}
}

func parseMethods(str string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range strings.Split(str, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.ToUpper(s)
		m[s] = true
	}
	if len(m) > 0 {
		return m
	}
	return nil
}

func main() {
	var host, aAllowedMethods, aDisallowedMethods string
	var port int
	flag.StringVar(&host, "h", "127.0.0.1", `Listen host.
If 0.0.0.0 will only listen all IPv4.
If [::] will only listen all IPv6.
If :: will listen all IPv4 and IPv6.
`)
	flag.IntVar(&port, "p", 8080, `Listen port.
If 0, random.
`)
	flag.StringVar(&aAllowedMethods, "m", "", "Allowed methods.\n "+formatMethods)
	flag.StringVar(&aDisallowedMethods, "d", "", "Disallowed methods.\n "+formatMethods)

	flag.Parse()

	allowedMethods = parseMethods(aAllowedMethods)
	disallowedMethods = parseMethods(aDisallowedMethods)

	var network string
	switch host {
	case "0.0.0.0":
		network = "tcp4"
	case "[::]":
		network = "tcp6"
	case "::":
		host = ""
		fallthrough
	default:
		network = "tcp"
	}

	listener, err := net.Listen(network, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	addr := listener.Addr()
	log.Println("Listening", addr.Network(), addr.String())
	if port == 0 {
		port = listener.Addr().(*net.TCPAddr).Port
	}

	http.HandleFunc("/", httpHandler)
	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
