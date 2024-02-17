package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func sliceFind(s []net.Conn, target net.Conn) int { 
    for i, v := range s{
        if v.RemoteAddr().String() == target.RemoteAddr().String() {
            return i
        }
    }
    return -1
}

const hostport = ":4444"
var wg sync.WaitGroup
var users []net.Conn

func broadcast(msg string, conn net.Conn){
    for _, user := range users {
        if user.RemoteAddr().String() != conn.RemoteAddr().String() {
            user.Write([]byte(msg))
        }
    }
}

func handleConnection(conn net.Conn){
    wg.Add(1)
    defer wg.Done()

    if sliceFind(users, conn) == -1{
        users = append(users, conn)
    }
    broadcast(conn.RemoteAddr().String() + " connected", conn)

    for {
        msg := make([]byte, 1024)
        _, e := conn.Read(msg)
        if e != nil {
            broadcast(conn.LocalAddr().String() + " disconnected", conn)
            idx := sliceFind(users, conn)
            var temp []net.Conn
            for i, v := range users {
                if i != idx {
                    temp = append(temp, v)
                }
            }
            users = temp
            break
        }
        broadcast(conn.RemoteAddr().String() + " : " + string(msg), conn)
    }
}

func main(){

    defer wg.Wait()

    socket, err := net.Listen("tcp", hostport)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Couldn't listen on port " + hostport)
        os.Exit(-1)
    }
    defer socket.Close()
    fmt.Println("Started listening on " + hostport)

    for {
        conn, err := socket.Accept()
        if err != nil {
            fmt.Println("Error accepting connection : " + err.Error())
            return
        }

        go handleConnection(conn)
    }
}
