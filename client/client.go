package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var localAddress = ""

func handleServerConnection() {
	serverAddr := "localhost:12345"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("LIST\n"))
	if err != nil {
		log.Fatal(err)
	}

	reply := make([]byte, 1024)

	length, err := conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}
	sreply := string(reply[:length])

	localAddress = conn.LocalAddr().String()

	fmt.Println(sreply)
	fmt.Println("myIP: ", localAddress)
	var newConn net.Conn
	if sreply == "EMPTY\n" {
		signOnList(conn)
		newConn = waitForOtherClient()
		conn.Close()
		startChat(newConn)
	} else {

		newConn = connectToClient(sreply)
		conn.Close()
		startChat(newConn)
	}
}

func giveInputOptions(sreply string) string {
	replyLines := strings.Split(sreply, "\n")
}

func signOnList(conn net.Conn) {
	fmt.Println("sign on")
	_, err := conn.Write([]byte("SIGNON\n"))
	if err != nil {
		log.Fatal(err)
	}
}

func waitForOtherClient() net.Conn {
	fmt.Println("start listening as", localAddress)
	ln, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("on list client")
	fmt.Println("client-local: ", conn.LocalAddr().String())
	fmt.Println("client-remote: ", conn.RemoteAddr().String())
	return conn
}

func connectToClient(sreply string) net.Conn {
	_, clientIP, found := strings.Cut(sreply, "- ")
	if clientIP[len(clientIP)-2] == '\r' {
		clientIP = strings.TrimSuffix(clientIP, "\r")
	} else {
		clientIP = strings.TrimSuffix(clientIP, "\n")
	}
	fmt.Println("clientIP", clientIP)
	if !found {
		fmt.Println("did not work...")
		return nil
	} else {
		conn, err := net.Dial("tcp", clientIP)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("without list client")
		fmt.Println("client-local: ", conn.LocalAddr().String())
		fmt.Println("client-remote: ", conn.RemoteAddr().String())
		return conn
	}
}

func listenToClient(conn net.Conn) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Client: %s", msg)

}

func messageToClient(conn net.Conn) {
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte(text))

}

func startChat(conn net.Conn) {
	fmt.Println("start Chat")

	for {
		go messageToClient(conn)
		go listenToClient(conn)
	}
}

func main() {
	handleServerConnection()
}
