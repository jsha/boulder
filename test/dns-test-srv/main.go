// Copyright 2015 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/miekg/dns"
)

var listen = flag.String("listen", ":8053", "port (and optionally address) to listen on")

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	// Normally this test DNS server will return 127.0.0.1 for everything.
	// However, in some situations (for instance Docker), it's useful to return a
	// different hardcoded host. You can do so by setting the FAKE_DNS environment
	// variable.
	fakeDNS := os.Getenv("FAKE_DNS")
	if fakeDNS == "" {
		fakeDNS = "127.0.0.1"
	}
	for _, q := range r.Question {
		if strings.Contains(q.Name, "sleep") {
			index := strings.Index(q.Name, ".")
			sleepTime, err := strconv.Atoi(q.Name[0:index])
			if err == nil {
				time.Sleep(time.Duration(sleepTime) * time.Second)
			} else {
				fmt.Printf("Parse error: %s", err)
			}
		}
		fmt.Printf("dns-srv: Query -- [%s] %s\n", q.Name, dns.TypeToString[q.Qtype])
		switch q.Qtype {
		case dns.TypeA:
			record := new(dns.A)
			record.Hdr = dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    86400,
			}
			record.A = net.ParseIP(fakeDNS)

			m.Answer = append(m.Answer, record)
		}
	}

	w.WriteMsg(m)
	return
}

func serveTestResolver() {
	dns.HandleFunc(".", dnsHandler)
	server := &dns.Server{
		Addr:         *listen,
		Net:          "udp",
		ReadTimeout:  time.Millisecond,
		WriteTimeout: time.Millisecond,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}

func main() {
	flag.Parse()
	fmt.Println("dns-srv: Starting test DNS server on", *listen)
	serveTestResolver()
	forever := make(chan bool, 1)
	<-forever
}
