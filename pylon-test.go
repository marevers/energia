package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goburrow/serial"
)

var (
	address  string
	baudrate int
	databits int
	stopbits int
	parity   string
)

func main() {
	flag.StringVar(&address, "a", "/dev/ttyACM0", "address")
	flag.IntVar(&baudrate, "b", 115200, "baud rate")
	flag.IntVar(&databits, "d", 8, "data bits")
	flag.IntVar(&stopbits, "s", 1, "stop bits")
	flag.StringVar(&parity, "p", "N", "parity (N/E/O)")
	flag.Parse()

	config := serial.Config{
		Address:  address,
		BaudRate: baudrate,
		DataBits: databits,
		StopBits: stopbits,
		Parity:   parity,
		Timeout:  30 * time.Second,
	}

	log.Printf("connecting %+v", config)
	port, err := serial.Open(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
	defer func() {
		err := port.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("closed")
	}()

	// This is the sample command from the protocol doc
	command := []byte{0x7E, 0x32, 0x30, 0x30, 0x31, 0x34, 0x36, 0x34, 0x32, 0x45, 0x30, 0x30, 0x32, 0x30, 0x31, 0x46, 0x44, 0x33, 0x35, 0x0D}
	if _, err = port.Write(command); err != nil {
		log.Fatal(err)
	}

	out, err := ioutil.TempFile("", "out485")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Output is here:", out.Name())

	go func() {
		if _, err = io.Copy(out, port); err != nil {
			log.Fatal(err)
		}
	}()

	// Ctrl+C to stop reading (sometimes EOF doesn't happen)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	close(sigChan)
	fmt.Println(sig, " stopping")

	// Close file to make sure we have flushed to disk
	_ = out.Close()
}
