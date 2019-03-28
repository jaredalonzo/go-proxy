// Jared Alonzo
// A simple proxy
// Version 1.0

package main

import (
  "fmt"
  "log"
  "net/http"
  "io"
  "os"
)

type HTTPConnection struct {
  Request *http.Request
  Response *http.Response
}

type Proxy struct {
}

type HTTPConnectionChannel chan *HTTPConnection

var connectionChannel = make(HTTPConnectionChannel)

func PrintUsage() {
  fmt.Println("Usage: ./proxy [port#]")
}

func PrintHTTP(conn *HTTPConnection) {
  fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
  for k, v := range conn.Request.Header {
    fmt.Println(k, ":", v)
  }
  fmt.Println("\n====================================\n")
  fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
  for k, v := range conn.Response.Header {
    fmt.Println(k, ":", v)
  }
  fmt.Println(conn.Response.Body)
  fmt.Println("\n====================================\n")
}

func HandleHTTP() {
  for {
    select {
    case conn := <-connectionChannel:
      PrintHTTP(conn)
    }
  }
}

func NewProxy() *Proxy {
  return &Proxy{}
}

func (p *Proxy) ServeHTTP(write http.ResponseWriter, read *http.Request) {
  var res *http.Response
  var err error
  var req *http.Request
  client := &http.Client{}

  // ** redirecting fun **
  //res, err = client.Get("https://www.google.com")
  //if err != nil {
  //  http.Error(write, err.Error(), http.StatusInternalServerError)
  //  return
  //}
  //write.WriteHeader(res.StatusCode)
  //io.Copy(write, res.Body)
  //res.Body.Close()

  req, err = http.NewRequest(read.Method, read.RequestURI, read.Body)
  // grab header: val of read and append to req instance
  for hdr, val := range read.Header {
    req.Header.Set(hdr, val[0])
  }
  res, err = client.Do(req)
  read.Body.Close() // close client to proxy

  // combined for GET/POST
  if err != nil {
    http.Error(write, err.Error(), http.StatusInternalServerError)
    return
  }

  // set up response to client
  conn := &HTTPConnection{read, res}

  // copy response from server headers to response to client headers
  for hdr, val := range res.Header {
    write.Header().Set(hdr, val[0])
  }
  write.WriteHeader(res.StatusCode)
  // copy response from server body to reposne to client body
  io.Copy(write, res.Body)
  res.Body.Close() // close server to proxy

  PrintHTTP(conn)
}

func main() {
  //go HandleHTTP()
  var port string 
  if len(os.Args) < 2 {
    PrintUsage()
    return
  } else {
    port = os.Args[1]
  }
  proxy := NewProxy()
  fmt.Println("\n====================================\n")
  err := http.ListenAndServe(":"+port, proxy)
  if err != nil {
    log.Fatal("ListenAndServe: ", err.Error())
  }
}
