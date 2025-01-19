package gostore

import (
	"golang.org/x/crypto/ssh"
)

type Key struct {
	ID      string
	StoreID string
	ssh.PublicKey
}

type API interface {
	GetKeys(storeID string, secretIDs []string) ([]Key, error)
	GetPrivateKey(storeID, secretID string) (ssh.Signer, error)
}
