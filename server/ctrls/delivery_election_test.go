package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/tendervoting/server/confs"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryFailOnSignature(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	// change uuid
	ed.ID = uuid.NewV4().String()
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestDeliverySuccessfulOnSignature(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = ed.From

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeOK, resp.Code)
}

func TestElectionDeliveryFailOnNotGonverment(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestElectionDeliveryFailOnNonHexVoter(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	// we will use the public key of the gonverment
	ed.Voters = []string{ed.From + "."}
	confs.Conf.GonvermentPublicKeyHex = ed.From

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestElectionDeliveryFailOnNonPublicKeyVoter(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	ed.Voters = []string{hex.EncodeToString([]byte("."))}
	confs.Conf.GonvermentPublicKeyHex = ed.From

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestElectionDeliveryFailOnTwiceTheSameVoter(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	// we will use the public key of the gonverment
	ed.Voters = []string{ed.From, ed.From}
	confs.Conf.GonvermentPublicKeyHex = ed.From

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestElectionDeliveryFailOnPuttingTheSameElectionID(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = ed.From

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	tvd.Data = &ed

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeOK, resp.Code)

	resp = app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}
