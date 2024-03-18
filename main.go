package main

import (
	"fmt"
	"net"
	"sync"
    "strings"
)

var wg sync.WaitGroup

var users []net.Conn

func broadcast(msg string){
    for _, u := range users {
        u.Write([]byte(msg))
    }
}

func handleClient(conn net.Conn){
    defer wg.Done()
    conn.Write([]byte("Please provider your username "))
    name := make([]byte, 1024)
    conn.Read(name)
    username := strings.Split(strings.TrimSpace(string(name)), "\n")[0]
    fmt.Println(username, "connected")
    for {
        data := make([]byte, 1024)
        _, err := conn.Read(data)
        if err != nil {
            fmt.Println(username, "disconnected")
            break
        }
        fmt.Println(username + " : "+string(data))
        broadcast(username + " : "+string(data))
    }
}

func main(){

    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            panic(err)
        }
        go handleClient(conn)
    }

}
