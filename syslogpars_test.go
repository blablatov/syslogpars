// Test of syslog UDP server
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/blablatov/syslogpars/beeper"
)

var strTests = []struct {
	chport string
	sport  string
}{
	{"'`", ",,`"},
	{"0001234", "45600000"},
	{"_+/__65534", "1"},
	{"NaN\null;%;№$", "-=90++123"},
	{":51444 _", ":514"},
	{"000_///123", "655349986"},
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
	{snet: "udp", saddr: ":514", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "192.168.1.1:51444", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp", taddr: "::1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp", taddr: "127.0.0.1"},
	{snet: "udp", saddr: ":0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "10.0.0.2:51444", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: ":0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "0.0.0.0:0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "[::ffff:0.0.0.0]:0", tnet: "udp6", taddr: "::1"},
	{snet: "udp", saddr: "[::]:0", tnet: "udp4", taddr: "127.0.0.1"},
	{snet: "udp", saddr: "127.0.0.1:514", tnet: "udp", taddr: "127.0.0.1"},
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

				serUDPAddr, err := net.ResolveUDPAddr(tt.snet, ptest.chport)
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
}

func TestReadPort(t *testing.T) {
	var sport string
	sp, err := os.Open("port.conf")
	if err != nil {
		log.Fatalf("Error open config port: %v", err)
	}
	defer sp.Close()
	input := bufio.NewScanner(sp)
	for input.Scan() {
		sport = input.Text()
	}
	if sport != "" {
		log.Println("Test get port from config file is ok: ", sport)
		switch strings.HasPrefix(sport, ":") {
		case false:
			log.Println("Need prefix ':' before number of port in the file of config")
		case true:
			log.Println("Test of the prefix ':' is ok")
		default:
			log.Println("Check prefix ':' before number of port in the file of config")
		}
	}
}

func TestReadDsn(t *testing.T) {
	var dsn string
	sm, err := os.Open("mongo.conf")
	if err != nil {
		log.Fatalf("Error open config mongo: %v", err)
	}
	defer sm.Close()
	input := bufio.NewScanner(sm)
	for input.Scan() {
		dsn = input.Text()
	}
	if dsn != "" {
		log.Println("Test get dsn from config file is ok: ", dsn)
	}
}

func TestEOF(t *testing.T) {
	var atEOF bool
	alarms := []byte("System")
	_, _, err := bufio.ScanLines(alarms, atEOF)
	if err == nil {
		log.Printf("EOF is to string of data: %v", err)
	}
}

func TestBeeper(t *testing.T) {
	beeper.Melody("***---***")
}

func BenchmarkBeeper(b *testing.B) {
	b.ReportAllocs()
	b.ReportAllocs()
	for i := 0; i < 10; i++ {
		go mainBeep()
	}
}

func BenchmarkMelody(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 10; i++ {
		beeper.Melody()
	}
}
