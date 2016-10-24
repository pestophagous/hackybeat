package deduper

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	lpkg "github.com/pestophagous/hackybeat/util/logger"
)

type serializer int

const (
	gobber serializer = iota
	jsoner
	md5er
)

const defaultSerializer = jsoner // gobber

var whichSerializer = defaultSerializer

type blobifier interface {
	makeBlob(interface{}) []byte
}

func getSerializer(log *lpkg.LogWithNilCheck) blobifier {
	switch whichSerializer {
	case gobber:
		return &gobMaker{logger: log}
	case jsoner:
		return &jsonMaker{logger: log}
	case md5er:
		panic("md5er style of serializer not yet provided.")
	}
	panic("should have exited function during switch/case above.")
	return nil
}

type gobMaker struct {
	logger *lpkg.LogWithNilCheck
}
type jsonMaker struct {
	logger *lpkg.LogWithNilCheck
}

func (this *gobMaker) InstanceIdForLogging() string {
	return "gobMaker"
}

func (this *jsonMaker) InstanceIdForLogging() string {
	return "jsonMaker"
}

func (this *gobMaker) makeBlob(object interface{}) []byte {
	var blob bytes.Buffer
	enc := gob.NewEncoder(&blob)
	err := enc.Encode(object)
	this.logger.LogPossibleFailureOf("enc.Encode blob", this, err)
	return blob.Bytes()
}

func (this *jsonMaker) makeBlob(object interface{}) []byte {
	rslt, err := json.Marshal(object)
	this.logger.LogPossibleFailureOf("json.Marshal", this, err)
	return rslt
}
