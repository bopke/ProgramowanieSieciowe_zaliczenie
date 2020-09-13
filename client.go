package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// Converts Unix Epoch from milliseconds to time.Time
func fromUnixMillis(ms int64) time.Time {
	return time.Unix(ms/int64(1000), (ms%int64(1000))*int64(1000000))
}

func toUnixMillis(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func handleConnection(conn net.Conn, sleeptime int, finishChannel chan bool) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		select {
		case <-finishChannel:
			_, err := conn.Write([]byte("DISCONNECT"))
			if err != nil {
				fmt.Println("Error sending request to server!", err)
			}
			return
		case <-time.After(time.Duration(sleeptime) * time.Millisecond):
			start := time.Now()
			_, err := conn.Write([]byte("TIME"))
			if err != nil {
				fmt.Println("Error sending request to server!", err)
				return
			}
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error receiving answer from server!", err)
				return
			}
			end := time.Now()
			if n == 8 {
				Tserv := int64(binary.LittleEndian.Uint64(buffer))
				delta := Tserv + (toUnixMillis(end)-toUnixMillis(start))/2 - toUnixMillis(end)
				fmt.Println("Czas serwera:", fromUnixMillis(delta+toUnixMillis(end)).Format(time.RFC3339))
				fmt.Println("Delta:", delta)
				//			log.Println("Server timestamp:", timestamp)
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	channel := make(chan bool)
	for {
		fmt.Print("Podaj adres lub wciśnij enter aby zakończyć: ")
		addr, _ := reader.ReadString('\n')
		addr = addr[:len(addr)-1]
		if len(addr) == 0 {
			return
		}

		conn, err := net.Dial("tcp4", addr)
		if err != nil {
			fmt.Println("Error connecting to server!", err)
			continue
		}

		fmt.Println("Co ile milisekund pytać serwer o czas? ")
		var sleeptime int
		_, err = fmt.Scanf("%d", &sleeptime)
		if err != nil {
			// potential stdin corruption
			fmt.Println("Not a correct number")
			return
		}
		if sleeptime < 10 || sleeptime > 1000 {
			fmt.Println("Dopuszczalny zakres 10-1000 :(")
			continue
		}

		go func() {
			fmt.Println("Wpisz cokolwiek aby zakończyć.")
			_, _ = reader.ReadString('\n')
			channel <- true

		}()

		handleConnection(conn, sleeptime, channel)
	}
}
