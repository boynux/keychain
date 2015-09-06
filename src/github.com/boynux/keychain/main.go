package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
)

var (
    keystore *KeyStore = new(KeyStore)
    store chan KeyValuePair
)

func GetHandler(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    key := vars["key"]

    value := keystore.Get(key)

    fmt.Fprintln(w, value)
}

func SetHandler(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)

    kv := KeyValuePair {
        Key: vars["key"],
        Value: vars["value"],
    }

    store <- kv
}


func main() {
    keystore = NewKeyStore()
    keystore.Load()

    store = make(chan KeyValuePair)
    go keystore.Serve(store)

    r := mux.NewRouter()
    r.HandleFunc("/mapstore/get/{key}", GetHandler).
      Methods("GET")

    r.HandleFunc("/mapstore/set/{key}/{value}", SetHandler).
      Methods("PUT")

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
