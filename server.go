package main

import (
	coviddata "./pkg"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)
var covidData = coviddata.Load("covid_final_data.csv")


func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,linux]")
	flag.Parse()

	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "linux":
	default:
		log.Fatalln("unsupported network protocol:", network)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer ln.Close()
	log.Printf("Service started: (%s) %s\n", network, addr)

	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\nUsage: Get <Region or Date.>\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		// buffer for client command
		cmdLine := make([]byte, 1024 * 4)
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("connection read error:", err)
			return
		}
		cmd, param := parseComand(string(cmdLine[0:n]))
		if cmd == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}

		// execute command
		switch (cmd) {
		case "Get":
			result := coviddata.Find(covidData, param)
			if len(result) == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}

			// Convert to JSON
			jsonData, err := json.Marshal(result)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			conn.Write([]byte("Response: "))
			conn.Write(jsonData)
			conn.Write([]byte("\n"))

			//Writing JSON into File
			jsonFile, err := os.Create("./covid_final_data.json")
			if err != nil {
				fmt.Println(err)
			}
			defer jsonFile.Close()

			jsonFile.Write(jsonData)
			jsonFile.Close()

		default:
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
		}
	}
}

func parseComand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}
	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])
	return
}