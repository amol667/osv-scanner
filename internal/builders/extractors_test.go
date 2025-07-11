package builders_test

import (
	"testing"

	"github.com/google/osv-scalibr/extractor/filesystem/language/golang/gomod"
	"github.com/google/osv-scalibr/extractor/filesystem/language/javascript/packagelockjson"
	"github.com/google/osv-scanner/v2/internal/builders"
)

func TestBuildExtractors(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name           string
		extractorNames []string
		expectedCount  int
	}{
		{
			name:           "valid extractors",
			extractorNames: []string{gomod.Name, packagelockjson.Name},
			expectedCount:  2,
		},
		{
			name:           "mixed valid and invalid",
			extractorNames: []string{gomod.Name, "unknown-extractor"},
			expectedCount:  1,
		},
		{
			name:           "all invalid",
			extractorNames: []string{"unknown1", "unknown2"},
			expectedCount:  0,
		},
		{
			name:           "empty list",
			extractorNames: []string{},
			expectedCount:  0,
		},
		{
			name:           "duplicate extractors",
			extractorNames: []string{gomod.Name, gomod.Name},
			expectedCount:  2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			extractors := builders.BuildExtractors(tt.extractorNames)
			
			if len(extractors) != tt.expectedCount {
				t.Errorf("got %d extractors, want %d", len(extractors), tt.expectedCount)
			}
			
			for i, extractor := range extractors {
				if extractor == nil {
					t.Errorf("extractor at index %d is nil", i)
				}
			}
		})
	}
}

func TestBuildExtractors_KnownExtractors(t *testing.T) {
	t.Parallel()
	
	knownExtractors := []string{
		gomod.Name,
		packagelockjson.Name,
	}
	
	extractors := builders.BuildExtractors(knownExtractors)
	
	if len(extractors) != len(knownExtractors) {
		t.Errorf("got %d extractors, want %d", len(extractors), len(knownExtractors))
	}
	
	for i, extractor := range extractors {
		if extractor == nil {
			t.Errorf("known extractor %q returned nil", knownExtractors[i])
		}
	}
}

func TestBuildExtractors_SingleExtractor(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name          string
		extractorName string
		shouldExist   bool
	}{
		{"go.mod", gomod.Name, true},
		{"package-lock.json", packagelockjson.Name, true},
		{"unknown", "completely-unknown-extractor", false},
		{"empty string", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			extractors := builders.BuildExtractors([]string{tt.extractorName})
			
			if tt.shouldExist {
				if len(extractors) != 1 {
					t.Errorf("got %d extractors, want 1", len(extractors))
				}
				if len(extractors) > 0 && extractors[0] == nil {
					t.Error("extractor should not be nil")
				}
			} else {
				if len(extractors) != 0 {
					t.Errorf("got %d extractors, want 0", len(extractors))
				}
			}
		})
	}
}

func TestBuildExtractors_CaseInsensitive(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name          string
		extractorName string
		shouldExist   bool
	}{
		{"uppercase", "GO.MOD", false},
		{"mixed case", "Package-Lock.Json", false},
		{"exact case", gomod.Name, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			extractors := builders.BuildExtractors([]string{tt.extractorName})
			
			if tt.shouldExist {
				if len(extractors) != 1 {
					t.Errorf("got %d extractors, want 1", len(extractors))
				}
			} else {
				if len(extractors) != 0 {
					t.Errorf("got %d extractors, want 0", len(extractors))
				}
			}
		})
	}
}
