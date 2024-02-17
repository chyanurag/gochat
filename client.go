package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func sendMessage(msg string, conn net.Conn){
    wg.Add(1)
    defer wg.Done()
    conn.Write([]byte(msg))
}

func getMessage(conn net.Conn){
    wg.Add(1)
    defer wg.Done()
    for {
        buffer := make([]byte, 1024)
        conn.Read(buffer)
        fmt.Println(string(buffer))
    }
}

func main(){

    conn, err := net.Dial("tcp", ":4444")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed connecting to server")
        return
    }

    go getMessage(conn)
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan(){
        msg := scanner.Text()
        sendMessage(msg, conn)
    }

    wg.Wait()

}
