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

package kusto

import (
	"sync"
	"time"
)

// Pool keeps one Client per ProviderConfig. An SDK client caches HTTP
// connections and tokens, so rebuilding it on every reconcile would be slow
// and hammer the token endpoint. Entries are rebuilt when the config
// fingerprint changes and evicted lazily after being idle.
type Pool struct {
	mu      sync.Mutex
	entries map[string]*poolEntry
	idle    time.Duration
	newFn   func(Config) (Client, error)
	now     func() time.Time
}

type poolEntry struct {
	client   Client
	fp       string
	lastUsed time.Time
}

// NewPool creates a pool. newFn builds clients; idle is the eviction timeout.
func NewPool(idle time.Duration, newFn func(Config) (Client, error)) *Pool {
	if idle <= 0 {
		idle = 30 * time.Minute
	}
	return &Pool{entries: map[string]*poolEntry{}, idle: idle, newFn: newFn, now: time.Now}
}

// Get returns the client for key, creating or replacing it if the config
// fingerprint changed.
func (p *Pool) Get(key string, cfg Config) (Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	p.evictLocked(now)
	fp := cfg.Fingerprint()
	if e, ok := p.entries[key]; ok && e.fp == fp {
		e.lastUsed = now
		return e.client, nil
	}
	c, err := p.newFn(cfg)
	if err != nil {
		return nil, err
	}
	if old, ok := p.entries[key]; ok {
		closeClient(old.client)
	}
	p.entries[key] = &poolEntry{client: c, fp: fp, lastUsed: now}
	return c, nil
}

// Remove drops the client for key.
func (p *Pool) Remove(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if e, ok := p.entries[key]; ok {
		closeClient(e.client)
		delete(p.entries, key)
	}
}

// Len returns the number of pooled clients.
func (p *Pool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.entries)
}

func (p *Pool) evictLocked(now time.Time) {
	for k, e := range p.entries {
		if now.Sub(e.lastUsed) > p.idle {
			closeClient(e.client)
			delete(p.entries, k)
		}
	}
}

func closeClient(c Client) {
	if cl, ok := c.(Closer); ok {
		_ = cl.Close()
	}
}
