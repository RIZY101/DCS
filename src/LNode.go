package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	//"os"
	"fmt"
	"crypto/rand"
	"net"
	"strconv"
	"sync"
	"math/big"
)
var mutex sync.Mutex

var NodeId string
var Key string
//Key2 used for retrieving files from another Node
var Key2 string
//Each nodeId is mapped to one of these structs
//This is the data you store
//A slice of local paths becuase each Data might have multiple files per NodeId
//Also local path will be data/nodeId for file
type Data struct {
	key string
	localPath string
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
	Key2 = genKey()
	//args := os.Args
	//6633 is node on the keypad
    cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
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
	log.Printf("Node Client(TLS) up and listening on port 6633")
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnection(conn)
	}

	//log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String()) 

	//conn.Write([]byte("ATL ipOfNode storageInGB"))
	//conn.Write([]byte(args[1]))
	//buffer := make([]byte, 64)
	//conn.Read(buffer)
	//log.Printf(string(buffer))
	//parseMsg(string(buffer))
	
	//defer conn.Close()
	//log.Printf("Connection Killed")    
}

func handleConnection(c net.Conn) {
    ip := c.RemoteAddr().String()
	log.Printf("Client connected! Their addr is: " + ip)
	buffer := make([]byte, 64)
	c.Read(buffer)
	msg := string(buffer)
	mutex.Lock()
	resp, size := parseMsg(msg)
	mutex.Unlock()
	c.Write([]byte(resp))
	if size > 0 {
		buffer := make([]byte, size)
		c.Read(buffer)
		data := string(buffer)
		log.Printf(data)
	}
}

func parseMsg (msg string) (string, int) {
	args := strings.Split(msg, " ")
	
	if args[0] == "STORE" && len(args) == 4 {
	    //TODO***This will need lots more work with data path and data tranfer***
	    log.Printf("VALID REQUEST")
	    mapOfData[args[1]] = Data{args[2], "data/" + args[1]}
		printMap2()
		sizeF := strings.Trim(args[3], "\x00")
		size, err := strconv.Atoi(sizeF)
		if err != nil {
			log.Fatal(err)
		}
		return "STORER yes", size
		//Send back STORER yesOrNo
	} else if args[0] == "RETRIEVE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		nodeKey := strings.Trim(args[2], "\x00")
		if mapOfData[args[1]].key == nodeKey {
			printMap2()
			//TODO Add in real size of data to expect back later 
			return "RETRIEVER yes dataSize", 0
		}
		return "RETRIEVER no 0", 0
	} else if args[0] == "REMOVE" && len(args) == 3 {
	 	log.Printf("VALID REQUEST")
	 	nodeKey := strings.Trim(args[2], "\x00")
		if mapOfData[args[1]].key == nodeKey {
			delete(mapOfData, args[1])
			printMap2()
			//TODO Add delete file on filesystem
			return "REMOVER yes", 0
		}
		return "REMOVER no", 0
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
	} else if args[0] == "CHECKR" && len(args) == 3 {
		//TODO Implement correct fileName and NodeId later
		if args[1] == "no" {
			log.Printf("No change in the IP address")
		} else {
			mapOfYourData["TestFileName"] = Node{args[2], "TestNodeId"}
			log.Printf("The IP address changed")
		}
		printMap()
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
	return "TEST RESPONSE", 0
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

//Only use this is in a mutex safe place
func printMap2() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfData {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

//crypto/rand implements a cryptographically secure random number generator so we can just use that 
//Also the key should be aprox 64 bits (This might need to get bigger later)
func genKey() string {
	key, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		log.Fatal(err)
	}
	keyStr := key.String()
	return keyStr
}

//TODO implement CHECKR
//TODO Before running retrieve run update to make sure its at the same ip
//TODO Remember to send a STORE after recieving a NODER
