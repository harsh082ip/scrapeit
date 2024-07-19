package helpers

import (
	"fmt"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.NewVerifier().DisableCatchAllCheck()
)

func VerifyEmail(email string) (bool, error) {

	ret, err := verifier.Verify(email)
	if err != nil {
		return false, fmt.Errorf("verification failed: %v", err)
	}

	// for invalid syntax
	if !ret.Syntax.Valid {
		return false, fmt.Errorf("verification failed: %v", "Invalid Syntax")
	}

	return true, nil
}
