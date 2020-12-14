package main

import (
	"bufio"
	"os"

	"./crypto"
	"./hmi"
	"./storage"
)

func main() {
	storage := storage.NewFileSystemPasswordStorage("~/.resources")
	cryptoProvider := crypto.NewPersistentCryptoProvider(storage)
	scanner := bufio.NewScanner(os.Stdin)
	hmi := hmi.NewHMI(cryptoProvider, scanner)
	hmi.Start()
}
