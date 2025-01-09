// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package lrutracedecisioncache // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/tracedecisioncache/lrutracedecisioncache"

import (
	"context"
	"encoding/binary"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/tracedecisioncache"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// lruDecisionCache implements Cache as a simple LRU cache.
// It holds trace IDs that had sampling decisions made on them.
// It does not specify the type of sampling decision that was made, only that
// a decision was made for an ID. You need separate DecisionCaches for caching
// sampled and not sampled trace IDs.
type lruDecisionCache struct {
	cache *lru.Cache[uint64, []byte]
}

var _ tracedecisioncache.CacheExtension = (*lruDecisionCache)(nil)

// NewLRUDecisionCache returns a new lruDecisionCache.
// The size parameter indicates the amount of keys the cache will hold before it
// starts evicting the least recently used key.
func NewLRUDecisionCache(size int) (tracedecisioncache.CacheExtension, error) {
	c, err := lru.New[uint64, []byte](size)
	if err != nil {
		return nil, err
	}
	return &lruDecisionCache{cache: c}, nil
}

func (c *lruDecisionCache) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (c *lruDecisionCache) Shutdown(ctx context.Context) error {
	return nil
}

func (c *lruDecisionCache) Get(id pcommon.TraceID) ([]byte, bool) {
	return c.cache.Get(rightHalfTraceID(id))
}

func (c *lruDecisionCache) Put(id pcommon.TraceID, v []byte) {
	_ = c.cache.Add(rightHalfTraceID(id), v)
}

// Delete is no-op since LRU relies on least recently used key being evicting automatically
func (c *lruDecisionCache) Delete(_ pcommon.TraceID) {}

func rightHalfTraceID(id pcommon.TraceID) uint64 {
	return binary.LittleEndian.Uint64(id[8:])
}
