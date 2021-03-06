package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/satori/go.uuid"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/tendervoting/server/confs"
	"github.com/stretchr/testify/assert"
)

func TestPollDeliveryFailOnNonGonverment(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)

}

func TestPollDeliveryFailOnElectionIDDoesNotExists(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = uuid.NewV4().String()
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnEmptyPollHash(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnDoesNotHavePollJson(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	ioutil.WriteFile(tmpFolder+"/example", []byte("example"), 0755)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnPollJsonFormatError(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", []byte("example"), 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnEmptyDescription(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "", // empty description
		Choices: map[string]string{
			"k": "k",
		},
	}
	bpj, _ := json.Marshal(pj)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnEmptyChoices(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "k",
		// empty choices
	}
	bpj, _ := json.Marshal(pj)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailOnPollExistsAlready(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "k",
		Choices: map[string]string{
			"k": "k",
		},
	}
	bpj, _ := json.Marshal(pj)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeOK, resp.Code)

	resp = app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliveryFailoOnNotLatestElection(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	oldElectionID := forTestCreateElection(t, app, privk, []string{})
	forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = oldElectionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "k",
		Choices: map[string]string{
			"k": "k",
		},
	}
	bpj, _ := json.Marshal(pj)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestPollDeliverySuccesful(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	electionID := forTestCreateElection(t, app, privk, []string{})

	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "k",
		Choices: map[string]string{
			"k": "k",
		},
	}
	bpj, _ := json.Marshal(pj)
	err = ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
	assert.Nil(t, err)
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	hash, err := sh.AddDir(tmpFolder)
	assert.Nil(t, err)
	os.RemoveAll(tmpFolder)

	pd.PollHash = hash
	b, _ := json.Marshal(pd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	confs.Conf.GonvermentPublicKeyHex = pd.From

	tvd := TVDelivery{}
	tvd.Type = POLL
	tvd.Signature = sign
	tvd.Data = &pd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeOK, resp.Code)
}
