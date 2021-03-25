package main

import (
	"crypto/tls"
	"log"
)

func main() {
	//6633 is node on the keypad
	//connThread, err := net.Dial("tcp", "localhost:6633")
	 //if err != nil {
	 	//fmt.Println("Error connecting:", err.Error())
	 	//os.Exit(1)
	 //}

	 //defer connThread.Close()
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

	defer conn.Close()
	log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String())    
}
