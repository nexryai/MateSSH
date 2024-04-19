package setup

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/nexryai/MateSSH/internal/config"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"time"
)

func isValidPublicKey(key string) bool {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	return err == nil
}

func ServeSetupWizard(initPassphrase string, hostKeys hostkey.Keyring) error {
	setupWizard := func(s ssh.Session) {
		var errs []error

		_, e := io.WriteString(s, "<Please enter the init passphrase for setup>\n")
		if e != nil {
			errs = append(errs, e)
		}

		term := terminal.NewTerminal(s, "> ")
		input := ""

		// Authentication loop
		passwordTries := 0
		for {
			input, _ = term.ReadPassword(" 🔑 Enter passphrase: ")
			if input == initPassphrase {
				break
			}

			// Wrong passphrase
			time.Sleep(5 * time.Second)
			_, e = io.WriteString(s, "Wrong passphrase\n")
			if e != nil {
				errs = append(errs, e)
			}

			if passwordTries >= 3 {
				_, e = io.WriteString(s, "Too many tries\n")
				if e != nil {
					errs = append(errs, e)
				}
				return
			} else {
				passwordTries++
			}
		}

		if errors.Join(errs...) != nil {
			log.Fatal(errors.Join(errs...))
			return
		}

		// Setup loop
		welcomeMsg := fmt.Sprintf("\n=========================================\n" +
			"Welcome to MateSSH!\n" +
			"Please enter the SSH public key used for authentication.\n" +
			"WARNING: If you make a mistake here, you will have to physically access the server, log in, and run the reset command. Please enter carefully!\n")

		_, e = io.WriteString(s, welcomeMsg)
		if e != nil {
			fmt.Println(e)
		}

		for {
			input, _ = term.ReadLine()
			if input == "quit" {
				break
			} else if isValidPublicKey(input) {
				//ToDo
				err := config.CreateConfig(hostKeys, input)
				if err != nil {
					fmt.Println(err)
					_, e = io.WriteString(s, "Sorry. Failed to save config.\n")
					if e != nil {
						fmt.Println(e)
						break
					}
				} else {
					_, e = io.WriteString(s, "Public key saved\n")
					if e != nil {
						fmt.Println(e)
					}

					// Exit the setup
					break
				}
			} else {
				_, e = io.WriteString(s, "Invalid public key\n")
				if e != nil {
					fmt.Println(e)
					break
				}
			}
		}
	}

	// Configure server
	server := ssh.Server{
		Addr:    ":2222",
		Handler: setupWizard,
	}

	// Add host key
	for _, s := range *hostKeys.Signers {
		server.AddHostKey(s)
	}

	log.Fatal(server.ListenAndServe())
	return nil
}
