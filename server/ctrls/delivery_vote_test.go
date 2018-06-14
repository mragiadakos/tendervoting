package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/stretchr/testify/assert"
)

func TestVoteFailOnEmptyPollHash(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)

	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)

}

func TestVoteFailOnPollHashDoesNotExists(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	vd.PollHash = "lalallafakehahahash"

	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestVoteFailOnVoterNotInTheElection(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	otherVoter, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	otherPubB, _ := otherVoter.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	vd.PollHash = forTestCreatePoll(t, app, privk, []string{hex.EncodeToString(otherPubB)}, map[string]string{"k": "k"})

	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestVoteFailOnChoiceThatDoesNotexists(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	vd.PollHash = forTestCreatePoll(t, app, privk, []string{vd.From}, map[string]string{"k": "k"})
	vd.Choice = "non-existent-choice"
	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestVoteFailOnReVotingOnTheSamePoll(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	pollHash := forTestCreatePoll(t, app, privk, []string{vd.From}, map[string]string{"a": "a", "b": "b"})

	vd.PollHash = pollHash
	vd.Choice = "a"
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

	// voting second time
	vd = VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	vd.PollHash = pollHash
	vd.Choice = "b"
	b, _ = json.Marshal(vd)
	sign, err = privk.Sign(b)
	assert.Nil(t, err)

	tvd = TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ = json.Marshal(tvd)
	resp = app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestVoteFailOnNotLatestPoll(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	oldPollHash := forTestCreatePoll(t, app, privk, []string{vd.From}, map[string]string{"a": "a", "b": "b"})
	// new poll hash
	forTestCreatePoll(t, app, privk, []string{vd.From}, map[string]string{"1": "1", "2": "2"})

	vd.PollHash = oldPollHash
	vd.Choice = "a"
	b, _ := json.Marshal(vd)
	sign, err := privk.Sign(b)
	assert.Nil(t, err)

	tvd := TVDelivery{}
	tvd.Type = VOTE
	tvd.Signature = sign
	tvd.Data = &vd

	tx, _ := json.Marshal(tvd)
	resp := app.DeliverTx(tx)
	assert.Equal(t, CodeTypeUnauthorized, resp.Code)
}

func TestVoteSuccessful(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()

	vd := VoteDeliveryData{}
	vd.From = hex.EncodeToString(pubB)
	pollHash := forTestCreatePoll(t, app, privk, []string{vd.From}, map[string]string{"a": "a", "b": "b"})

	vd.PollHash = pollHash
	vd.Choice = "a"
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
	ps, _ := app.state.GetPoll(pollHash)
	assert.Equal(t, 1, ps.Choices[vd.Choice])
}
