package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func main() {
	c, err := net.Dial("udp", "127.0.0.1:5040")
	if err != nil {
		log.Fatal(err)
	}
	c.SetDeadline(time.Now().Add(time.Second))
	var buf = make([]byte, 8)
	c.Write(buf)
	_, err = c.Read(buf)
	if err != nil {
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			log.Println("net timeout")
		}
	}

	req, err := http.NewRequest(http.MethodGet, "127.0.0.1:5040", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Timeout: time.Second}
	_, err = client.Do(req)
	if err != nil {
		urlErr, ok := err.(*url.Error)
		if ok && urlErr.Timeout() {
			log.Println("http timeout")
		}
	}
	fmt.Println(err)
}
