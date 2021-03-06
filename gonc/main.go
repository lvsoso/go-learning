package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	verbose bool
	timeout int
	host    string
	port    string
)

const DefaultTimeout = 0

func init() {
	flag.IntVar(&timeout, "w", DefaultTimeout, "Connections which cannot be established or are idle timeout after timeout seconds.")
	flag.BoolVar(&verbose, "v", false, "Produce more verbose output.")
	flag.StringVar(&host, "h", "", "Host/Address/Doamin/IP.")
	flag.StringVar(&port, "p", "", "port.")
	flag.Parse()
}

func checkError(err error) {
	if err == nil {
		return
	}

	// output error when verbose mode is on
	if verbose {
		fmt.Fprint(os.Stderr, err)
	}
	os.Exit(1)
}

func main() {
	timeout := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if verbose {
		fmt.Printf("Succeeded to connect to %s %s port!\n", host, port)
	}

	go func() {
		io.Copy(conn, os.Stdin)
	}()

	_, err = io.Copy(os.Stdout, conn)
	checkError(err)
}
