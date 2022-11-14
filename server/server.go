package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// var clients = make([]string, 0)
var clients = make(map[string]int)

func handleConnection(conn net.Conn) {
	fmt.Println("new conn", conn.RemoteAddr())
	for {
		//fmt.Println(conn.RemoteAddr().String())
		response := ""
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			if msg == "LIST\n" {
				if len(clients) == 0 {
					response = "EMPTY\n"
				} else {
					for client, clientNmbr := range clients {
						response += fmt.Sprintf("%d- %s\n", clientNmbr, client)
						fmt.Println(response)
					}
				}
			} else if msg == "SIGNON\n" {
				fmt.Println("signing on ")
				clients[conn.RemoteAddr().String()] = len(clients)
			} else if msg == "" {
				// TODO
				// client sends number to choose
			}
			conn.Write([]byte(response))
		}
		if err == io.EOF {
			fmt.Println("connection to ", conn.RemoteAddr().String(), " has been closed")
			delete(clients, conn.RemoteAddr().String())
			conn.Close()
			return
		}
		msg = ""
	}
}

func main() {
	fmt.Println("starting server...")
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
