package main

import (
	"flag"
	"fmt"
	"net"
)

// Types of network
const (
	NetIP string = "ip+net"
)

var (
	IPv4bcast     = net.IPv4(255, 255, 255, 255) // broadcast
	IPv4allsys    = net.IPv4(224, 0, 0, 1)       // all systems
	IPv4allrouter = net.IPv4(224, 0, 0, 2)       // all routers
	IPv4zero      = net.IPv4(0, 0, 0, 0)         // all zeros
)

var (
	IPv6zero                   = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	IPv6unspecified            = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	IPv6loopback               = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	IPv6interfacelocalallnodes = net.IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
	IPv6linklocalallnodes      = net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}
	IPv6linklocalallrouters    = net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}
)

// Basic network information about this machine
var interfaces []net.Interface
var interfaceAddrs []net.Addr

var Options struct {
	listInterfaces     bool
	listInterfaceAddrs bool
	rdns               string
	ip                 string
	mx                 string
	ns                 string
	txt                string
}

func main() {
	var err error
	cmdFlags()
	flag.Parse()
	if interfaces, err = net.Interfaces(); err != nil {
		panic(err)
	}
	if interfaceAddrs, err = net.InterfaceAddrs(); err != nil {
		panic(err)
	}
	if Options.listInterfaces {
		listInterfaces(interfaces)
	} else if Options.listInterfaceAddrs {
		listInterfaceAddrs(interfaceAddrs)
	} else if Options.rdns != "" {
		hostnames, _ := net.LookupAddr(Options.rdns)
		for _, hostname := range hostnames {
			fmt.Println(hostname)
		}
	} else if Options.ip != "" {
		ips, _ := net.LookupIP(Options.ip)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	} else if Options.mx != "" {
		mxs, _ := net.LookupMX(Options.mx)
		for _, mx := range mxs {
			fmt.Println(mx.Pref, mx.Host)
		}
	} else if Options.ns != "" {
		nss, _ := net.LookupNS(Options.ns)
		for _, ns := range nss {
			fmt.Println(ns.Host)
		}
	} else if Options.txt != "" {
		txts, _ := net.LookupTXT(Options.txt)
		for _, txt := range txts {
			fmt.Println(txt)
		}
	}
}

func listInterfaces(interfaces []net.Interface) {
	for _, x := range interfaces {
		fmt.Printf(
			"%-15v %-6v %-20v %v\n",
			x.Name, x.MTU, x.HardwareAddr, x.Flags,
		)
	}
}

func listInterfaceAddrs(interfaceAddrs []net.Addr) {
	for _, x := range interfaceAddrs {
		fmt.Println(x.Network(), x.String())
	}
}

func cmdFlags() {
	flag.BoolVar(&Options.listInterfaces, "if", false,
		"list interfaces on this machine")
	flag.BoolVar(&Options.listInterfaceAddrs, "ifaddr", false,
		"list interface addrs on this machine")
	flag.StringVar(&Options.rdns, "rdns", "",
		"Reverse Lookup of ip address to name")
	flag.StringVar(&Options.ip, "ip", "",
		"Lookup of name address to ip address")
	flag.StringVar(&Options.mx, "mx", "",
		"Lookup of mx record for a domain name")
	flag.StringVar(&Options.ns, "ns", "",
		"Lookup of ns record for a domain name")
	flag.StringVar(&Options.txt, "txt", "",
		"Lookup of txt record for a domain name")
	//flag.IntVar( &options.seed, "s", seed, "Seed value" )
	//flag.IntVar( &options.count, "n", 1, "Generate n combinations" )
	//flag.BoolVar( &options.help, "h", false, "Print usage and default options" )
}
