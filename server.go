package main
import (
    "fmt"
    "net"
    "bufio"
    "os"
    "io/ioutil"
    "strings"
    "strconv"
    "path/filepath"
)

const BUFFERSIZE = 1024

func handleConnection(conn net.Conn) {
    defer conn.Close()
    curDir, err := filepath.Abs(os.Getwd())
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
            conn.Write([]byte(curDir))
            conn.Write([]byte(fmt.Sprintf("%d\n", len(files))))
            for _, file := range files {
                conn.Write([]byte(fmt.Sprintf("%d", file.Size()) + " " + file.Name() + "\n"))
                fmt.Println(file.Name())
            }
        case 'D':
            //
            filename := filepath.Join(curDir, filepath.Base(strings.TrimSpace(message[1:])))
            file, err := os.Open(filename)
            if err != nil {
                fmt.Println("Error " + err.Error())
                conn.Write([]byte("Error " + err.Error() + "\n"))
            } else {
                defer file.Close()
                fileinfo, err := os.Stat(filename)
                if err != nil {
                    fmt.Println("Error " + err.Error())
                    conn.Write([]byte("Error " + err.Error() + "\n"))
                } else {
                    conn.Write([]byte(fmt.Sprintf("%d\n", fileinfo.Size())))
                    sendBuffer := make([]byte, BUFFERSIZE)
                    for {
                        _, err = file.Read(sendBuffer)
                        if err == io.EOF {
                            break
                        }
                        conn.Write(sendBuffer)
                    }
                }
            }
        case 'U':
            //
            filename := filepath.Join(curDir, filepath.Base(strings.TrimSpace(message[1:])))
            file, err := os.Create(filename)
            if err != nil {
                fmt.Println("Error " + err.Error())
                conn.Write([]byte("Error " + err.Error() + "\n"))
            } else {
                defer file.Close()
                filelenStr := strings.TrimSpace(bufio.NewReader(conn).ReadString('\n'))
                filelen, err := strconv.ParseInt(filelenStr, 10, 64)
                if err != nil {
                    fmt.Println("Error " + err.Error())
                    conn.Write([]byte("Error " + err.Error() + "\n"))
                } else {
                    var receivedBytes int64
                    for {
                        if (fileSize - receivedBytes) < BUFFERSIZE {
                            io.CopyN(newFile, connection, (fileSize - receivedBytes))
                            connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
                            break
                        }
                        io.CopyN(newFile, connection, BUFFERSIZE)
                        receivedBytes += BUFFERSIZE
                    }
                }
            }
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
    defer ln.Close()
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
