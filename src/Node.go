package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	"os"
)

func main() {
	args := os.Args
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

	log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String()) 

	//conn.Write([]byte("ATL ipOfNode storageInGB"))
	conn.Write([]byte(args[1]))
	
	buffer := make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	parseMsg(string(buffer))
	
	defer conn.Close()
	log.Printf("Connection Killed")    
}

func parseMsg (msg string) {
	args := strings.Split(msg, " ")
	
	if args[0] == "STORE" && len(args) == 4 {
		//***This will need lots more work***
		log.Printf("VALID REQUEST")
		//Send back STORER yesOrNo
	} else if args[0] == "RETRIEVE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		//Send back RETRIEVER yesOrNo data
	} else if args[0] == "REMOVE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		//Send back REMOVER yesOrNo
	} else if args[0] == "ATLR" && len(args) == 4 {
			log.Printf("VALID RESPONSE")
	} else if args[0] == "RFLR" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
	} else if args[0] == "NODER" && len(args) == 3 {
		log.Printf("VALID RESPONSE")
	} else if args[0] == "UPDATER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
	} else if args[0] == "STORER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
	} else if args[0] == "RETRIEVER" && len(args) == 3 {
		log.Printf("VALID RESPONSE")
	} else if args[0] == "REMOVER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
	} else {
		log.Printf("NOT A VALID REQUEST OR RESPONSE")
	}
}
