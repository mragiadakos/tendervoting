package ctrls

import (
	"github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var _ types.Application = (*TVApplication)(nil)

type TVApplication struct {
	types.BaseApplication

	state State
}

func NewTVApplication() *TVApplication {
	state := loadState(dbm.NewMemDB())
	return &TVApplication{state: state}
}
