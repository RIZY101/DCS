package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	"os"
	"os/exec"
	"fmt"
	"crypto/rand"
	"strconv"
	"math/big"
	"io/ioutil"
	"net"
	"bufio"
	//"io"
)

var NodeId string
//Your key to give to the LNode storing your file
var Key string
var fileName string
var currentNodeId string

type Node struct {
	ip string
	NodeId string
}
//Global variables are bad dont use them...
//This is where your data is stored
//Each file name is the key
var mapOfYourData map[string]Node = make(map[string]Node)


func main() {
	fileName = ""
	NodeId = strings.Trim(genNodeId(), " ")
	NodeId = strings.Trim(NodeId, "\x0a")
	currentNodeId = ""
	Key = genKey()
    
    fmt.Println("Welcome to the DCS Client CLI!\nPlease type HELP for a list of commands.")
    reader := bufio.NewReader(os.Stdin)
    for {
    	fmt.Printf("> ")
    	text, _ := reader.ReadString('\n')
    	// convert CRLF to LF
    	text = strings.Replace(text, "\n", "", -1)
    	parseCmd(text)
    } 
}

func parseMsg (msg string, conn net.Conn) int {
	args := strings.Split(msg, " ")
		//Send back STORER yesOrNo
	if args[0] == "NODER" && len(args) == 3 {
		str := strings.Split(args[1], ":")
		ip := str[0]
		mapOfYourData[fileName] = Node{ip, args[2]}
		log.Printf("VALID RESPONSE")
		//printMap()
	} else if args[0] == "CHECKR" && len(args) == 3 {
		if args[1] == "no" {
			log.Printf("No change in the IP address")
		} else {
			str := strings.Split(args[2], ":")
			ip := str[0]
			mapOfYourData[fileName] = Node{ip, currentNodeId}
			log.Printf("The IP address changed")
		}
		//printMap()
	} else if args[0] == "STORER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("LNode ready for the data")
		} else {
			log.Printf("LNode not ready for your data")
		}
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(data)
	} else if args[0] == "RETRIEVER" && len(args) == 3 {
		log.Printf("VALID RESPONSE")
		sizeF := strings.Trim(args[2], "\x00")
		size, err := strconv.Atoi(sizeF)
		if err != nil {
			log.Fatal(err)
		}
		return size
	} else if args[0] == "REMOVER" && len(args) == 2 {
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("Data Removed")
		} else {
			log.Printf("Failed to remove your data. Maybe the key you sent was wrong.\n")
		}
	} else {
		log.Printf("NOT A VALID REQUEST OR RESPONSE")
	}
	return 0
}

func parseCmd(cmd string) {
	cmds := strings.Split(cmd, " ")
	if cmds[0] == "STORE" {
		fileName = cmds[1]
		store()
	} else if cmds[0] == "RETRIEVE" {
		fileName = cmds[1]
		retrieve()	
	} else if cmds[0] == "REMOVE" {
		fileName = cmds[1]
		remove()
	} else if cmds[0] == "FILES" {
		files()
	} else if cmds[0] == "MAP" {
		mapp()
	} else if cmds[0] == "HELP" {
		fmt.Println("STORE <filename>: Stores the file you specify on the network.\nRETRIEVE <filename>: Retrieves the file you specify from the network.\nREMOVE <filename>: Removes that file from the network.\nFILES: Lists all your files stored on the network.\nMAP: Shows the physical location of where all your files are stored on the network.\nHELP: Shows this list of commands.")
	} else {
		fmt.Println("Not a valid command! Please try using the HELP command to see all possible commands.")
	}
}

func store() {
	//Conn setup
	//6633 is node on the keypad
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatal(err)
	}
	ip := "localhost"
	//For MNode
	port := "6633"
	//log.Printf("Connecting to %s\n", ip)
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ip+":"+port, &config)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String())

	//Ask MasterNode for new Node
	info, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fSize := info.Size()
	sSize := strconv.FormatInt(fSize, 10)
	conn.Write([]byte("NODE " + sSize))
	buffer := make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	size := parseMsg(string(buffer), conn)
	//log.Printf("%d",size)
	defer conn.Close()
	//log.Printf("Connection to Master Node Killed")

	//Ask Node to store your stuff there
	//TODO Make it based off struct later
	//ip = 0
	//For LNode
	port = "6634"
	conn, err = tls.Dial("tcp", ip+":"+port, &config)
	if err != nil {
		log.Fatal(err)
	}
	writeStr := "STORE " + NodeId + " " + Key + " " + sSize
	conn.Write([]byte(writeStr))
	buffer = make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	size = parseMsg(string(buffer), conn)
	//log.Printf("%d", size)
	defer conn.Close()
	//log.Printf("Connection to Node Killed")
	_ = size
}

func retrieve() {
	//Conn setup
	//6633 is node on the keypad
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatal(err)
	}
	ip := "localhost"
	//For LNode
	port := "6634"
	//log.Printf("Connecting to %s\n", ip)
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ip+":"+port, &config)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String())
	
	//Ask MasterNode for new Node
	conn.Write([]byte("RETRIEVE " + NodeId + " " + Key))
	buffer := make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	size := parseMsg(string(buffer), conn)
	if size > 0 {
		buffer2 := make([]byte, size)
		conn.Read(buffer2)
		pwd, _ := os.Getwd()
		f, err := os.Create(pwd + "/data2/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		num, err := f.Write(buffer2)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d bytes", num)
	}
	defer conn.Close()
	//log.Printf("Connection to Node Killed")
}

func remove() {
	//Conn setup
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		log.Fatal(err)
	}
	ip := "localhost"
	//For LNode
	port := "6634"
	//log.Printf("Connecting to %s\n", ip)
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ip+":"+port, &config)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String())
	
	//Ask MasterNode for new Node
	conn.Write([]byte("REMOVE " + NodeId + " " + Key))
	buffer := make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	size := parseMsg(string(buffer), conn)
	//log.Printf("%d", size)
	defer conn.Close()
	//log.Printf("Connection to Node Killed")
	_ = size
}

func files() {
	for k := range mapOfYourData {
		fmt.Printf("%s\n", k)
	}
}

func mapp() {
	
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

//Please note this function uses a unix utility to generate a standard uuid
//You can learn more about uuid's here: https://en.wikipedia.org/wiki/Universally_unique_identifier
func genNodeId() string {
	nodeId, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(nodeId)
}
