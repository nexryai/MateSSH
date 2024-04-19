package main

import (
	"github.com/nexryai/MateSSH/internal/config"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"github.com/nexryai/MateSSH/internal/logger"
	"github.com/nexryai/MateSSH/internal/server"
	"github.com/nexryai/MateSSH/internal/setup"
	"github.com/sethvargo/go-diceware/diceware"
	"strings"
)

func main() {
	log := logger.GetLogger("main")

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
