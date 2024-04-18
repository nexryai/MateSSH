package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"github.com/nexryai/MateSSH/internal/setup"
	"github.com/sethvargo/go-diceware/diceware"
	"io"
	"log"
	"strings"
)

func main() {
	configIsExist := false

	if !configIsExist {
		// Generate host key
		hostKeyring := hostkey.Keyring{}
		err := hostKeyring.Generate()
		if err != nil {
			log.Fatal(err)
		}

		// Generate a passphrase
		passPhrasesList, err := diceware.Generate(8)
		if err != nil {
			log.Fatal(err)
		}

		initPassphrase := strings.Join(passPhrasesList, "-")
		fmt.Println("Your init passphrase is: ", initPassphrase)

		err = setup.ServeSetupWizard(initPassphrase, hostKeyring)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ssh.Handle(func(s ssh.Session) {
			io.WriteString(s, "Hello from MateSSH\n")
		})

		log.Fatal(ssh.ListenAndServe(":2222", nil))
	}
}
