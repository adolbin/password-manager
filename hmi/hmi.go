package hmi

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"../crypto"
)

type HMICommand interface {
	Execute(cryptoStorage crypto.CryptoStorage, scanner *bufio.Scanner, input string, password []byte) bool
	GetPriority() int
}

var commands = make([]HMICommand, 0)

func register(command HMICommand) {
	commands = append(commands, command)
}

type HMI struct {
	cryptoStorage crypto.CryptoStorage
	scanner       *bufio.Scanner
}

func (hmi *HMI) Start() {
	sort.Slice(commands, func(i int, j int) bool {
		return commands[i].GetPriority() < commands[j].GetPriority()
	})

	password := hmi.requestPassword()
	hmi.cryptoStorage.Init(password)
	fmt.Println("What do you want?")
	hmi.repl(password)
}

func (hmi *HMI) requestPassword() []byte {
	fmt.Println("Who are you?")

	setSecureTerminalMode(true)

	hmi.scanner.Scan()
	password := hmi.scanner.Text()

	setSecureTerminalMode(false)
	return []byte(password)
}

func (hmi *HMI) repl(password []byte) {
	for hmi.scanner.Scan() {
		input := hmi.scanner.Text()
		commandExecuted := false
		for i := 0; i < len(commands) && !commandExecuted; i++ {
			commandExecuted = commands[i].Execute(hmi.cryptoStorage, hmi.scanner, input, password)
		}

		if !commandExecuted {
			fmt.Println("Dunno :(")
		}
		fmt.Println("What do you want?")
	}
}

func setSecureTerminalMode(secureMode bool) {
	var arg string
	if secureMode {
		arg = "-echo"
	} else {
		arg = "echo"
	}
	passwordTerminalMode := exec.Command("stty", arg)
	passwordTerminalMode.Stdin = os.Stdin
	err := passwordTerminalMode.Run()
	if err != nil {
		panic(err)
	}
}

func NewHMI(cryptoStorage crypto.CryptoStorage, scanner *bufio.Scanner) *HMI {
	return &HMI{
		cryptoStorage,
		scanner,
	}
}
