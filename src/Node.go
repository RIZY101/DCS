package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	"os"
	"fmt"
)

var NodeId string
var Key string

//Each nodeId is mapped to one of these structs
//This is the data you store
//A slice of local paths becuase each Data might have multiple files per NodeId
type Data struct {
	key string
	localPath []string
}

//Global variables are bad dont use them...
//This is the data you store
//NodeId of who it belongs to is the key
var mapOfData map[string]Data = make(map[string]Data)

type Node struct {
	ip string
	NodeId string
}
//Global variables are bad dont use them...
//This is where your data is stored
//Each file name is the key
var mapOfYourData map[string]Node = make(map[string]Node)


func main() {
	NodeId = ""
	Key = ""
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
	    //TODO***This will need lots more work***
		log.Printf("VALID REQUEST")
		//Send back STORER yesOrNo
	} else if args[0] == "RETRIEVE" && len(args) == 3 {
		//TODO
		log.Printf("VALID REQUEST")
		//Send back RETRIEVER yesOrNo data
	} else if args[0] == "REMOVE" && len(args) == 3 {
	 	//TODO
		log.Printf("VALID REQUEST")
		//Send back REMOVER yesOrNo
	} else if args[0] == "ATLR" && len(args) == 4 {
		if args[1] == "no" {
			log.Printf("Failed to add you to the network. Your IP may be blacklisted...\n")
		} else {
			NodeId = args[2]
			Key = args[3]
		}
		printGlobals()
		log.Printf("VALID RESPONSE")
	} else if args[0] == "RFLR" && len(args) == 2 { 
		if args[1] == "no" {
			log.Printf("Failed to remove you from the network. Maybe the key you sent was wrong\n")
		} else {
			log.Printf("You were removed from the network\n")
		}
		printGlobals()
		log.Printf("VALID RESPONSE")
	} else if args[0] == "NODER" && len(args) == 3 {
		//TODO implment adding correct file name later
		mapOfYourData["TestFileName"] = Node{args[1], args[2]}
		log.Printf("VALID RESPONSE")
		printMap()
	} else if args[0] == "UPDATER" && len(args) == 2 {
		if args[1] == "yes" {
			log.Printf("Your IP was updated on the MasterNode\n")
		} else {
			log.Printf("Your IP did not change or you gave a bad NodeId or Key\n")
		}
		log.Printf("VALID RESPONSE")
	} else if args[0] == "STORER" && len(args) == 2 {
		//TODO if no maybe retransmit
		log.Printf("VALID RESPONSE")
	} else if args[0] == "RETRIEVER" && len(args) == 3 {
		//TODO if no tell the user maybe their key was invalid
		log.Printf("VALID RESPONSE")
	} else if args[0] == "REMOVER" && len(args) == 2 {
		//TODO If no tell user maybe key was wrong
		log.Printf("VALID RESPONSE")
	} else {
		log.Printf("NOT A VALID REQUEST OR RESPONSE")
	}
}

//Usef for testing only
func printGlobals () {
	log.Printf(NodeId + "\n")
	log.Printf(Key + "\n")
}

//Only use this is in a mutex safe place
func printMap() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfYourData {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

//TODO implement CHECKR
//TODO Before running retrieve run update to make sure its at the same ip
//TODO Remember to send a STORE after recieving a NODER
