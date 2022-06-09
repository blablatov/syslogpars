// syslog UDP-сервер.
// Base of idea to https://github.com/alash3al/go-beeper/tree/v1.0.0.
// Thanks him!
package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"mainbeep"
	"sync"
	"time"
)

func main() {
	// Starting cycle listen the udp-server. Работа udp-сервера в цикле.
	servport := ":51444"
	serUDPAddr, err := net.ResolveUDPAddr("udp", servport)
	if err != nil {
		log.Fatal(err)
	}
	for {
		sUDPConn, err := net.ListenUDP("udp", serUDPAddr)
		if err != nil {
			log.Fatal(err)
		}
		handleConn(sUDPConn)
	}
}

// Handler of connect. Обработчик подключения.
func handleConn(cn *net.UDPConn) {
	defer cn.Close()
	for {
		var cnbuf [1024]byte
		dn, addr, err := cn.ReadFromUDP(cnbuf[0:])
		if err != nil {
			continue
		}
		// System time. Время сервера.
		cntime := time.Now().String() //"15:04:05\n"
		// Data the host syslog. Syslog данные с хоста.
		fmt.Println("APC client: ", string(cnbuf[0:dn]))
		alarm := string(cnbuf[0:dn])
		// If data contains need the string, call method beeper via goroutine.
		// Если данные содержат реперную строку, вызываем метод beeper через goroutine
		if strings.Contains(alarm, "High temperature") {
			var wg sync.WaitGroup // Synchronization of goroutines. Синхронизация горутин.
			wg.Add(1)             // Counter of goroutines. Значение счетчика горутин
			go mainbeep.MainBeep(wg)
			// Wait of counter. Ожидание счетчика
			go func() {
				wg.Wait()
			}()
		}
		cn.WriteToUDP([]byte(cntime), addr)
		if err != nil {
			return // For example, disabling the client. Например, отключение клиента.
		}
	}
}
