package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

func Sha256Example(data string) {
	h := sha256.New()
	h.Write([]byte(data))
	fmt.Printf("%x\n", h.Sum(nil))
}
func Sha256Example2(data string) {
	//h := sha256.New()
	//h.Write([]byte("hello world\n"))
	fmt.Printf("%x\n", sha256.Sum256([]byte(data)))
}

func Sha224Example(data string) {
	//h := sha256.New()
	//h.Write([]byte("hello world\n"))
	fmt.Printf("%x\n", sha256.Sum224([]byte(data)))
}

type Block2 struct {
	Nonce        int
	PreviousHash [32]byte
	timestamp    int64
	transactions []string
}

func (b *Block2) Hash() [32]byte {
	m, _ := json.Marshal(b)
	log.Println("json mashall:", m)
	log.Println("string() json mashall:", string(m))
	return sha256.Sum256([]byte(m))
}
