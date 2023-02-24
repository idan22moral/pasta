package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/idan22moral/pasta/internal"
	"github.com/idan22moral/pasta/internal/server"
)

func getMyPrimaryIP() (string, error) {
	const googleDNSAddr string = "8.8.8.8:53"
	googleAddr, err := net.ResolveTCPAddr("tcp", googleDNSAddr)
	if err != nil {
		return "", err
	}
	conn, err := net.DialTCP("tcp", nil, googleAddr)
	if err != nil {
		return "", err
	}
	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP, nil
}

func main() {
	const port string = "8080"
	const serverIP string = "0.0.0.0"
	const uploadsDir string = "./uploads/"

	exitSignal := make(chan interface{})

	go func() {
		serverAddrPort := fmt.Sprintf("%s:%s", serverIP, port)
		server.RunServer(serverAddrPort, uploadsDir)
		close(exitSignal)
	}()

	primaryIP, err := getMyPrimaryIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	serverURL := fmt.Sprintf("http://%s:%s/", primaryIP, port)
	err = internal.PrintQR(serverURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Server available at: %s\n", serverURL)
	fmt.Println("(available in all other interfaces too)")

	<-exitSignal
}
