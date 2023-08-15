// Simple syslog UDP-server.
// Base of idea to https://github.com/alash3al/go-beeper.git
// Thanks dude!

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"syslog2mongo"
	"time"

	"github.com/blablatov/syslogpars/beeper"
	"github.com/blablatov/syslogpars/syslog2mongo"
)

type embmongo struct {
	syslog2mongo.ApcMongo
}

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Getting set udp port from config file.
	// Получение заданного udp порта из конфига.
	chport := make(chan string)
	go func() {
		chport <- readPort()
	}()

	// Starting cycle listen the udp-server.
	// Старт udp-сервера в цикле прослушивания.
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

		// Data the syslog client. Syslog данные клиента.
		fmt.Println("APC client:", addr.String(), string(cnbuf[0:dn]))

		alarm := string(cnbuf[0:dn])

		// If data of client contains the string with set message, calls func beep in goroutine
		// Если данные клиента содержат строку с заданным сообщением, вызывается метод beep в goroutine
		if strings.Contains(alarm, "High temperature") {
			go mainBeep()
			log.Println("\n\nAPC client high temp:", addr.String(), alarm)
		}
		if strings.Contains(alarm, "Maximum temperature") {
			go mainBeep()
			log.Println("\n\nAPC client max temp:", addr.String(), alarm)
		}
		if strings.Contains(alarm, "Configuration") {
			log.Println("\n\nAPC client any:", addr.String(), alarm)
		}

		///////////////////////////////////////
		// Writting data of modbus to MongoDB via method SendMongo of interface.
		// Запись данных APC в MongoDB через метод SendMongo интерфейса.
		start := time.Now()

		// Formating data of structure ApcMongo. Заполнение структуры.
		var m embmongo
		m.AddrData = addr.String()
		m.AlarmData = alarm

		// Getting DSN MongoDB from config. Получение DSN из конфига.
		dsnmgo := make(chan string)
		go func() {
			dsnmgo <- readDsn()
		}()
		DsnMongo := <-dsnmgo

		var b bool
		// Вызов метода с передачей аргументов. Calls the method
		b, err = embmongo.SendMongo(m, DsnMongo)
		if err != nil {
			log.Fatalf("Error of method via type embedding: %v", err)
		}
		fmt.Println("\nResult of request via type embedding: ", b)

		secs := time.Since(start).Seconds()
		fmt.Printf("\n%.2fs Request execution time to MongoDB via type embedding\n", secs)

		////////////////////////////////////////////////////////////////////////
		// Sending data to MongoDB via goroutine.
		// Отправка данных в MongoDB через горутину.
		start2 := time.Now()

		// Formating data of structure. Заполнение структуры.
		mg := syslog2mongo.ApcMongo{
			AddrData:  addr.String(),
			AlarmData: alarm,
		}

		adr := make(chan string) // Channel to data send addr of sensor. Канал передачи адреса APC.
		alr := make(chan string) // Channel to data send alarm. Канал данных тревоги.

		var wg sync.WaitGroup // Synchronization of goroutines. Синхронизация горутин.
		wg.Add(1)             // Counter of goroutines. Значение счетчика горутин

		go syslog2mongo.SendgMongo(mg.AddrData, DsnMongo, mg.AlarmData, adr, alr, wg)

		// Getting data from goroutine. Получение данных из канала горутины.
		log.Println("\nSensor of system: ", <-adr, "\nData of sensor: ", <-alr)

		// Wait of counter. Ожидание счетчика
		go func() {
			wg.Wait() // Waiting of counter. Ожидание счетчика.
			close(adr)
			close(alr)
		}()
		secs2 := time.Since(start2).Seconds()
		fmt.Printf("\n%.2fs Request execution time to MongoDB via goroutine\n", secs2)

		cn.WriteToUDP([]byte(cntime), addr)
		if err != nil {
			return // For example, disabling the client. Например, при отключении клиента.
		}
	}
}

// Func reads udp port of syslogd from the file ./port.conf
// Метод чтения udp порта syslog-сервера из локального конфига
func readPort() string {
	var sport string
	sp, err := os.Open("port.conf")
	if err != nil {
		log.Fatalf("\nError open config port: %v", err)
	}
	defer sp.Close()
	input := bufio.NewScanner(sp)
	for input.Scan() {
		sport = input.Text()
	}
	return sport
}

// Reads DSN from file the ./mongo.conf. Получение DSN из конфига
func readDsn() string {
	var dsn string
	rf, err := os.Open("mongo.conf")
	if err != nil {
		log.Fatalf("\nError open a config mongo: %v", err)
	}
	defer rf.Close()
	input := bufio.NewScanner(rf)
	for input.Scan() {
		dsn = input.Text()
	}
	return dsn
}

func mainBeep() {
	// beep once. Подать звуковой сигнал один раз
	//beeper.Beep()

	// beep three times. Звуковой сигнал три раза
	//beeper.Beep(3)

	// beep, beep, pause, pause, beep, pause, pause, etc
	// Мелодия в цикле (*бипер, -пауза)
	beeper.Melody("**--**--**--**")
	time.Sleep(5 * time.Second)
}
