/*
Copyright © 2023 Idan Moral idan22moral@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/idan22moral/pasta/internal"
	"github.com/idan22moral/pasta/internal/server"
	"github.com/spf13/cobra"
)

const defaultServerAddr = "0.0.0.0"
const defaultServerPort = "8080"

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

var rootCmd = &cobra.Command{
	Use:   "pasta uploads_directory",
	Short: "Copy-paste files between devices",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("uploads directory must be specified")
		}

		stat, err := os.Stat(args[0])

		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err == nil && !stat.IsDir() {
			return fmt.Errorf("the specified path is not a directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		uploadsDir := args[0]
		serverAddr, err := cmd.Flags().GetString("address")
		if err != nil {
			fmt.Println(err)
			return
		}

		serverPort, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Println(err)
			return
		}

		serverAddrPort := fmt.Sprintf("%s:%s", serverAddr, serverPort)
		exitSignal := make(chan interface{})
		go func() {
			server.RunServer(serverAddrPort, uploadsDir)
			close(exitSignal)
		}()

		displayServerAddr := serverAddr
		if serverAddr == defaultServerAddr {
			primaryIP, err := getMyPrimaryIP()
			if err != nil {
				fmt.Println(err)
				return
			}
			displayServerAddr = primaryIP
		}

		serverURL := fmt.Sprintf("http://%s:%s/", displayServerAddr, serverPort)
		err = internal.PrintQR(serverURL)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Server available at: %s\n", serverURL)
		fmt.Println("(available in all other interfaces too)")

		uploadsDirAbs, err := filepath.Abs(uploadsDir)
		if err != nil {
			uploadsDirAbs = uploadsDir
		}
		fmt.Printf("Uploaded files will be stored at: %s\n", uploadsDirAbs)

		<-exitSignal
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("address", "a", defaultServerAddr, "bind the HTTP server to this address")
	rootCmd.Flags().StringP("port", "p", defaultServerPort, "bind the HTTP server to this port")
}
