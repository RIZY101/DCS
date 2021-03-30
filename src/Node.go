package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
)

func main() {
	//6633 is node on the keypad
    cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
    if err != nil {
    	log.Fatal(err)
    }

	ip := "localhost"
	port := "6633"

	log.Printf("Connecting to %s\n", ip)

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ip+":"+port, &config)

	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte("TEST FROM NODE"))
	defer conn.Close()
	log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String())    
}
