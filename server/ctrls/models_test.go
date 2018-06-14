package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/stretchr/testify/assert"
)

func TestModelTVDeliverySuccessfulSignature(t *testing.T) {
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	ed := ElectionDeliveryData{}
	pubB, _ := privk.GetPublic().Bytes()
	ed.From = hex.EncodeToString(pubB)
	b, _ := json.Marshal(ed)
	er := TVDelivery{}
	er.Type = ELECTION
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
	b, _ := json.Marshal(pd)
	pr := TVDelivery{}
	pr.Type = POLL
	pr.Signature, _ = privk.Sign(b)
	pr.Data = &pd
	ver, err := pr.VerifySignature()
	assert.Nil(t, err)
	assert.Equal(t, true, ver)
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
