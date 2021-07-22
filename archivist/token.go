package archivist

import (
	"encoding/base64"
	"fmt"
	"os"
)

func CheckTokenExists() error {
	f, err := os.Open(".jitsuin")
	if err != nil {
		return err
	}
	fmt.Printf("%v", f)
	return err
}

func ValidateToken(token string) error {
	_, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	return nil
}
