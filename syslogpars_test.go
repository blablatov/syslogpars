// Test of syslog UDP server
package main

import (
	"fmt"
	"mainbeep"
	"net"
	"sync"
	"testing"
	"github.com/blablatov/syslogpars/beeper/mainbeep"
)

var strTests = []struct {
	chport  string
	sport string
}{
	{"'`", ",,`"},
	{"0001234", "45600000"},
	{"_+/__65534", "1"},
	{"NaN\null\n\n", "NaN\t\t\t123"},
	{":51444 _", ":514"},
	{"\n\123", "65534"},
	{"Number 9,78.000", "NumNum<>Num"},
}

var udpServerTests = []struct {
	snet, saddr string // server endpoint. Конечная точка сервера
	tnet, taddr string // target endpoint for client. Целевая конечная точка для клиента
	dial        bool   // test with Dial. Тест с Dial
}{
	{snet: "udp", saddr: ":0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "0.0.0.0:0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp", taddr: "::1"},

	{snet: "udp", saddr: ":0", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "0.0.0.0:0", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp", taddr: "127.0.0.1"},

	{snet: "udp", saddr: ":0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "0.0.0.0:0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp6", taddr: "::1"},

	{snet: "udp", saddr: ":0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "0.0.0.0:0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp4", taddr: "127.0.0.1"},

	{snet: "udp", saddr: "127.0.0.1:0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::ffff:127.0.0.1]:0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::1]:0", tnet: "udp", taddr: "::1"},

	{snet: "udp4", saddr: ":0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp4", saddr: "0.0.0.0:0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp4", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp4", taddr: "127.0.0.1"},

	{snet: "udp4", saddr: "127.0.0.1:0", tnet: "udp4", taddr: "127.0.0.1"},

	{snet: "udp6", saddr: ":0", tnet: "udp6", taddr: "::1"},
	{snet: "udp6", saddr: "[::]:0", tnet: "udp6", taddr: "::1"},

	{snet: "udp6", saddr: "[::1]:0", tnet: "udp6", taddr: "::1"},

	{snet: "udp", saddr: "127.0.0.1:0", tnet: "udp", taddr: "127.0.0.1", dial: true},

	{snet: "udp", saddr: "[::1]:0", tnet: "udp", taddr: "::1", dial: true},
}

func TestSyslog(t *testing.T) {

	var prevchport string
	for _, chtest := range strTests {
		if chtest.chport != prevchport {
			fmt.Printf("\n%s\n", chtest.chport)
			prevchport = chtest.chport
		}
		
		var prevsport string
	for _, ptest := range strTests {
		if ptest.sport != prevsport {
			fmt.Printf("\n%s\n", ptest.sport)
			prevsport = ptest.sport
		}

		for _, tt := range udpServerTests {
			t.Logf("skipping %s test", tt.snet+" "+tt.saddr+"<-"+tt.taddr)
			continue

			serUDPAddr, err := net.ResolveUDPAddr(tt.snet, stest.servport)
			if err != nil {
				t.Fatal(err)
			}

			for {
				sUDPConn, err := net.ListenUDP(tt.snet, serUDPAddr)
				if err != nil {
					t.Fatal(err)
				}
				handleConn(sUDPConn)
			}
		}
	}
}

func BenchmarkGoroutine(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 10; i++ {
		go mainbeep.MainBeep()
	}
}
