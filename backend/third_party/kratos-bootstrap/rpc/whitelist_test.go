package rpc

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWhiteListBasic(t *testing.T) {
	// ensure clean state
	ClearWhiteList()

	AddWhiteList("Health.Check", "Public.Ping")

	if !DefaultWhiteList.IsWhitelisted("Health.Check") {
		t.Fatalf("expected Health.Check to be whitelisted")
	}
	if !DefaultWhiteList.IsWhitelisted("Public.Ping") {
		t.Fatalf("expected Public.Ping to be whitelisted")
	}
	// unknown should not be whitelisted
	if DefaultWhiteList.IsWhitelisted("Other.Op") {
		t.Fatalf("expected Other.Op NOT to be whitelisted")
	}

	// MatchFunc should return false (skip middleware) when whitelisted
	m := NewWhiteListMatcher()
	if m(context.Background(), "Health.Check") {
		t.Fatalf("expected matcher to be false for whitelisted op")
	}
	if !m(context.Background(), "Other.Op") {
		t.Fatalf("expected matcher to be true for non-whitelisted op")
	}

	ClearWhiteList()
	if DefaultWhiteList.IsWhitelisted("Health.Check") {
		t.Fatalf("expected Health.Check to be cleared")
	}
}

func TestWhiteListNormalization(t *testing.T) {
	ClearWhiteList()
	// add with leading slash (as gRPC may provide)
	AddWhiteList("/pkg.Service/MethodX")

	// normalized lookup should work for various representations
	cases := []string{"/pkg.Service/MethodX", "pkg.Service/MethodX", "MethodX"}
	for _, c := range cases {
		if !DefaultWhiteList.IsWhitelisted(c) {
			t.Fatalf("expected %q to be whitelisted (normalization/fallback)", c)
		}
	}

	ClearWhiteList()
}

func TestWhiteListConcurrent(t *testing.T) {
	ClearWhiteList()
	SetWhiteList([]string{"A", "B", "C"})

	var wg sync.WaitGroup
	const goroutines = 50
	const iterations = 200

	// readers
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = DefaultWhiteList.IsWhitelisted("A")
				_ = DefaultWhiteList.IsWhitelisted("NonExistent")
				_ = NewWhiteListMatcher()(context.Background(), "B")
			}
		}(i)
	}

	// writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for k := 0; k < 50; k++ {
				AddWhiteList("X", "Y")
				SetWhiteList([]string{"A", "Z"})
				ClearWhiteList()
				SetWhiteList([]string{"A", "B", "C"})
				// small sleep to increase interleaving
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}
