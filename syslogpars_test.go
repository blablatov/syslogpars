// syslog UDP-сервер
package main

import (
	//"os"
	"testing"
)

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

func TestUDPServer(t *testing.T) {
	for i, tt := range udpServerTests {
		if !main() {
			t.Logf("skipping %s test", tt.snet+" "+tt.saddr+"<-"+tt.taddr)
			continue
		}

		c1, err := main()
		if err != nil {
			if perr := err(err); perr != nil {
				t.Error(perr)
			}
			t.Fatal(err)
		}

		ls, err := (&packetListener{PacketConn: c1}).newLocalServer()
		if err != nil {
			t.Fatal(err)
		}
		defer ls.teardown()
		tpch := make(chan error, 1)
		handler := func(ls *sUDPConn, c UDPConn) { packetTransponder(c, tpch) }
		if err := ls.buildup(handler); err != nil {
			t.Fatal(err)
		}

		trch := make(chan error, 1)
		_, port, err := SplitHostPort(ls.PacketConn.LocalAddr().String())
		if err != nil {
			t.Fatal(err)
		}
		if tt.dial {
			d := Dialer{Timeout: someTimeout}
			c2, err := d.Dial(tt.tnet, JoinHostPort(tt.taddr, port))
			if err != nil {
				if perr := parseDialError(err); perr != nil {
					t.Error(perr)
				}
				t.Fatal(err)
			}
			defer c2.Close()
			go transceiver(c2, []byte("UDP SERVER TEST"), trch)
		} else {
			c2, err := ListenPacket(tt.tnet, JoinHostPort(tt.taddr, "0"))
			if err != nil {
				if perr := parseDialError(err); perr != nil {
					t.Error(perr)
				}
				t.Fatal(err)
			}
			defer c2.Close()
			dst, err := ResolveUDPAddr(tt.tnet, JoinHostPort(tt.taddr, port))
			if err != nil {
				t.Fatal(err)
			}
			go packetTransceiver(c2, []byte("UDP SERVER TEST"), dst, trch)
		}

		for err := range trch {
			t.Errorf("#%d: %v", i, err)
		}
		for err := range tpch {
			t.Errorf("#%d: %v", i, err)
		}
	}
}
