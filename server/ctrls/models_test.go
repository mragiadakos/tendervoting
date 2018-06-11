package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/stretchr/testify/assert"
)

func TestModelTVDeliverySuccessfulSignature(t *testing.T) {
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	ed := ElectionDeliveryData{}
	pubB, _ := privk.GetPublic().Bytes()
	ed.From = hex.EncodeToString(pubB)
	ed.StartTime = time.Now()
	ed.EndTime = time.Now()
	b, _ := json.Marshal(ed)
	er := TVDelivery{}
	er.Signature, _ = privk.Sign(b)
	er.Data = &ed
	ver, err := er.VerifySignature()
	assert.Nil(t, err)
	assert.Equal(t, true, ver)
}

func TestModelTVDeliveryFailSignature(t *testing.T) {
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	pd := PollDeliveryData{}
	pubB, _ := privk.GetPublic().Bytes()
	pd.From = hex.EncodeToString(pubB)
	pd.StartTime = time.Now()
	pd.EndTime = time.Now()
	b, _ := json.Marshal(pd)
	pr := TVDelivery{}
	pr.Signature, _ = privk.Sign(b)
	pr.Data = &pd
	ver, err := pr.VerifySignature()
	assert.Nil(t, err)
	assert.Equal(t, true, ver)
}

func TestModelElectionDeliveryFailOnTime(t *testing.T) {
	ed := ElectionDeliveryData{}
	ed.EndTime = time.Now()
	time.Sleep(1 * time.Second)
	ed.StartTime = time.Now()
	err := ed.ValidateTime()
	assert.NotNil(t, err)
}

func TestModelPollDeliveryFailOnTime(t *testing.T) {
	pd := PollDeliveryData{}
	pd.EndTime = time.Now()
	time.Sleep(1 * time.Second)
	pd.StartTime = time.Now()
	err := pd.ValidateTime()
	assert.NotNil(t, err)
}

func TestModelElectionDeliveryFailOnNonGonverment(t *testing.T) {
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	ed := ElectionDeliveryData{}
	pubB, _ := privk.GetPublic().Bytes()
	ed.From = hex.EncodeToString(pubB)
	err = ed.ValidateGonverment()
	assert.NotNil(t, err)
}

func TestModelPollDeliveryFailOnNonGonverment(t *testing.T) {
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	pd := PollDeliveryData{}
	pubB, _ := privk.GetPublic().Bytes()
	pd.From = hex.EncodeToString(pubB)
	err = pd.ValidateGonverment()
	assert.NotNil(t, err)
}

func TestModelElectionDeliveryFailOnHexVoter(t *testing.T) {
	_, pubk, err := crypto.GenerateEd25519Key(rand.Reader)
	b, _ := pubk.Bytes()
	voterHex := hex.EncodeToString(b)
	ed := ElectionDeliveryData{}
	ed.Voters = append([]string{}, voterHex+".")
	err = ed.ValidateVoters()
	assert.NotNil(t, err)
}

func TestModelElectionDeliveryFailOnPublicKeyVoter(t *testing.T) {
	voterHex := hex.EncodeToString([]byte("."))
	ed := ElectionDeliveryData{}
	ed.Voters = append([]string{}, voterHex)
	err := ed.ValidateVoters()
	assert.NotNil(t, err)
}
