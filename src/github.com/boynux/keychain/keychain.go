package main

import (
    "os"
    "fmt"
    "bufio"
    "log"
    "strings"
)

const Filepath string = "/var/lib/keychain/keychain.data"

type KeyValuePair struct {
    Key string
    Value string
}

type KeyStore struct {
    Filepath string
    keymap map[string]string
    sep string
}

func NewKeyStore() *KeyStore {
    return &KeyStore {
        keymap:  make(map[string]string),
        sep: "::",
    }
}

func (k *KeyStore) splitItem(s string) KeyValuePair {
    result := strings.SplitN(s, k.sep, 2)

    return KeyValuePair {
        Key: result[0],
        Value: result[1],
    }
}

func (k *KeyStore) Get(key string) string {
    return k.keymap[key]
}

func (k *KeyStore) set(key, value string) {
    k.keymap[key] = value
}

func (k *KeyStore) append(pair KeyValuePair) {
    f, _ := os.OpenFile(Filepath, os.O_APPEND|os.O_WRONLY, 0600)
    defer f.Close()

    f.WriteString(fmt.Sprintf("%s%s%s\n", pair.Key, k.sep, pair.Value))
}

func (k *KeyStore) Serve(c chan KeyValuePair) {
    for pair := range c {
        k.append(pair)
        k.set(pair.Key, pair.Value)
    }
}

func (k *KeyStore) Load() {
    log.Println("Loading file data")
    reader, err := os.Open(Filepath)
    defer reader.Close()

    if err != nil {
        log.Fatal("could not load file data: ", err)
    }

    buffer := bufio.NewScanner(reader)
    buffer.Split(bufio.ScanLines)

    for buffer.Scan() {
        pair := k.splitItem(buffer.Text())
        k.keymap[pair.Key] = pair.Value
    }
}
