package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

	b, _ := json.Marshal(ed)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = ELECTION
	tvd.Signature = sign
	// change time
	ed.EndTime = time.Now().UTC()
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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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

func TestElectionDeliveryFailOnTime(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)
	ed.StartTime = time.Now().UTC()
	time.Sleep(1 * time.Second)
	ed.EndTime = ed.EndTime

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

func TestElectionDeliveryFailOnNonHexVoter(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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
	ed.StartTime = time.Now().UTC()
	ed.EndTime = time.Now().Add(1 * time.Hour).UTC()

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
