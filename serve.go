// Serve a directory over HTTP, with a twist: prints a QR code for your mobile testing needs.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/mdp/qrterminal/v3"
)

// Get preferred outbound ip of this machine
// From https://stackoverflow.com/a/37382208/17901280
func outboundIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

var qrConfig = qrterminal.Config{
	Level:     qrterminal.M,
	Writer:    os.Stdout,
	BlackChar: qrterminal.WHITE,
	WhiteChar: qrterminal.BLACK,
	QuietZone: 1,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Serve a directory over HTTP, with a twist: prints a QR code for your mobile testing needs.\n\nUsage:\n")
		flag.PrintDefaults()
	}
	port := flag.String("p", "7999", "port")
	directory := flag.String("d", ".", "directory")
	hideQr := flag.Bool("q", false, "don't print qr code")
	dontOpenCurrentPage := flag.Bool("w", false, "don't try to open url in your local browser")
	flag.Parse()

	url := fmt.Sprintf("http://%s:%s", outboundIp(), *port)

	if !*hideQr {
		qrterminal.GenerateWithConfig(url, qrConfig)
	}

	log.Printf("Serving %s on %s\n", *directory, url)
	if !*dontOpenCurrentPage {
		time.AfterFunc(500*time.Millisecond, func() {
			cmd := exec.Command("open", url)
			if err := cmd.Start(); err != nil {
				log.Println("Tried everything (open) but couldn't open your browser")
			}
		})
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
