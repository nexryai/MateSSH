package main

import (
	"fmt"
	"github.com/nexryai/MateSSH/internal/config"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"github.com/nexryai/MateSSH/internal/logger"
	"github.com/nexryai/MateSSH/internal/server"
	"github.com/nexryai/MateSSH/internal/setup"
	"github.com/sethvargo/go-diceware/diceware"
	"os"
	"strings"
)

func main() {
	fmt.Printf("Starting MateSSH server...\n\n")

	log := logger.GetLogger("main")

	if os.Geteuid() == 0 {
		log.Fatal("MateSSH must not be run as root")
		os.Exit(1)
	}

	if !config.IsExist() {
		// Generate host key
		hostKeyring := hostkey.Keyring{}
		err := hostKeyring.Generate()
		if err != nil {
			log.FatalWithDetail("Failed to generate host keys", err)
		}

		// Generate a passphrase
		passPhrasesList, err := diceware.Generate(8)
		if err != nil {
			log.FatalWithDetail("Failed to generate passphrase", err)
		}

		initPassphrase := strings.Join(passPhrasesList, "-")
		log.Info("Your init passphrase is: ", initPassphrase)

		err = setup.ServeSetupWizard(initPassphrase, hostKeyring)
		if err != nil {
			log.FatalWithDetail("Failed to start setup wizard", err)
		}
	} else {
		err := server.Start()
		log.FatalWithDetail("Failed to start server", err)
	}
}
