/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package net

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	host   string
	client = http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{Timeout: 2 * time.Second}).Dial,
		},
	}
)

func ping(domain string) (int, error) {
	url := "http://" + domain
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "This pings a remotre host",
	Long:  `This pings a remote host and returns the results to the user in the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := ping(host); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Response:", resp)
		}
	},
}

func init() {

	pingCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to ping")

	if err := pingCmd.MarkFlagRequired("host"); err != nil {
		panic(err)
	}
	NetCmd.AddCommand(pingCmd)
}
