package gostore

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func NewLogAPI(next API, logger logrus.FieldLogger) API {
	return &logAPI{next: next, logger: logger}
}

type logAPI struct {
	next   API
	logger logrus.FieldLogger
}

func (l logAPI) GetKeys(storeID string, secretIDs []string) ([]Key, error) {
	keys, err := l.next.GetKeys(storeID, secretIDs)
	l.logger.
		WithFields(logrus.Fields{
			"target":    "gostore-api",
			"method":    "GetKeys",
			"storeID":   storeID,
			"secretIDs": secretIDs,
			"err":       err,
		}).
		Debugln()
	return keys, err
}

func (l logAPI) GetPrivateKey(storeID, secretID string) (ssh.Signer, error) {
	key, err := l.next.GetPrivateKey(storeID, secretID)
	l.logger.
		WithFields(logrus.Fields{
			"target":   "gostore-api",
			"method":   "GetPrivateKey",
			"storeID":  storeID,
			"secretID": secretID,
			"err":      err,
		}).
		Debugln()
	return key, err
}
