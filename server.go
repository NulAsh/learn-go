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

func listDir(curDir string, conn net.Conn) (err error) {
    files, err := ioutil.ReadDir(curDir)
    if err != nil {
        return
    }
    conn.Write([]byte(curDir))
    conn.Write([]byte(fmt.Sprintf("\n%d\n", len(files))))
    for _, file := range files {
        if file.IsDir() {
            conn.Write([]byte("D " + file.Name() + "\n"))
        } else {
            conn.Write([]byte(fmt.Sprintf("%d", file.Size()) + " " + file.Name() + "\n"))
        }
        fmt.Println(file.Name())
    }
    return
}

func downloadFile(filename string, conn net.Conn) (err error){
    file, err := os.Open(filename)
    if err != nil {
        return
    }
    defer file.Close()
    fileinfo, err := os.Stat(filename)
    if err != nil {
        return
    }
    conn.Write([]byte(fmt.Sprintf("%d\n", fileinfo.Size())))
    wbytes, err := io.Copy(conn, file)
    if err != nil {
        return
    }
    if wbytes != fileinfo.Size() {
        fmt.Printf("%d bytes of %d written\n", wbytes, fileinfo.Size())
    }
    return
}

func uploadFile(filename string, conn net.Conn) (err error){
    file, err := os.Create(filename)
    if err != nil {
        return
    }
    defer file.Close()
    filelenStr, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        return
    }
    filelen, err := strconv.ParseInt(strings.TrimSpace(filelenStr), 10, 64)
    if err != nil {
        return
    }
    rbytes, err := io.CopyN(file, conn, filelen)
    if err != nil {
        return
    }
    if rbytes != filelen {
        fmt.Printf("%d bytes of %d read\n", rbytes, filelen)
    }
    return
}

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
            err := listDir(curDir, conn)
            if err != nil {
                fmt.Println("Error: " + err.Error())
                return
            }
        case 'D':
            //
            filename := filepath.Join(curDir, filepath.Base(strings.TrimSpace(message[1:])))
            err := downloadFile(filename, conn)
            if err != nil {
                fmt.Println("Error " + err.Error())
                conn.Write([]byte("Error " + err.Error() + "\n"))
                return
            }
        case 'U':
            //
            filename := filepath.Join(curDir, filepath.Base(strings.TrimSpace(message[1:])))
            err := uploadFile(filename, conn)
            if err != nil {
                fmt.Println("Error " + err.Error())
                conn.Write([]byte("Error " + err.Error() + "\n"))
                return
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
