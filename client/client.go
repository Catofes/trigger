package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	url := flag.String("c", "", "Connect URL.")
	command := flag.String("e", "", "Command.")
	timeout := flag.Int64("t", 3600, "Timeout")
	flag.Parse()
	if url == nil || command == nil {
		log.Fatal("Wrong input.")
	}
	for {
		ctx, cancel := context.WithCancel(context.Background())
		requestURL := fmt.Sprintf("%s/wait/%d", *url, time.Now().Unix())
		request, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		done := make(chan *http.Response)
		go func(req *http.Request, done chan *http.Response) {
			response, _ := http.DefaultClient.Do(request)
			done <- response
		}(request, done)
		select {
		case <-time.After(time.Duration(*timeout) * time.Second):
			log.Println("Timeout, restart.")
		case response := <-done:
			if response.StatusCode == http.StatusOK {
				log.Println("Update, execute command.")
				command := strings.Split(*command, " ")
				var cmd *exec.Cmd
				if len(command) == 1 {
					cmd = exec.Command(command[0])
				} else {
					cmd = exec.Command(command[0], command[1:]...)
				}
				log.Printf("Excuting command : %s.", cmd.Path)
				out, err := cmd.CombinedOutput()
				if err != nil {
					log.Fatalf("err: %s\nout: %s\n", err, out)
				} else {
					log.Printf("Excuted, result:\n %s", string(out))
				}
			} else if response.StatusCode == http.StatusNoContent {
				log.Println("Connect lost, restart.")
			}
		}
		cancel()
		time.Sleep(5 * time.Second)
	}
}
