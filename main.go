package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/nexryai/MateSSH/internal/setup"
	"io"
	"log"
)

func main() {
	configIsExist := false

	if !configIsExist {
		err := setup.ServeSetupWizard("password")
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
