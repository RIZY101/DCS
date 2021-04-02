package main

//@Author: Richard Zins

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"log"
	"strings"
)

func main() {

	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")

	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
     
	listen, err := tls.Listen("tcp", "localhost:6633", &config)
	if err != nil {
		log.Fatal(err)
	}
	//6633 is node on the keypad 	
	log.Printf("Server(TLS) up and listening on port 6633")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Printf("Client connected! Their addr is: " + c.RemoteAddr().String())
	buffer := make([]byte, 32)
	c.Read(buffer)
	msg := string(buffer)
	resp := parseMsg(msg)
	c.Write([]byte(resp))
}

func parseMsg (msg string) string {
	args := strings.Split(msg, " ")
	
	if args[0] == "ATL" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		return "ATLR yesOrNo nodeId key"
	} else if args[0] == "RFL" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		return "RFLR yesOrNo"
	} else if args[0] == "NODE" && len(args) == 1 {
		log.Printf("VALID REQUEST")
		return "NODER ipOfNewNode nodeId\n"
	} else if args[0] == "UPDATE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		return "UPDATER yesOrNO"
	} else {
		log.Printf("NOT A VALID REQUEST")
	}
	return "NOT A VALID REQUEST"
}


