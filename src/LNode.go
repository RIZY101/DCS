package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	"os"
	"fmt"
	"crypto/rand"
	"net"
	"strconv"
	"sync"
	//"math/big"
	"bufio"
	"io/ioutil"
)
var mutex sync.Mutex

var NodeId string
var Key string
var currentNodeId string

//Each nodeId is mapped to one of these structs
//This is the data you store
//Also local path will be data/nodeId for file
type Data struct {
	key string
	localPath string
}

//Global variables are bad dont use them...
//This is the data you store
//NodeId of who it belongs to is the key
var mapOfData map[string]Data = make(map[string]Data)


func main() {
	NodeId = ""
	Key = ""
	currentNodeId = ""
	//6633 is node on the keypad and + 1 is 6634
    cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
    if err != nil {
    	log.Fatal(err)
    }

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
	     
	listen, err := tls.Listen("tcp", "localhost:6634", &config)
	if err != nil {
		log.Fatal(err)
	}
	//6633 is node on the keypad 	
	log.Printf("Node (TLS) up and listening on port 6634")
	go userInput()
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
    ip := c.RemoteAddr().String()
	log.Printf("Client connected! Their addr is: " + ip)
	buffer := make([]byte, 128)
	c.Read(buffer)
	msg := string(buffer)
	log.Printf(msg)
	mutex.Lock()
	resp, size := parseMsg(msg)
	mutex.Unlock()
	c.Write([]byte(resp))
	if size > 0 {
		buffer2 := make([]byte, size)
		c.Read(buffer2)
		//TODO Get rid of lines bellow after testing because printing string data of png messes up cli
		//data := string(buffer2)
		//log.Printf(data)
		pwd, _ := os.Getwd()
		f, err := os.Create(pwd + "/data/" + currentNodeId)
		if err != nil {
			log.Fatal(err)
		}
		num, err := f.Write(buffer2)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d bytes", num)
	}
	if size == -1 {
		pwd, _ := os.Getwd()
		data, err := ioutil.ReadFile(pwd + "/data/" + currentNodeId)
		if err != nil {
			log.Fatal(err)
		}
		c.Write(data)
	}
}

func parseMsg (msg string) (string, int) {
	args := strings.Split(msg, " ")
	
	if args[0] == "STORE" && len(args) == 4 {
	    log.Printf("VALID REQUEST")
	    mapOfData[args[1]] = Data{args[2], "/data/" + args[1]}
		printMap()
		sizeF := strings.Trim(args[3], "\x00")
		size, err := strconv.ParseFloat(sizeF, 64)
		if err != nil {
			log.Fatal(err)
		}
		currentNodeId = args[1]
		return "STORER yes", int(size)
		//Send back STORER yesOrNo
	} else if args[0] == "RETRIEVE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		nodeKey := strings.Trim(args[2], "\x00")
		if mapOfData[args[1]].key == nodeKey {
			printMap()
			pwd, _ := os.Getwd()
			fi, err := os.Stat(pwd + "/data/" + args[1])
			if err != nil {
			    log.Fatal(err)
			}
			size := fi.Size()
			sizeStr := strconv.FormatInt(size, 10)
			return "RETRIEVER yes " + sizeStr, -1
		}
		return "RETRIEVER no 0", 0
	} else if args[0] == "REMOVE" && len(args) == 3 {
	 	log.Printf("VALID REQUEST")
	 	nodeKey := strings.Trim(args[2], "\x00")
		if mapOfData[args[1]].key == nodeKey {
			delete(mapOfData, args[1])
			printMap()
			pwd, _ := os.Getwd()
			err := os.Remove(pwd + "/data/" + args[1])
			if err != nil {
				log.Fatal(err)
			}
			return "REMOVER yes", 0
		}
		return "REMOVER no", 0
	} else if args[0] == "ATLR" && len(args) == 4 {
		log.Printf("VALID RESPONSE")
		if args[1] == "no" {
			log.Printf("Failed to add you to the network. Your IP may be blacklisted...\n")
		} else {
			NodeId = args[2]
			Key = args[3]
		}
		printGlobals()
	} else if args[0] == "RFLR" && len(args) == 2 { 
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("You were removed from the network\n")
		} else {
			log.Printf("Failed to remove you from the network. Maybe the key you sent was wrong\n")
		}
		printGlobals()
	} else if args[0] == "UPDATER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("Your IP was updated on the MasterNode\n")
		} else {
			log.Printf("Your IP did not change or you gave a bad NodeId or Key\n")
		}
	} else {
		log.Printf("NOT A VALID REQUEST OR RESPONSE")
	}
	return "TEST RESPONSE", 0
}

//Use for testing only
func printGlobals () {
	log.Printf(NodeId + "\n")
	log.Printf(Key + "\n")
}

//Only use this is in a mutex safe place
func printMap() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfData {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

func userInput() {
	fmt.Println("Type Commands Below")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
	    text = strings.Replace(text, "\n", "", -1)
	    sendCmd(text)
	}
}
func sendCmd(str string) {
	slice := strings.Split(str, " ")
	if slice[0] == "ATL" || slice[0] == "RFL" || slice[0] == "UPDATE" {
		cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
		if err != nil {
			log.Fatal(err)
		}
		ip := "localhost"
		port := "6633"
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		conn, err := tls.Dial("tcp", ip+":"+port, &config)
		if err != nil {
   			log.Fatal(err)
		}
		conn.Write([]byte(str))
		buffer := make([]byte, 64)
		conn.Read(buffer)
		parseMsg(string(buffer))
		defer conn.Close()
	} else {
		fmt.Println("Not a valid command")
	}	
}

