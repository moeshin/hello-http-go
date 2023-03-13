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

type HelloHttp struct {
	Host              string
	Port              int
	port              int
	IPv4              bool
	IPv6              bool
	AllowedMethods    map[string]bool
	DisallowedMethods map[string]bool
	listener          net.Listener
}

func (h *HelloHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.Method, r.URL)
	if !(h.AllowedMethods != nil && h.AllowedMethods[r.Method]) &&
		(h.DisallowedMethods != nil && h.DisallowedMethods[r.Method]) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf("Hello %s: %s\n", r.Method, r.RequestURI)))
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

func (h *HelloHttp) GetPort() int {
	if h.Port != 0 {
		return h.Port
	}
	return h.port
}

func (h *HelloHttp) Listen() error {
	var network string
	if h.IPv4 == h.IPv6 {
		network = "tcp"
	} else if h.IPv4 {
		network = "tcp4"
	} else {
		network = "tcp6"
	}
	listener, err := net.Listen(network, fmt.Sprintf("%s:%d", h.Host, h.Port))
	if err != nil {
		return err
	}
	h.listener = listener
	if h.Port == 0 {
		h.port = h.listener.Addr().(*net.TCPAddr).Port
	}
	return nil
}

func (h *HelloHttp) Serve() error {
	return http.Serve(h.listener, h)
}

func ParseMethods(str string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range strings.Split(str, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.ToUpper(s)
		m[s] = true
	}
	return m
}

func main() {
	host := flag.String("h", "", "Listen host. Default all host.")
	port := flag.Int("p", 8080, "Listen port. If 0, random.")
	v4 := flag.Bool("4", false, "Listen IPv4.")
	v6 := flag.Bool("6", false, "Listen IPv6.")
	allowedMethods := flag.String("m", "", "Allowed methods.")
	disallowedMethods := flag.String("d", "", "Disallowed methods.")

	flag.Parse()

	hh := HelloHttp{
		Host:              *host,
		Port:              *port,
		IPv4:              *v4,
		IPv6:              *v6,
		AllowedMethods:    ParseMethods(*allowedMethods),
		DisallowedMethods: ParseMethods(*disallowedMethods),
	}
	err := hh.Listen()
	if err != nil {
		panic(err)
	}

	h := hh.Host
	if h == "" {
		h = "all"
	}
	log.Printf("Listen host: %s, port: %d", h, hh.GetPort())

	err = hh.Serve()
	if err != nil {
		panic(err)
	}
}
