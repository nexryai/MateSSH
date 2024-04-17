package setup

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
)

func ServeSetupWizard(initPasswordHash string) error {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Please enter the init password for setup:\n")

		term := terminal.NewTerminal(s, "> ")
		line := ""

		// Authentication loop
		passwordTries := 0
		for {
			line, _ = term.ReadPassword("Enter password: ")
			if line == "quit" {
				break
			}

			io.WriteString(s, "Wrong password\n")

			if passwordTries >= 3 {
				io.WriteString(s, "Too many tries\n")
				return
			} else {
				passwordTries++
			}
		}

		io.WriteString(s, "=======================================================\n")
		io.WriteString(s, "Welcome to MateSSH!\n")
		io.WriteString(s, "Please enter the SSH public key used for authentication.\n")
		io.WriteString(s, "WARNING: If you make a mistake here, you will have to physically access the server, log in, and run the reset command. Please enter carefully!\n")
		for {
			line, _ = term.ReadLine()
			if line == "quit" {
				break
			}
			io.WriteString(s, fmt.Sprintf("You wrote: %s\n", line))
		}
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
	return nil
}
