package main
import (
    "fmt"
    "net"
    "bufio"
    "os"
    "io/ioutil"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    curDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error: " + err.Error())
        return
    }
    for {
        message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println("Error: " + err.Error())
            return
        }
        switch message[0] {
        case 'L':
            // Listing of the current directory
            files, err := ioutil.ReadDir(curDir)
            if err != nil {
                fmt.Println("Error: " + err.Error())
                return
            }
            conn.Write([]byte(fmt.Sprintf("%d\n", len(files))))
            for _, file := range files {
                conn.Write([]byte(fmt.Sprintf("%d", file.Size()) + " " + file.Name() + "\n"))
                fmt.Println(file.Name())
            }
        case 'D':
            //
        case 'U':
            //
        }
        fmt.Printf("%#v\n", message)
    }
}
func main() {
    ln, err := net.Listen("tcp", ":1234")
    if err != nil {
        // handle error
        fmt.Println(err.Error())
        return
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            // handle error
            fmt.Println(err.Error())
            continue
        }
        go handleConnection(conn)
    }
}