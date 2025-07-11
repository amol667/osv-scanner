package cachedregexp_test

import (
	"regexp"
	"sync"
	"testing"

	"github.com/google/osv-scanner/v2/internal/cachedregexp"
)

func TestMustCompile(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name    string
		pattern string
	}{
		{"simple pattern", "test"},
		{"complex pattern", `\d+\.\d+\.\d+`},
		{"anchored pattern", "^start.*end$"},
		{"character class", "[a-zA-Z0-9]+"},
		{"escaped characters", `\.\*\+\?\[\]\(\)\{\}\|\^`},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			regex := cachedregexp.MustCompile(tt.pattern)
			if regex == nil {
				t.Fatal("MustCompile returned nil")
			}
			
			expected := regexp.MustCompile(tt.pattern)
			if regex.String() != expected.String() {
				t.Errorf("got %q, want %q", regex.String(), expected.String())
			}
		})
	}
}

func TestMustCompile_Caching(t *testing.T) {
	t.Parallel()
	
	pattern := "test-pattern-for-caching"
	
	regex1 := cachedregexp.MustCompile(pattern)
	regex2 := cachedregexp.MustCompile(pattern)
	
	if regex1 != regex2 {
		t.Error("MustCompile should return cached instance for same pattern")
	}
}

func TestMustCompile_DifferentPatterns(t *testing.T) {
	t.Parallel()
	
	pattern1 := "pattern-one"
	pattern2 := "pattern-two"
	
	regex1 := cachedregexp.MustCompile(pattern1)
	regex2 := cachedregexp.MustCompile(pattern2)
	
	if regex1 == regex2 {
		t.Error("MustCompile should return different instances for different patterns")
	}
}

func TestMustCompile_Concurrent(t *testing.T) {
	t.Parallel()
	
	pattern := "concurrent-test-pattern"
	var wg sync.WaitGroup
	results := make([]*regexp.Regexp, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = cachedregexp.MustCompile(pattern)
		}(i)
	}
	
	wg.Wait()
	
	for i := 1; i < len(results); i++ {
		if results[0] != results[i] {
			t.Error("Concurrent MustCompile calls should return same cached instance")
		}
	}
}

func TestMustCompile_ConcurrentDifferentPatterns(t *testing.T) {
	t.Parallel()
	
	patterns := []string{
		"pattern-a",
		"pattern-b", 
		"pattern-c",
		"pattern-d",
		"pattern-e",
	}
	
	var wg sync.WaitGroup
	results := make(map[string]*regexp.Regexp)
	mu := sync.Mutex{}
	
	for _, pattern := range patterns {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			regex := cachedregexp.MustCompile(p)
			mu.Lock()
			results[p] = regex
			mu.Unlock()
		}(pattern)
	}
	
	wg.Wait()
	
	if len(results) != len(patterns) {
		t.Errorf("got %d results, want %d", len(results), len(patterns))
	}
	
	for pattern, regex := range results {
		if regex == nil {
			t.Errorf("pattern %q returned nil regex", pattern)
		}
		
		cached := cachedregexp.MustCompile(pattern)
		if regex != cached {
			t.Errorf("pattern %q: cached instance differs from concurrent result", pattern)
		}
	}
}
