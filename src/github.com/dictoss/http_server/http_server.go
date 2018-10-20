package main

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "encoding/json"
)

func main() {
    const c_read_timeout_sec = 10 * time.Second
    const c_write_timeout_sec = 10 * time.Second
    const c_max_header_bytes = 1 << 20 // 1MB

    m := http.NewServeMux()

    m.Handle("/hello", http.HandlerFunc(handler_hello))
    m.Handle("/rest/hello", http.HandlerFunc(handler_rest_hello))

    s := http.Server{
        Addr: ":8090",
        Handler: m,
        ReadTimeout:    c_read_timeout_sec,
        WriteTimeout:   c_write_timeout_sec,
        MaxHeaderBytes: c_max_header_bytes,
    }

    log.Fatal(s.ListenAndServe())
}

func handler_hello(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Hello World from Go.")
}

// require StructTag UpperCase in Member.
type RequestRestMessage struct {
    ReqMsg string `json:"req_msg"`
}

type ResponseRestMessage struct {
    ReqMsg string `json:"req_msg"`
    ResMsg string `json:"res_msg"`
}

func handler_rest_hello(w http.ResponseWriter, r *http.Request){
    fmt.Print("IN handler_rest_hello()\n")

    req_json_str := `{"req_msg": "Hello."}`
    req_json_bytes := ([]byte)(req_json_str)
    req_rrm := new(RequestRestMessage)
    
    if err := json.Unmarshal(req_json_bytes, &req_rrm); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    fmt.Print("in param:\n")
    fmt.Printf("  ReqMsg: %s\n", req_rrm.ReqMsg)    

    res_rrm := ResponseRestMessage{
       ReqMsg: req_rrm.ReqMsg,
       ResMsg: "Me too !!"}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)    
    json.NewEncoder(w).Encode(res_rrm)
}
