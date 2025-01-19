package sshagent

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"

	"github.com/UsingCoding/gostore-agent/internal/config"
	"github.com/UsingCoding/gostore-agent/internal/gostore"
)

func NewAgent(
	stores config.Stores,
	api gostore.API,
) sshagent.ExtendedAgent {
	return &agent{
		stores: stores,
		api:    api,
	}
}

var (
	errOperationUnsupported = errors.New("operation unsupported")
)

type agent struct {
	stores config.Stores
	api    gostore.API
}

func (a agent) List() ([]*sshagent.Key, error) {
	keys, err := a.getAllKeys()
	if err != nil {
		return nil, err
	}

	agentKeys := make([]*sshagent.Key, 0, len(keys))
	for _, key := range keys {
		agentKeys = append(agentKeys, &sshagent.Key{
			Format:  key.Type(),
			Blob:    key.Marshal(),
			Comment: fmt.Sprintf("gostore %s secret %s", key.StoreID, key.ID),
		})
	}

	return agentKeys, nil
}

func (a agent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return a.SignWithFlags(key, data, 0)
}

func (a agent) Add(sshagent.AddedKey) error {
	return errOperationUnsupported
}

func (a agent) Remove(_ ssh.PublicKey) error {
	return errOperationUnsupported
}

func (a agent) RemoveAll() error {
	return nil
}

func (a agent) Lock([]byte) error {
	return errOperationUnsupported
}

func (a agent) Unlock([]byte) error {
	return errOperationUnsupported
}

func (a agent) Signers() ([]ssh.Signer, error) {
	keys, err := a.getAllKeys()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keys")
	}

	signers := make([]ssh.Signer, 0, len(keys))
	for _, key := range keys {
		signer, err2 := a.api.GetPrivateKey(key.StoreID, key.ID)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to get private key for %s", key.ID)
		}
		signers = append(signers, signer)
	}

	return signers, nil
}

func (a agent) SignWithFlags(key ssh.PublicKey, data []byte, flags sshagent.SignatureFlags) (*ssh.Signature, error) {
	signers, err := a.Signers()
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	for _, s := range signers {
		if bytes.Equal(s.PublicKey().Marshal(), key.Marshal()) {
			signer = s
		}
	}
	if signer == nil {
		return nil, errors.Errorf("no signer found for %s", key.Marshal())
	}

	alg := key.Type()
	switch {
	case alg == ssh.KeyAlgoRSA && flags&sshagent.SignatureFlagRsaSha256 != 0:
		alg = ssh.KeyAlgoRSASHA256
	case alg == ssh.KeyAlgoRSA && flags&sshagent.SignatureFlagRsaSha512 != 0:
		alg = ssh.KeyAlgoRSASHA512
	}

	return signer.(ssh.AlgorithmSigner).SignWithAlgorithm(rand.Reader, data, alg)
}

func (a agent) Extension(_ string, _ []byte) ([]byte, error) {
	return nil, sshagent.ErrExtensionUnsupported
}

func (a agent) getAllKeys() ([]gostore.Key, error) {
	var ret []gostore.Key
	for storeID, store := range a.stores {
		if len(store.IDs) == 0 {
			continue // skip empty stores
		}

		keys, err := a.api.GetKeys(storeID, store.IDs)
		if err != nil {
			return nil, err
		}
		ret = append(ret, keys...)
	}

	return ret, nil
}
