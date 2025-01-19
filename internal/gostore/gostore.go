package gostore

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewGostore(p string) API {
	return &gostore{p: p}
}

type gostore struct {
	p string
}

func (g gostore) GetKeys(storeID string, secretIDs []string) ([]Key, error) {
	keys := make([]Key, 0, len(secretIDs))
	for _, id := range secretIDs {
		signer, err := g.getSigner(storeID, id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get signer for %s", id)
		}

		keys = append(keys, Key{
			ID:        id,
			StoreID:   storeID,
			PublicKey: signer.PublicKey(),
		})
	}

	return keys, nil
}

func (g gostore) GetPrivateKey(storeID, secretID string) (ssh.Signer, error) {
	return g.getSigner(storeID, secretID)
}

func (g gostore) getSigner(storeID, id string) (ssh.Signer, error) {
	data, err := g.exec(storeID, []string{"cat", id})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get private key")
	}

	publicKey, err := ssh.ParsePrivateKey(data)
	return publicKey, errors.Wrap(err, "failed to parse private key")
}

func (g gostore) exec(storeID string, args []string) ([]byte, error) {
	cmd := exec.Command(g.p, args...) //nolint:gosec
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOSTORE_STORE_ID=%s", storeID), // use specific store
	)

	return cmd.Output()
}
