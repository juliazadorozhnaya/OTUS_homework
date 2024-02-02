package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var ts time.Duration
	flag.DurationVar(&ts, "timeout", time.Second*10, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Usage: go-telnet --timeout=10s host port")
	}

	address := net.JoinHostPort(args[0], args[1])
	tClient := NewTelnetClient(address, ts, os.Stdin, os.Stdout)
	if err := tClient.Connect(); err != nil {
		log.Fatalf("Error connecting to %s: %s", address, err)
	}
	defer tClient.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go handleSignals(tClient)

	go func() {
		defer wg.Done()
		if err := tClient.Send(); err != nil {
			log.Printf("Error sending data: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := tClient.Receive(); err != nil {
			log.Printf("Error receiving data: %s", err)
		}
	}()

	wg.Wait()
}

func handleSignals(tClient TelnetClient) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT)

	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT:
			fmt.Println("\nSignal received, closing connection...")
			tClient.Close()
			os.Exit(0)
		}
	}
}
