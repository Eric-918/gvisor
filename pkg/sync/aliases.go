// Copyright 2020 The gVisor Authors.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"
)

// Aliases of standard library types.
type (
	// Mutex is an alias of sync.Mutex.
	Mutex = sync.Mutex

	// RWMutex is an alias of sync.RWMutex.
	RWMutex = sync.RWMutex

	// Cond is an alias of sync.Cond.
	Cond = sync.Cond

	// Locker is an alias of sync.Locker.
	Locker = sync.Locker

	// Once is an alias of sync.Once.
	Once = sync.Once

	// Pool is an alias of sync.Pool.
	Pool = sync.Pool

	// WaitGroup is an alias of sync.WaitGroup.
	WaitGroup = sync.WaitGroup

	// Map is an alias of sync.Map.
	Map = sync.Map
)
