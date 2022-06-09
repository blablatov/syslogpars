// Test syslog UDP-сервер
package main

import (
	"fmt"
	"mainbeep"
	"net"
	"sync"
	"testing"
)

var strTests = []struct {
	servport string
	//SelectReqSql string
}{
	{" "},
	{"0001234"},
	{"_+/__65534"},
	{"NaN\null\n\n"},
	{":51444 _"},
	{"\n\123"},
	{"Number 9,78.000"},
}

var udpServerTests = []struct {
	snet, saddr string // server endpoint
	tnet, taddr string // target endpoint for client
	dial        bool   // test with Dial
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

	var prevservport string
	for _, stest := range strTests {
		if stest.servport != prevservport {
			fmt.Printf("\n%s\n", stest.servport)
			prevservport = stest.servport
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
		var wg sync.WaitGroup // Synchronization of goroutines. Синхронизация горутин.
		wg.Add(1)             // Counter of goroutines. Значение счетчика горутин
		go mainbeep.MainBeep(wg)
		go func() {
			wg.Wait() // Waiting of counter. Ожидание счетчика.
		}()
	}
}
