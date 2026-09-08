/*
Copyright 2026 The provider-azure-adx Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package snapshot is the batch-observe cache. One reconcile of a Table does
// not run ".show table X"; it asks the cache for the "tables" section of its
// database, and a cache miss loads all tables of that database with a single
// command. With thousands of managed resources and a 10 minute poll this
// turns O(resources) commands into O(databases x sections) per TTL.
package snapshot

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/singleflight"
)

// Key identifies a cached section.
type Key struct {
	Endpoint string
	Database string
	Section  string
}

func (k Key) String() string { return k.Endpoint + "|" + k.Database + "|" + k.Section }

// Cache holds loaded sections with a TTL. It is safe for concurrent use.
type Cache struct {
	ttl      time.Duration
	disabled bool
	loadTO   time.Duration

	mu      sync.Mutex
	entries map[Key]entry
	sf      singleflight.Group
	now     func() time.Time
}

type entry struct {
	val any
	exp time.Time
}

// Option configures a Cache.
type Option func(*Cache)

// WithClock overrides the clock (tests).
func WithClock(now func() time.Time) Option { return func(c *Cache) { c.now = now } }

// WithLoadTimeout bounds one loader run (default 2 minutes).
func WithLoadTimeout(d time.Duration) Option { return func(c *Cache) { c.loadTO = d } }

// New creates a cache. ttl <= 0 or disabled=true turns caching off: every Get
// runs its loader (still deduplicated through singleflight).
func New(ttl time.Duration, disabled bool, opts ...Option) *Cache {
	c := &Cache{ttl: ttl, disabled: disabled || ttl <= 0, loadTO: 2 * time.Minute, entries: map[Key]entry{}, now: time.Now}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Disabled reports whether caching is off.
func (c *Cache) Disabled() bool { return c.disabled }

// Get returns the cached value for key or runs load. Concurrent callers of
// the same key share one load. The loader runs detached from the caller's
// cancellation so one aborted reconcile cannot fail the others waiting on it.
func (c *Cache) Get(ctx context.Context, key Key, load func(ctx context.Context) (any, error)) (any, error) {
	if !c.disabled {
		c.mu.Lock()
		e, ok := c.entries[key]
		now := c.now()
		c.mu.Unlock()
		if ok && now.Before(e.exp) {
			cacheRequests.WithLabelValues(key.Section, "hit").Inc()
			return e.val, nil
		}
	}
	cacheRequests.WithLabelValues(key.Section, "miss").Inc()
	v, err, _ := c.sf.Do(key.String(), func() (any, error) {
		lctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), c.loadTO)
		defer cancel()
		val, err := load(lctx)
		if err != nil {
			return nil, err
		}
		if !c.disabled {
			c.mu.Lock()
			c.entries[key] = entry{val: val, exp: c.now().Add(c.ttl)}
			c.mu.Unlock()
		}
		return val, nil
	})
	if err != nil {
		return nil, err
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return v, nil
}

// Put stores a value directly (e.g. after a write that returned fresh state).
func (c *Cache) Put(key Key, val any) {
	if c.disabled {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{val: val, exp: c.now().Add(c.ttl)}
}

// Invalidate drops the given sections of a database. With no sections every
// section of the database is dropped.
func (c *Cache) Invalidate(endpoint, database string, sections ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(sections) == 0 {
		for k := range c.entries {
			if k.Endpoint == endpoint && k.Database == database {
				delete(c.entries, k)
			}
		}
		return
	}
	for _, s := range sections {
		delete(c.entries, Key{Endpoint: endpoint, Database: database, Section: s})
	}
}

// Len returns the number of live entries (tests).
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Load is a typed wrapper around Cache.Get.
func Load[T any](ctx context.Context, c *Cache, key Key, load func(ctx context.Context) (T, error)) (T, error) {
	v, err := c.Get(ctx, key, func(ctx context.Context) (any, error) { return load(ctx) })
	if err != nil {
		var zero T
		return zero, err
	}
	t, ok := v.(T)
	if !ok {
		var zero T
		return zero, &TypeError{Key: key}
	}
	return t, nil
}

// TypeError is returned when a section holds a value of an unexpected type.
type TypeError struct{ Key Key }

func (e *TypeError) Error() string { return "snapshot: unexpected value type for " + e.Key.String() }

var cacheRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "adx_observe_cache_requests_total",
	Help: "Batch-observe cache lookups by section and result (hit|miss).",
}, []string{"section", "result"})

// RegisterMetrics registers the cache metrics with reg.
func RegisterMetrics(reg prometheus.Registerer) error {
	if err := reg.Register(cacheRequests); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok { //nolint:errorlint // prometheus returns the struct by value
			return err
		}
	}
	return nil
}
