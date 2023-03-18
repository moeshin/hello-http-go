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

const formatMethods = "(format: <method>[,<methods>...])"

func handle(w http.ResponseWriter, r *http.Request) {
	write := func(s string) {
		_, err := w.Write([]byte(s))
		if err != nil {
			log.Println("Error writing to response", err)
		}
	}

	startLine := fmt.Sprintf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	log.Print(startLine)
	if (disallowedMethods != nil && disallowedMethods[r.Method]) ||
		(allowedMethods != nil && !allowedMethods[r.Method]) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "HEAD" {
		return
	}

	write("Hello HTTP\n\n")
	write(startLine)

	if r.ProtoAtLeast(1, 1) {
		write(fmt.Sprintf("Host: %s\n", r.Host))
	}

	var names []string
	for name := range r.Header {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		values := r.Header[name]
		for _, value := range values {
			w.Header()
			write(fmt.Sprintf("%s: %s\n", name, value))
		}
	}

	write("\n")
	_, err := io.Copy(w, r.Body)
	if err != nil {
		log.Println("Error copy request body to response", err)
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
	flag.StringVar(&host, "h", "*", `Listen host.
If 0.0.0.0 will only listen all IPv4.
If [::] will only listen all IPv6.
If * will listen all IPv4 and IPv6.
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
	case "*":
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

	http.HandleFunc("/", handle)
	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
