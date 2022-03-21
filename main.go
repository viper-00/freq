package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// BASH PS1 prompt
// Makes reading code easier than escape sequences, from here -> https://gist.github.com/vratiu/9780109
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

func main() {
	fmt.Println("freq by zhong-my")
	fmt.Println("--- ---")
	fmt.Print("Please input your domain: ")

	jobs := make(chan string)
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for domain := range jobs {
			resp, err := http.Get(domain)
			if err != nil {
				fmt.Println(err)
				fmt.Print("Please input your domain: ")
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				fmt.Print("Please input your domain: ")
				continue
			}
			s := string(body)
			if strings.Contains(s, "alert(1)") || strings.Contains(s, "alert('XSS'") {
				fmt.Println(string(colorRed), "Vulnerable To XSS:", domain, string(colorReset))
			} else {
				fmt.Println(string(colorGreen), "Not Vulnerable To XSS:", domain, string(colorReset))
			}
			fmt.Print("Please input your domain: ")
		}
	}()

	for scanner.Scan() {
		domain := scanner.Text()
		jobs <- domain
	}

	close(jobs)
	wg.Wait()
}
