package request

import (
	"fmt"
	"strconv"
	"time"
)

// AddSecret is structure for add struct request form variables
type AddSecret struct {
	Secret           string
	ExpireAfterViews int
	ExpireAfter      time.Duration
}

// NewAddSecret builds AddSecret on given parameters
func NewAddSecret(secret, eav, ea string) (res AddSecret, err error) {
	res.Secret = secret
	res.ExpireAfterViews, err = strconv.Atoi(eav)
	if err != nil {
		return res, err
	}
	expireAfter, err := strconv.Atoi(ea)
	if err != nil {
		return res, err
	}
	res.ExpireAfter = time.Duration(expireAfter) * time.Minute
	return res, nil
}

// Validate checks for correct values of the AddSecret struct
func (as AddSecret) Validate() error {
	if as.Secret == "" {
		return fmt.Errorf("empty secret")
	}
	if as.ExpireAfterViews < 1 {
		return fmt.Errorf("ExpireAfterViews not valid")
	}
	return nil
}
