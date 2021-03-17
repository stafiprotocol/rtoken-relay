package chains

// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

import (
	"github.com/stafiprotocol/rtoken-relay/core"
)

type Router interface {
	Send(msg *core.Message) error
}
