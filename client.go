package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func listenForMessages(conn net.Conn){
    defer wg.Done()
    for {
        data := make([]byte, 1024)
        _, e := conn.Read(data)
        if e != nil {
            fmt.Println("Disconnected")
            break
        }
        fmt.Print(string(data))
    }
}

func sendMessage(conn net.Conn, msg string){
    defer wg.Done()
    conn.Write([]byte(msg))
}

func main(){

    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }

    unamePrompt := make([]byte, 1024)
    conn.Read(unamePrompt)
    fmt.Print(string(unamePrompt))
    
    var uname string
    fmt.Scanln(&uname)
    conn.Write([]byte(uname))

    wg.Add(1)
    go listenForMessages(conn)
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print(">>> ")
        scanner.Scan()
        msg := scanner.Text()
        if msg == "/exit" {
            break
        }
        sendMessage(conn, msg)
    }
    wg.Wait()


}
