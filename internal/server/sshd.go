package server

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/nexryai/MateSSH/internal/config"
	"github.com/nexryai/MateSSH/internal/logger"
	"io"
	"os/exec"
)

func cmdPTYLinux(s ssh.Session) {
	log := logger.GetLogger("Session")
	log.Info(fmt.Sprintf("Session opened: %v", s.RemoteAddr()))

	ptyReq, winCh, isPty := s.Pty()
	if !isPty {
		io.WriteString(s, "Something went wrong.\n")
		return
	}

	// Start the shell
	shell := exec.Command("nu")

	ptmx, err := pty.Start(shell)
	if err != nil {
		log.FatalWithDetail("Failed to create pty", err)
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
	log.Info("Set TERM to ", ptyReq.Term)
	log.Info("Starting shell: ", shell.Path)

	done := s.Context().Done()
	go func() {
		<-done
		log.Info(fmt.Sprintf("Session closed: %s", s.RemoteAddr()))
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
	log := logger.GetLogger("Server")

	// Load configuration
	log.ProgressInfo("Loading configuration...")
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	log.ProgressOk()

	// Create server
	log.ProgressInfo("Creating server...")
	handler := func(s ssh.Session) {
		_, _ = s.Write([]byte("Hello world\n"))
		cmdPTYLinux(s)
	}

	keyHandler := func(ctx ssh.Context, key ssh.PublicKey) bool {
		log := logger.GetLogger("Auth")
		log.Debug("Authenticating user: ", ctx.User())

		for _, authorizedKey := range conf.AuthorizedKeys {
			allowed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKey))
			if err != nil {
				log.Warn("Failed to parse authorized key: %v\n", err.Error())
				continue
			}

			if ssh.KeysEqual(key, allowed) {
				log.Info("Authorized key for user: ", ctx.User())
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

	log.ProgressOk()
	fmt.Println("")

	log.Info("Server started on", server.Addr)
	return server.ListenAndServe()
}
