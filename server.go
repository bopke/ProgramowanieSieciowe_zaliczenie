package main

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func getTimeMillis(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	timestampBinary := make([]byte, 8)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Read from", conn.RemoteAddr(), "failed. Closing connection.")
			return
		}
		if string(buffer[:n]) == "DISCONNECT" {
			log.Println("Caught disconnect code from", conn.RemoteAddr())
			return
		}
		if string(buffer[:n]) != "TIME" {
			// this is so sad can we hit 3 likes?
			log.Println("Incorrect query from", conn.RemoteAddr(), "'"+string(buffer[:n])+"'")
			_, err = conn.Write([]byte("ERROR"))
			if err != nil {
				log.Println("Write to", conn.RemoteAddr(), "failed. Closing connection.")
				return
			}
		}
		timestamp := getTimeMillis(time.Now())
		log.Println("Sending time to", conn.RemoteAddr(), timestamp)
		binary.LittleEndian.PutUint64(timestampBinary, uint64(timestamp))
		_, err = conn.Write(timestampBinary)
		if err != nil {
			log.Println("Write to", conn.RemoteAddr(), "failed. Closing connection.")
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	port := 8000 + rand.Intn(1000)
	ln, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Panicln("Unable to open socket", err)
	}
	log.Println("Server running on", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panicln("Unable to open connection", err)
		}
		log.Println("Opened connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
