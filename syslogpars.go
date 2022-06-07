// syslog UDP-сервер
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
func handleConn(cn *net.UDPConn) {
	defer cn.Close()
	for {
		var cnbuf [1024]byte
		dn, addr, err := cn.ReadFromUDP(cnbuf[0:])
		if err != nil {
			continue
		}
		cntime := time.Now().String() //"15:04:05\n"
		//cntime := time.Now().Format("01/02 03:04:05PM '06 -0700")
		fmt.Println("APC client: ", string(cnbuf[0:dn]))
		alarm := string(cnbuf[0:dn])
		if strings.Contains(alarm, "Alarm") {
			var wg sync.WaitGroup // Synchronization of goroutines. Синхронизация горутин.
			wg.Add(1)             // Counter of goroutines. Значение счетчика горутин
			go mainbeep.MainBeep(wg)
			// Wait of counter. Ожидание счетчика
			go func() {
				wg.Wait()
			}()
		}
		//fmt.Println("time of server: ", cntime)
		cn.WriteToUDP([]byte(cntime), addr)
		if err != nil {
			return // For example, disabling the client. Например, отключение клиента.
		}
		//time.Sleep(1 * time.Second)
	}
}
