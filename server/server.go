package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
			split_msg := strings.Split(msg, "\n")
			/* LIST request
			 * sending all the client IPs that signed on to the requesting client
			 */
			if split_msg[0] == "LIST" {
				if len(clients) == 0 {
					response = "EMPTY\n"
				} else {
					for client, clientNmbr := range clients {
						response += fmt.Sprintf("%d- %s\n", clientNmbr, client)
						fmt.Println(response)
					}
				}

				/* SIGNON request
				 * putting the requesting clients IP on the list
				 */
			} else if split_msg[0] == "SIGNON" {
				fmt.Println("signing on ")
				clients[conn.RemoteAddr().String()] = len(clients)

				/* CHOOSE request
				 * requesting client sends client IP (from list) that then gets removed from the list
				 */
			} else if split_msg[0] == "CHOOSE" {
				if len(split_msg) > 1 {
					clientIP := strings.TrimSuffix(split_msg[1], "\n")
					if _, ok := clients[clientIP]; ok {
						delete(clients, clientIP)
						response = "200"
					}
				} else {
					response = "ERROR: foreign client IP not provided\n"
				}
			} else {
				response = "400 Bad Request\n"
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
