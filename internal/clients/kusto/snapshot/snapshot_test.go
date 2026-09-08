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

package snapshot

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheHitMissTTL(t *testing.T) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	c := New(time.Minute, false, WithClock(func() time.Time { return now }))
	key := Key{Endpoint: "e", Database: "db", Section: "tables"}
	loads := 0
	load := func(context.Context) (map[string]string, error) {
		loads++
		return map[string]string{"T": "x"}, nil
	}
	v, err := Load(context.Background(), c, key, load)
	if err != nil || v["T"] != "x" || loads != 1 {
		t.Fatalf("first load: %v %v %d", v, err, loads)
	}
	if _, err := Load(context.Background(), c, key, load); err != nil || loads != 1 {
		t.Fatalf("expected hit, loads=%d", loads)
	}
	now = now.Add(61 * time.Second)
	if _, err := Load(context.Background(), c, key, load); err != nil || loads != 2 {
		t.Fatalf("expected reload after TTL, loads=%d", loads)
	}
	c.Invalidate("e", "db", "tables")
	if _, err := Load(context.Background(), c, key, load); err != nil || loads != 3 {
		t.Fatalf("expected reload after invalidate, loads=%d", loads)
	}
	// Invalidate everything for the database.
	Load(context.Background(), c, Key{"e", "db", "functions"}, load) //nolint:errcheck // exercised above
	c.Invalidate("e", "db")
	if c.Len() != 0 {
		t.Errorf("Len after db invalidate = %d", c.Len())
	}
	c.Put(key, map[string]string{"P": "y"})
	if v, _ := Load(context.Background(), c, key, load); v["P"] != "y" {
		t.Error("Put not visible")
	}
}

func TestCacheErrorsNotCached(t *testing.T) {
	c := New(time.Minute, false)
	key := Key{"e", "db", "tables"}
	calls := 0
	load := func(context.Context) (int, error) {
		calls++
		if calls == 1 {
			return 0, errors.New("boom")
		}
		return 7, nil
	}
	if _, err := Load(context.Background(), c, key, load); err == nil {
		t.Fatal("expected error")
	}
	if v, err := Load(context.Background(), c, key, load); err != nil || v != 7 {
		t.Fatalf("second load: %v %v", v, err)
	}
}

func TestCacheDisabled(t *testing.T) {
	c := New(time.Minute, true)
	if !c.Disabled() {
		t.Fatal("expected disabled")
	}
	key := Key{"e", "db", "tables"}
	calls := 0
	load := func(context.Context) (int, error) { calls++; return calls, nil }
	Load(context.Background(), c, key, load) //nolint:errcheck // exercised below
	if v, _ := Load(context.Background(), c, key, load); v != 2 || c.Len() != 0 {
		t.Errorf("disabled cache must always load: v=%d len=%d", v, c.Len())
	}
	c.Put(key, 99)
	if c.Len() != 0 {
		t.Error("Put on disabled cache must be a no-op")
	}
	if New(0, false).Disabled() != true {
		t.Error("ttl 0 must disable")
	}
}

func TestCacheSingleflight(t *testing.T) {
	c := New(time.Minute, false)
	key := Key{"e", "db", "tables"}
	var loads int32
	release := make(chan struct{})
	load := func(context.Context) (int, error) {
		atomic.AddInt32(&loads, 1)
		<-release
		return 1, nil
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Load(context.Background(), c, key, load) //nolint:errcheck // asserted via loads counter
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()
	if got := atomic.LoadInt32(&loads); got != 1 {
		t.Errorf("expected one load, got %d", got)
	}
}

func TestCacheLoaderDetachedFromCaller(t *testing.T) {
	c := New(time.Minute, false)
	key := Key{"e", "db", "tables"}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	loaderSawCancel := false
	_, err := Load(ctx, c, key, func(lctx context.Context) (int, error) {
		loaderSawCancel = lctx.Err() != nil
		return 1, nil
	})
	if loaderSawCancel {
		t.Error("loader must not inherit caller cancellation")
	}
	if err == nil {
		t.Error("caller with cancelled context must get an error")
	}
	// The value was still cached for the others.
	if c.Len() != 1 {
		t.Error("value should be cached")
	}
}

func TestTypeError(t *testing.T) {
	c := New(time.Minute, false)
	key := Key{"e", "db", "tables"}
	c.Put(key, "a string")
	_, err := Load(context.Background(), c, key, func(context.Context) (int, error) { return 0, nil })
	var te *TypeError
	if !errors.As(err, &te) {
		t.Errorf("expected TypeError, got %v", err)
	}
}
