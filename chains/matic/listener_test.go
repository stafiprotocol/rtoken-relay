// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"testing"

	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
)

func TestListener(t *testing.T) {
	ctx, err := createCliContext("", []string{config.ConfigFileFlag.Name}, []interface{}{"../../config_template_matic.json"})
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	chain := cfg.Chains[1]

	l := NewListener(chain.Name, core.RSymbol(chain.Rsymbol), chain.Opts, nil, 5, nil, nil, nil, nil)
	t.Log(l.eraSeconds, l.eraOffset)

}
