package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	_ "net/http/pprof"
	"sync"
	"tcp_proxy/parser"
)

func handleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// Connect to the target server
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Error connecting to target server: %v", err)
		return
	}
	defer targetConn.Close()

	// Create buffers to capture data
	var clientBuf, targetBuf bytes.Buffer

	var wg sync.WaitGroup
	defer wg.Wait()

	// we know the goroutine just have 2 so just add two in here
	wg.Add(2)

	// copy we copy target data from clientConn to target and also write into client buffer
	go func() {
		defer wg.Done()
		_, err := io.Copy(io.MultiWriter(targetConn, &clientBuf), clientConn)
		if err != nil {
			log.Printf("Error copying data from client to target: %v", err)
		}
	}()

	// after recv from target we copy client data from targetConn to client and also write into target buffer
	go func() {
		defer wg.Done()
		_, err := io.Copy(io.MultiWriter(clientConn, &targetBuf), targetConn)
		if err != nil {
			log.Printf("Error copying data from target to client: %v", err)
		}
	}()

	// process data and we need parsed into sql
	go func() {
		var buffer bytes.Buffer

		// forever loop
		for {

			if _, err := clientBuf.WriteTo(&buffer); err != nil {
				log.Printf("Error writing data from client to target to stdout: %v", err)
			}

			buf := bufio.NewScanner(&buffer)

			// with scan will blocking until found new
			for buf.Scan() {
				// cast to string
				data := string(buf.Bytes())
				// send to parser
				parser.ParseData(data)
			}
		}
	}()

}

var (
	target       = flag.String("target", "localhost:5432", "target address")
	listenerAddr = flag.String("listener", "localhost:8080", "listener address")
)

func main() {

	flag.Parse()
	// Start listening for incoming connections
	listener, err := net.Listen("tcp", *listenerAddr)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP proxy listening on %s, proxying to %s", *listenerAddr, *target)

	// Accept and handle incoming connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(clientConn, *target)
	}
}
