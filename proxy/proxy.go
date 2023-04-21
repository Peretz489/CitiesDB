package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const NODECOUNT = 2 //set number of nodes here

func proxyHandler(currentNode *int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nodePort := ":" + strconv.Itoa(9000+*currentNode)
		*currentNode += 1
		if *currentNode > NODECOUNT {
			*currentNode = 1
		}
		byteBody, _ := io.ReadAll(r.Body)
		buff := bytes.NewBuffer(byteBody)
		defer r.Body.Close()
		q := ""
		if r.URL.RawQuery != "" {
			q = "?" + r.URL.RawQuery
		}
		requestPath := "http://127.0.0.1" + nodePort + r.URL.Path + q
		fmt.Println("Redirected to",requestPath)
		req, err := http.NewRequest(r.Method, requestPath, buff)
		if err != nil {
			log.Fatal(err.Error())
		}
		client := &http.Client{}
		responce, err := client.Do(req)
		if err != nil {
			log.Fatal(err.Error())
		}

		w.WriteHeader(responce.StatusCode)
		w.Header().Set("Content-Type", responce.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", responce.Header.Get("Content-Length"))
		io.Copy(w, responce.Body)
		defer responce.Body.Close()

	}
}

func main() {
	var currentNode int = 1
	for i := 1; i <= NODECOUNT; i++ {
		nodeArg := "-p=" + strconv.Itoa(9000+i)
		cmd := exec.Command("./node.exe", nodeArg)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", proxyHandler(&currentNode))

	srv := &http.Server{
		Addr:    "127.0.0.1:9000",
		Handler: nil,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Proxy shutdown Failed:%+v", err)
	}
	log.Print("Proxy stopped")
}
