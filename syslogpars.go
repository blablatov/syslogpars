// Simple syslog UDP-server.
// Base of idea to https://github.com/alash3al/go-beeper.git
// Thanks dude!

package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/blablatov/syslogpars/beeper"
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Getting set udp port from config file.
	// Получение заданного udp порта из конфига.
	chport := make(chan string)
	go func() {
		chport <- readIp()
	}()

	// Starting cycle listen the udp-server.
	// Старт udp-сервера в цикле прослушивания.
	//servport := ":51444"
	saddr, err := net.ResolveUDPAddr("udp", <-chport)
	if err != nil {
		log.Fatal(err)
	}
	for {
		slis, err := net.ListenUDP("udp", saddr)
		if err != nil {
			log.Fatal(err)
		}
		handleConn(slis)
	}
}

// Handler of connect syslog clients.
// Обработчик подключения syslog-клиентов.
func handleConn(cn *net.UDPConn) {
	defer cn.Close()
	for {
		var cnbuf [1024]byte
		dn, addr, err := cn.ReadFromUDP(cnbuf[0:])
		if err != nil {
			continue
		}

		// System time. Время сервера.
		cntime := time.Now().String()

		// Data the syslog client. Syslog данные с клиента.
		fmt.Println("APC client: ", string(cnbuf[0:dn]))

		alarm := string(cnbuf[0:dn])

		// If data of client contains the string with some message, calls func beep in goroutine
		// Если данные клиента содержат строку с каким-либо сообщением, вызывается метод beep в goroutine
		switch alarm {
		case "High temperature":
			go mainbeep.MainBeep()
			log.Println("APC client: ", alarm)
		case "Maximum temperature":
			go mainbeep.MainBeep()
			log.Println("APC client: ", alarm)
		case "System:":
			fmt.Println("APC client: %v", alarm)
			log.Println("APC client: ", alarm)
		case "Configuration":
			fmt.Println("APC client: %v", alarm)
		default:
			Println("APC client any: ", alarm)
		}

		cn.WriteToUDP([]byte(cntime), addr)
		if err != nil {
			return // For example, disabling the client. Например, при отключении клиента.
		}
	}
}

// Func reads udp port of syslogd from the file ./port.conf
// Метод чтения udp порта syslog-сервера из локального конфига
func readIp() string {
	var sport string
	sp, err := os.Open("port.conf")
	if err != nil {
		log.Fatalf("Error open config: %v", err)
	}
	defer sp.Close()
	input := bufio.NewScanner(sp)
	for input.Scan() {
		sport = input.Text()
	}
	return sport
}
