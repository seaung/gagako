package smb

import "github.com/stacktitan/smb/smb"

func smbAuthorized(host, port, username, password string) (ok bool, err error) {
	options := smb.Options{
		Host:        host,
		Port:        445,
		User:        username,
		Password:    password,
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err != nil {
		return false, err
	}

	if session.IsAuthenticated {
		ok = true
	}

	defer session.Close()
	return ok, err
}
