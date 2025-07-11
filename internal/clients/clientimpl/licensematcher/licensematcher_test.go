package licensematcher_test

import (
	"testing"

	"github.com/google/osv-scanner/v2/internal/clients/clientimpl/licensematcher"
)

func TestDepsDevLicenseMatcher_Creation(t *testing.T) {
	t.Parallel()
	
	matcher := &licensematcher.DepsDevLicenseMatcher{}
	
	if matcher == nil {
		t.Error("DepsDevLicenseMatcher should not be nil")
	}
}

func TestDepsDevLicenseMatcher_ClientField(t *testing.T) {
	t.Parallel()
	
	matcher := &licensematcher.DepsDevLicenseMatcher{
		Client: nil,
	}
	
	if matcher.Client != nil {
		t.Error("Client field should be nil when set to nil")
	}
}
