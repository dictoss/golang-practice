// go version <= 1.11 or higher
package main

import (
    "log"
    "fmt"
    "net"
    "net/http"
    "net/http/fcgi"
    "encoding/json"
)

//
// struct
//

// require StructTag UpperCase in Member.
type ResponseHelloMessage struct {
    Message string `json:"message"`
}

type RequestHello2Message struct {
    ReqMsg string `json:"req_msg"`
}

type ResponseHello2Message struct {
    ReqMsg string `json:"req_msg"`
    ResMsg string `json:"res_msg"`
}
        

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "root fcgi with golang !!")
}

func handler_hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello fcgi with golang !!")
}

func handler_json_hello(w http.ResponseWriter, r *http.Request) {
    res_rrm := ResponseHelloMessage{
        Message: "Hello fcgi with golang !!"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res_rrm)
}

// see https://golang.org/pkg/net/http/
func handler_json_hello2(w http.ResponseWriter, r *http.Request) {
    //
    // request json format {"req_msg": ""}
    //
    res_rrm := ResponseHello2Message{
        ReqMsg: "",
        ResMsg: ""}

    log.Println("----------")
    log.Println(r.Method, ",", r.UserAgent())

    if r.Method == "POST" {
        req_rrm := new(RequestHello2Message)

        if err := json.NewDecoder(r.Body).Decode(&req_rrm); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        res_rrm.ReqMsg = req_rrm.ReqMsg
        res_rrm.ResMsg = "Me too !!"
    } else if r.Method == "GET" {
        // result default message
        res_rrm.ResMsg = "Hello !!"
    } else {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    log.Println(res_rrm)
        
    res_json_str, err := json.Marshal(&res_rrm)
    if err != nil {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")

        res_len := len(res_json_str)
        w.Header().Set("Content-Length", fmt.Sprint(res_len))

        //log.Println(w.Header())

        w.WriteHeader(http.StatusOK)
        w.Write(res_json_str)
        
        log.Println("==========")
    }
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    l, err := net.Listen("tcp", "127.0.0.1:9000")

    if err != nil {
        return
    }

    m := http.NewServeMux()
 
    fcgi_proxy_prefix := "/gofcgi"
    
    m.Handle(fcgi_proxy_prefix + "/hello/", http.HandlerFunc(handler_hello))
    m.Handle(fcgi_proxy_prefix + "/json/hello/", http.HandlerFunc(handler_json_hello))
    m.Handle(fcgi_proxy_prefix + "/json/hello2/", http.HandlerFunc(handler_json_hello2))

    // fallback routing
    m.Handle("/", http.HandlerFunc(handler))
    
    fcgi.Serve(l, m)
}
