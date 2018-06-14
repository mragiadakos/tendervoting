package ctrls

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/tendervoting/server/confs"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func forTestCreateElection(t *testing.T, app *TVApplication, privk crypto.PrivKey, voters []string) string {
	pubB, _ := privk.GetPublic().Bytes()

	ed := ElectionDeliveryData{}
	ed.ID = uuid.NewV4().String()
	ed.From = hex.EncodeToString(pubB)
	ed.Voters = voters

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
	return ed.ID
}

func forTestCreatePoll(t *testing.T, app *TVApplication, privk crypto.PrivKey, electionID string, choices map[string]string) string {
	pubB, _ := privk.GetPublic().Bytes()
	pd := PollDeliveryData{}
	pd.From = hex.EncodeToString(pubB)
	pd.ElectionID = electionID

	// we will use a temporary folder
	tmpFolder := "temporary"
	os.MkdirAll(tmpFolder, 0755)
	pj := PollJson{
		Description: "k",
		Choices:     choices,
	}
	bpj, _ := json.Marshal(pj)
	err := ioutil.WriteFile(tmpFolder+"/poll.json", bpj, 0755)
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
	return pd.PollHash
}

func forTestCreateVote(t *testing.T, app *TVApplication, privk crypto.PrivKey, election, poll, choice string) {
	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	vd.PollHash = poll
	vd.Choice = choice
	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeOK, resp.Code)
}
