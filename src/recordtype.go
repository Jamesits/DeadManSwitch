package main

import (
	"github.com/miekg/dns"
	"strings"
)

func stringToDnsType(s string) uint16 {
	s = strings.ToUpper(s)

	switch s {
	case "A":
		return dns.TypeA
	case "AAAA":
		return dns.TypeAAAA
	case "TXT":
		return dns.TypeTXT
	case "MX":
		return dns.TypeMX
	case "SRV":
		return dns.TypeSRV
	}

	// default
	return dns.TypeA
}