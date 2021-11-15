// Copyright 2021 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	chtest "perun.network/go-perun/channel/test"
	"perun.network/go-perun/client"
	ctest "perun.network/go-perun/client/test"
	"perun.network/go-perun/wire"
	"polycry.pt/poly-go/test"
)

func TestSubChannelDispute(t *testing.T) {
	rng := test.Prng(t)

	setups := NewSetups(rng, []string{"DisputeSusie", "DisputeTim"})
	roles := [2]ctest.Executer{
		ctest.NewDisputeSusie(t, setups[0]),
		ctest.NewDisputeTim(t, setups[1]),
	}

	baseCfg := ctest.MakeBaseExecConfig(
		[2]wire.Address{setups[0].Identity.Address(), setups[1].Identity.Address()},
		chtest.NewRandomAsset(rng),
		[2]*big.Int{big.NewInt(100), big.NewInt(100)},
		client.WithoutApp(),
	)
	cfg := &ctest.DisputeSusieTimExecConfig{
		BaseExecConfig:  baseCfg,
		SubChannelFunds: [2]*big.Int{big.NewInt(10), big.NewInt(10)},
		TxAmount:        big.NewInt(1),
	}

	ctx, cancel := context.WithTimeout(context.Background(), twoPartyTestTimeout)
	defer cancel()
	assert.NoError(t, ctest.ExecuteTwoPartyTest(ctx, roles, cfg))
}
