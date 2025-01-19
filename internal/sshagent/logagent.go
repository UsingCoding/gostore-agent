package sshagent

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

func NewLogAgent(next sshagent.ExtendedAgent, logger logrus.FieldLogger) sshagent.ExtendedAgent {
	return &logAgent{next: next, logger: logger}
}

type logAgent struct {
	next   sshagent.ExtendedAgent
	logger logrus.FieldLogger
}

func (a logAgent) List() ([]*sshagent.Key, error) {
	list, err := a.next.List()
	a.log("List", nil, err)
	return list, err
}

func (a logAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	sign, err := a.next.Sign(key, data)
	a.log("Sign", nil, err)
	return sign, err
}

func (a logAgent) Add(key sshagent.AddedKey) error {
	err := a.next.Add(key)
	a.log("Add", []any{key}, err)
	return err
}

func (a logAgent) Remove(key ssh.PublicKey) error {
	err := a.next.Remove(key)
	a.log("Remove", nil, err)
	return err
}

func (a logAgent) RemoveAll() error {
	err := a.next.RemoveAll()
	a.log("RemoveAll", nil, err)
	return err
}

func (a logAgent) Lock(passphrase []byte) error {
	err := a.next.Lock(passphrase)
	a.log("Lock", []any{passphrase}, err)
	return err
}

func (a logAgent) Unlock(passphrase []byte) error {
	err := a.next.Unlock(passphrase)
	a.log("Unlock", []any{passphrase}, err)
	return err
}

func (a logAgent) Signers() ([]ssh.Signer, error) {
	signers, err := a.next.Signers()
	a.log("Signers", []any{signers}, err)
	return signers, err
}

func (a logAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags sshagent.SignatureFlags) (*ssh.Signature, error) {
	withFlags, err := a.next.SignWithFlags(key, data, flags)
	a.log("SignWithFlags", []any{key, data, flags}, err)
	return withFlags, err
}

func (a logAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	extension, err := a.next.Extension(extensionType, contents)
	a.log("Extension", []any{extension}, err)
	return extension, err
}

func (a logAgent) log(method string, args []interface{}, respErr error) {
	a.logger.
		WithFields(logrus.Fields{
			"method": method,
			"args":   fmt.Sprint(args),
			"err":    respErr,
		}).
		Debugln()
}
