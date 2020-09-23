package signer

import (
	"sync"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

// PvGuard guards access to an underlying PrivValidator by using mutexes
// for each of the PrivValidator interface functions
type PvGuard struct {
	PrivValidator types.PrivValidator
	pvMutex       sync.Mutex
}

// GetPubKey implementes types.PrivValidator
func (pv *PvGuard) GetPubKey() (crypto.PubKey, error) {
	pv.pvMutex.Lock()
	defer pv.pvMutex.Unlock()
	return pv.PrivValidator.GetPubKey()
}

// SignVote implementes types.PrivValidator
func (pv *PvGuard) SignVote(chainID string, vote *types.Vote) error {
	pv.pvMutex.Lock()
	defer pv.pvMutex.Unlock()
	return pv.PrivValidator.SignVote(chainID, vote)
}

// SignProposal implementes types.PrivValidator
func (pv *PvGuard) SignProposal(chainID string, proposal *types.Proposal) error {
	pv.pvMutex.Lock()
	defer pv.pvMutex.Unlock()
	return pv.PrivValidator.SignProposal(chainID, proposal)
}
