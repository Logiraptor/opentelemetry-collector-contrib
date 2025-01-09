// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracedecisioncache // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/cache"

import (
	"encoding/json"

	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// CacheExtension is an extension that caches sampling decisions.
type CacheExtension interface {
	extension.Extension

	Cache[[]byte]
}

// Cache is a cache using a pcommon.TraceID as the key and any generic type as the value.
type Cache[V any] interface {
	// Get returns the value for the given id, and a boolean to indicate whether the key was found.
	// If the key is not present, the zero value is returned.
	Get(id pcommon.TraceID) (V, bool)
	// Put sets the value for a given id
	Put(id pcommon.TraceID, v V)
	// Delete deletes the value for the given id
	Delete(id pcommon.TraceID)
}

type Codec[V any] interface {
	// Encode encodes the value into a byte slice.
	Encode(v V) ([]byte, error)
	// Decode decodes the byte slice into a value.
	Decode(data []byte) (V, error)
}

func NewTypedCache[V any](codec Codec[V], cache Cache[[]byte]) Cache[V] {
	return &typedCache[V]{cache: cache, codec: codec}
}

func NewJsonCodec[V any]() Codec[V] {
	return jsonCodec[V]{}
}

type jsonCodec[V any] struct{}

func (jsonCodec[V]) Encode(v V) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec[V]) Decode(data []byte) (V, error) {
	var v V
	err := json.Unmarshal(data, &v)
	return v, err
}

type typedCache[V any] struct {
	cache Cache[[]byte]
	codec Codec[V]
}

func (c *typedCache[V]) Get(id pcommon.TraceID) (V, bool) {
	buf, ok := c.cache.Get(id)
	if !ok {
		var zero V
		return zero, false
	}

	v, err := c.codec.Decode(buf)
	if err != nil {
		var zero V
		return zero, false
	}

	return v, true
}

func (c *typedCache[V]) Put(id pcommon.TraceID, v V) {
	buf, err := c.codec.Encode(v)
	if err != nil {
		return
	}

	c.cache.Put(id, buf)
}

func (c *typedCache[V]) Delete(id pcommon.TraceID) {
	c.cache.Delete(id)
}
