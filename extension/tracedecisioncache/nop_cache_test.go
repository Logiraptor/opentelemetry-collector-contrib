// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracedecisioncache

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestNopCache(t *testing.T) {
	c := NewNopDecisionCache[bool]()
	id, err := traceIDFromHex("12341234123412341234123412341234")
	require.NoError(t, err)
	c.Put(id, true)
	v, ok := c.Get(id)
	assert.False(t, v)
	assert.False(t, ok)
}

func traceIDFromHex(idStr string) (pcommon.TraceID, error) {
	id := pcommon.NewTraceIDEmpty()
	_, err := hex.Decode(id[:], []byte(idStr))
	return id, err
}
