package smb

import (
	"log"
	"testing"

	"github.com/stacktitan/smb/smb"
)

func TestSMB(t *testing.T) {

	host := "192.168.100.21"
	options := smb.Options{
		Host:        host,
		Port:        445,
		User:        "GUEST",
		Domain:      "",
		Workstation: "",
		Password:    "",
	}
	debug := false
	session, err := smb.NewSession(options, debug)
	if err != nil {
		log.Fatalln("[!]", err)
	}
	defer session.Close()

	if session.IsSigningRequired {
		log.Println("[-] Signing is required")
	} else {
		log.Println("[+] Signing is NOT required")
	}

	if session.IsAuthenticated {
		log.Println("[+] Login successful")
	} else {
		log.Println("[-] Login failed")
	}

	if err := session.TreeConnect("99_Corporate-IT"); err != nil {
		log.Println(err)
		return
	}
	defer session.TreeDisconnect("99_Corporate-IT")
	session.NewNegotiateReq()

	log.Println("OK!")
}
