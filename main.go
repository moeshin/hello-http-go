package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
)

var allowedMethods, disallowedMethods map[string]bool

func handle(w http.ResponseWriter, r *http.Request) {
	startLine := fmt.Sprintf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	log.Print(startLine)
	if (disallowedMethods != nil && disallowedMethods[r.Method]) ||
		(allowedMethods != nil && !allowedMethods[r.Method]) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("Hello: " + startLine))
	var names []string
	for name := range r.Header {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		values := r.Header[name]
		for _, value := range values {
			w.Header()
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s\n", name, value)))
		}
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
	var v4, v6 bool
	flag.StringVar(&host, "h", "127.0.0.1", "Listen host.")
	flag.IntVar(&port, "p", 8080, "Listen port. If 0, random.")
	flag.BoolVar(&v4, "4", false, "Listen all IPv4.")
	flag.BoolVar(&v6, "6", false, "Listen all IPv6.")
	flag.StringVar(&aAllowedMethods, "m", "", "Allowed methods.")
	flag.StringVar(&aDisallowedMethods, "d", "", "Disallowed methods.")

	flag.Parse()

	allowedMethods = parseMethods(aAllowedMethods)
	disallowedMethods = parseMethods(aDisallowedMethods)

	var network string
	if v4 == v6 {
		network = "tcp"
	} else if v4 {
		network = "tcp4"
	} else {
		network = "tcp6"
	}
	if v4 || v6 {
		host = ""
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
