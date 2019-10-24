package main
import (
    "fmt"
    "net"
    "bufio"
    "os"
    "io"
    "io/ioutil"
    "strings"
    "strconv"
    "path/filepath"
)

const BUFFERSIZE = 1024

func handleConnection(conn net.Conn) {
    defer conn.Close()
    curDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error: " + err.Error())
        return
    }
    curDir, err = filepath.Abs(curDir)
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
            conn.Write([]byte(fmt.Sprintf("\n%d\n", len(files))))
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
                    wbytes, err := io.Copy(conn, file)
                    if err != nil {
                        fmt.Println("Error " + err.Error())
                        conn.Write([]byte("Error " + err.Error() + "\n"))
                    } else if wbytes != fileinfo.Size() {
                        fmt.Printf("%d bytes of %d written\n", wbytes, fileinfo.Size())
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
                filelenStr, err := bufio.NewReader(conn).ReadString('\n')
                if err != nil {
                    fmt.Println("Error " + err.Error())
                    conn.Write([]byte("Error " + err.Error() + "\n"))
                } else {
                    filelen, err := strconv.ParseInt(strings.TrimSpace(filelenStr), 10, 64)
                    if err != nil {
                        fmt.Println("Error " + err.Error())
                        conn.Write([]byte("Error " + err.Error() + "\n"))
                    } else {
                        rbytes, err := io.CopyN(file, conn, filelen)
                        if err != nil {
                            fmt.Println("Error " + err.Error())
                            conn.Write([]byte("Error " + err.Error() + "\n"))
                        } else if rbytes != filelen {
                            fmt.Printf("%d bytes of %d read\n", rbytes, filelen)
                        }
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
