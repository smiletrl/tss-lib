// Copyright © 2019 Binance
//
// This file is part of Binance. The full Binance copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package keygen

import (
	"errors"
	"time"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/crypto/paillier"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

func (round *round4) Start() *tss.Error {
	if round.started {
		return round.WrapError(errors.New("round already started"))
	}
	round.logger.Infof("round 4 starts")
	start := time.Now()
	round.number = 4
	round.started = true
	round.resetOK()

	i := round.PartyID().Index
	Ps := round.Parties().IDs()
	PIDs := Ps.Keys()
	ecdsaPub := round.save.ECDSAPub

	// 1-3. (concurrent)
	// r3 messages are assumed to be available and != nil in this function
	r3msgs := round.temp.kgRound3Messages
	chs := make([]chan bool, len(r3msgs))
	for i := range chs {
		chs[i] = make(chan bool)
	}
	for j, msg := range round.temp.kgRound3Messages {
		if j == i {
			continue
		}
		r3msg := msg.Content().(*KGRound3Message)
		go func(prf paillier.Proof, j int, ch chan<- bool) {
			start2 := time.Now()
			ppk := round.save.PaillierPKs[j]
			ok, err := prf.Verify(ppk.N, PIDs[j], ecdsaPub)
			if err != nil {
				common.Logger.Error(round.WrapError(err, Ps[j]).Error())
				ch <- false
				return
			}
			ch <- ok
			round.logger.Infof("round 4 after paillier proof, taking: %d milliseconds", time.Since(start2).Milliseconds())
		}(r3msg.UnmarshalProofInts(), j, chs[j])
	}

	// consume unbuffered channels (end the goroutines)
	for j, ch := range chs {
		if j == i {
			round.ok[j] = true
			continue
		}
		round.ok[j] = <-ch
	}
	culprits := make([]*tss.PartyID, 0, len(Ps)) // who caused the error(s)
	for j, ok := range round.ok {
		if !ok {
			culprits = append(culprits, Ps[j])
			common.Logger.Warningf("paillier verify failed for party %s", Ps[j])
			continue
		}
		common.Logger.Debugf("paillier verify passed for party %s", Ps[j])

	}
	if len(culprits) > 0 {
		return round.WrapError(errors.New("paillier verify failed"), culprits...)
	}

	round.end <- round.save

	spent := time.Since(start).Milliseconds()
	// 10 seconds
	if spent > 10000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 10 seconds", time.Since(start).Milliseconds())
	} else if spent > 7000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 7 seconds", time.Since(start).Milliseconds())
	} else if spent > 4000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 4 seconds", time.Since(start).Milliseconds())
	} else if spent > 3000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 3 seconds", time.Since(start).Milliseconds())
	} else if spent > 2000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 2 seconds", time.Since(start).Milliseconds())
	} else if spent > 1000 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds over 1 seconds", time.Since(start).Milliseconds())
	} else if spent > 500 {
		round.logger.Infof("round 4 completes, taking: %d milliseconds under 1 second but above 0.5 second", time.Since(start).Milliseconds())
	} else {
		round.logger.Infof("round 4 completes, taking: %d milliseconds under 1 second and under 0.5 second", time.Since(start).Milliseconds())
	}
	return nil
}

func (round *round4) CanAccept(msg tss.ParsedMessage) bool {
	// not expecting any incoming messages in this round
	return false
}

func (round *round4) Update() (bool, *tss.Error) {
	// not expecting any incoming messages in this round
	return false, nil
}

func (round *round4) NextRound() tss.Round {
	return nil // finished!
}
