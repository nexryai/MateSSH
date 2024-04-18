package setup

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"time"
)

func isValidPublicKey(key string) bool {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	return err == nil
}

func ServeSetupWizard(initPassphrase string) error {
	ssh.Handle(func(s ssh.Session) {
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
				break
			} else {
				_, e = io.WriteString(s, "Invalid public key\n")
				if e != nil {
					fmt.Println(e)
					break
				}
			}

			_, e = io.WriteString(s, fmt.Sprintf("You wrote: %s\n", input))
			if e != nil {
				fmt.Println(e)
				break
			}
		}
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
	return nil
}
