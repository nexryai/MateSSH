package main

import (
	"fmt"
	"github.com/nexryai/MateSSH/internal/config"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"github.com/nexryai/MateSSH/internal/server"
	"github.com/nexryai/MateSSH/internal/setup"
	"github.com/sethvargo/go-diceware/diceware"
	"log"
	"strings"
)

func main() {
	if !config.IsExist() {
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
		log.Fatal(server.Start())
	}
}
