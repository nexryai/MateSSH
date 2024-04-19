package server

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/nexryai/MateSSH/internal/config"
	"io"
	"log"
	"os/exec"
)

func cmdPTYLinux(s ssh.Session) {
	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		io.WriteString(s, "Something went wrong.\n")
		return
	}

	// Start the shell
	shell := exec.Command("nu")

	ptmx, err := pty.Start(shell)
	if err != nil {
		log.Printf("Failed to create pty: %v\n", err)
		io.WriteString(s, "Failed to start shell.\n")
		return
	}
	defer func() {
		_ = ptmx.Close()
		_ = shell.Wait()
	}()

	// Set the size of the pty
	pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(ptyReq.Window.Height),
		Cols: uint16(ptyReq.Window.Width),
	})
	log.Printf("Set TERM to %s\n", ptyReq.Term)
	log.Printf("Starting shell: %s\n", shell.Path)

	done := s.Context().Done()
	go func() {
		<-done
		log.Println("Session done:", s.RemoteAddr())
		ptmx.Close()
	}()

	// loop to handle window size changes
	go func() {
		for win := range winCh {
			pty.Setsize(ptmx, &pty.Winsize{
				Rows: uint16(win.Height),
				Cols: uint16(win.Width),
			})
		}
	}()

	// Copy stdin to the pty and the pty to stdout
	go func() {
		io.Copy(ptmx, s) // stdin
	}()
	io.Copy(s, ptmx) // stdout
}

func Start() error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	handler := func(s ssh.Session) {
		_, _ = s.Write([]byte("Hello world\n"))
		fmt.Print("New session from ", s.RemoteAddr(), "\n")
		cmdPTYLinux(s)
	}

	keyHandler := func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, authorizedKey := range conf.AuthorizedKeys {
			allowed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKey))
			if err != nil {
				continue
			}

			if ssh.KeysEqual(key, allowed) {
				return true
			}
		}
		return false
	}

	// Configure server
	server := &ssh.Server{
		Addr:             ":2222",
		Handler:          handler,
		PublicKeyHandler: keyHandler,
	}

	// Set host key
	hostKeyring := conf.HostKeys
	err = hostKeyring.Parse()
	if err != nil {
		return err
	}

	err = hostKeyring.GenSigners()
	if err != nil {
		return err
	}

	for _, s := range *hostKeyring.Signers {
		server.AddHostKey(s)
	}

	return server.ListenAndServe()
}
