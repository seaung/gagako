package routeros

import (
	"fmt"

	"github.com/go-routeros/routeros"
)

func RouterosAuthorized(host, port, username, password string) (bool, error) {
	_, err := routeros.Dial(fmt.Sprintf("%s:%s", host, port), username, password)
	if err != nil {
		return false, err
	}

	return true, nil
}
