// go version <= 1.11 or higher
package main

import (
    "flag"
    "errors"
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

type GlobalConfig struct {
    logpath string
    fcgi_listen_addr string
    fcgi_url_prefix string
}

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

// global var
var g_conf GlobalConfig


// function
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

func init_prog() error {
     flag.StringVar(&g_conf.logpath, "logpath", "", "log output path. If you set to empty string, output log to stdout.")
     flag.StringVar(&g_conf.fcgi_listen_addr, "fcgi_listen_addr", "127.0.0.1:9000" , "fast cgi listen address and port. Default 127.0.0.0:9000 .")
     flag.StringVar(&g_conf.fcgi_url_prefix, "fcgi_url_prefix", "/gofcgi", "fast cgi prefix url path. Default /gofcgi .")

    // set log
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // valdate execute parameter
    if true {
        log.Println("[g_conf.logpath]", g_conf.logpath)
        log.Println("[g_conf.fcgi_listen_addr]", g_conf.fcgi_listen_addr)
        log.Println("[g_conf.fcgi_url_prefix]", g_conf.fcgi_url_prefix)

        return nil
    } else {
        return errors.New("invalid paramer.")
    }
}

func main() {
    err := init_prog()
    if err != nil {
        log.Fatalln(err)
        return
    } else {
        log.Println("load config.")
    }

    l, err := net.Listen("tcp", g_conf.fcgi_listen_addr)
    if err != nil {
        return
    }

    m := http.NewServeMux()
 
    fcgi_proxy_prefix := g_conf.fcgi_url_prefix
    
    m.Handle(fcgi_proxy_prefix + "/hello/", http.HandlerFunc(handler_hello))
    m.Handle(fcgi_proxy_prefix + "/json/hello/", http.HandlerFunc(handler_json_hello))
    m.Handle(fcgi_proxy_prefix + "/json/hello2/", http.HandlerFunc(handler_json_hello2))

    // fallback routing
    m.Handle("/", http.HandlerFunc(handler))
    
    fcgi.Serve(l, m)
}
