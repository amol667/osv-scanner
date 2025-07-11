package baseimagematcher_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/osv-scanner/v2/internal/clients/clientimpl/baseimagematcher"
	"github.com/google/osv-scanner/v2/pkg/models"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	
	config := baseimagematcher.DefaultConfig()
	
	if config.MaxRetryAttempts <= 0 {
		t.Error("MaxRetryAttempts should be positive")
	}
	
	if config.UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}
	
	if config.BackoffDurationMultiplier <= 0 {
		t.Error("BackoffDurationMultiplier should be positive")
	}
	
	if config.BackoffDurationExponential <= 0 {
		t.Error("BackoffDurationExponential should be positive")
	}
	
	if config.JitterMultiplier < 0 {
		t.Error("JitterMultiplier should be non-negative")
	}
	
	if config.MaxConcurrentBatchRequests <= 0 {
		t.Error("MaxConcurrentBatchRequests should be positive")
	}
}

func TestDepsDevBaseImageMatcher_Config(t *testing.T) {
	t.Parallel()
	
	config := baseimagematcher.DefaultConfig()
	config.MaxRetryAttempts = 1
	config.BackoffDurationMultiplier = 0.001
	
	matcher := &baseimagematcher.DepsDevBaseImageMatcher{
		HTTPClient: http.Client{Timeout: 5 * time.Second},
		Config:     config,
	}
	
	if matcher.Config.MaxRetryAttempts != 1 {
		t.Errorf("got MaxRetryAttempts %d, want 1", matcher.Config.MaxRetryAttempts)
	}
	
	if matcher.Config.BackoffDurationMultiplier != 0.001 {
		t.Errorf("got BackoffDurationMultiplier %f, want 0.001", matcher.Config.BackoffDurationMultiplier)
	}
}

func TestDepsDevBaseImageMatcher_HTTPClient(t *testing.T) {
	t.Parallel()
	
	config := baseimagematcher.DefaultConfig()
	
	matcher := &baseimagematcher.DepsDevBaseImageMatcher{
		HTTPClient: http.Client{Timeout: 5 * time.Second},
		Config:     config,
	}
	
	if matcher.HTTPClient.Timeout != 5*time.Second {
		t.Errorf("got timeout %v, want %v", matcher.HTTPClient.Timeout, 5*time.Second)
	}
}

func TestDepsDevBaseImageMatcher_ConfigValidation(t *testing.T) {
	t.Parallel()
	
	config := baseimagematcher.DefaultConfig()
	config.MaxRetryAttempts = 2
	config.BackoffDurationMultiplier = 0.001
	
	matcher := &baseimagematcher.DepsDevBaseImageMatcher{
		HTTPClient: http.Client{Timeout: 5 * time.Second},
		Config:     config,
	}
	
	if matcher.Config.MaxRetryAttempts != 2 {
		t.Errorf("got MaxRetryAttempts %d, want 2", matcher.Config.MaxRetryAttempts)
	}
	
	if matcher.Config.BackoffDurationMultiplier != 0.001 {
		t.Errorf("got BackoffDurationMultiplier %f, want 0.001", matcher.Config.BackoffDurationMultiplier)
	}
}

func TestDepsDevBaseImageMatcher_UserAgent(t *testing.T) {
	t.Parallel()
	
	config := baseimagematcher.DefaultConfig()
	
	if config.UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}
	
	if !strings.Contains(config.UserAgent, "osv-scanner") {
		t.Errorf("UserAgent should contain 'osv-scanner', got %q", config.UserAgent)
	}
}

func TestBuildBaseImageDetails(t *testing.T) {
	t.Parallel()
	
	layerMetadata := []models.LayerMetadata{
		{DiffID: "layer1"},
		{DiffID: "layer2"},
		{DiffID: "layer3"},
	}
	
	baseImagesToLayersMap := [][]models.BaseImageDetails{
		{},
		{{Name: "ubuntu:20.04"}},
		{{Name: "ubuntu:20.04"}},
	}
	
	if len(layerMetadata) != 3 {
		t.Errorf("got %d layers, want 3", len(layerMetadata))
	}
	
	if len(baseImagesToLayersMap) != 3 {
		t.Errorf("got %d base image maps, want 3", len(baseImagesToLayersMap))
	}
	
	if len(baseImagesToLayersMap[0]) != 0 {
		t.Errorf("first base image map should be empty, got %d items", len(baseImagesToLayersMap[0]))
	}
	
	if len(baseImagesToLayersMap[1]) != 1 || baseImagesToLayersMap[1][0].Name != "ubuntu:20.04" {
		t.Errorf("second base image map should contain ubuntu:20.04")
	}
}

func TestBuildBaseImageDetails_EmptyLayers(t *testing.T) {
	t.Parallel()
	
	layerMetadata := []models.LayerMetadata{}
	baseImagesToLayersMap := [][]models.BaseImageDetails{}
	
	if len(layerMetadata) != 0 {
		t.Errorf("got %d layers, want 0", len(layerMetadata))
	}
	
	if len(baseImagesToLayersMap) != 0 {
		t.Errorf("got %d base image maps, want 0", len(baseImagesToLayersMap))
	}
}

func TestBuildBaseImageDetails_DifferentBaseImages(t *testing.T) {
	t.Parallel()
	
	layerMetadata := []models.LayerMetadata{
		{DiffID: "layer1"},
		{DiffID: "layer2"},
		{DiffID: "layer3"},
	}
	
	baseImagesToLayersMap := [][]models.BaseImageDetails{
		{{Name: "alpine:3.14"}},
		{{Name: "ubuntu:20.04"}},
		{{Name: "ubuntu:20.04"}},
	}
	
	if len(layerMetadata) != 3 {
		t.Errorf("got %d layers, want 3", len(layerMetadata))
	}
	
	if len(baseImagesToLayersMap) != 3 {
		t.Errorf("got %d base image maps, want 3", len(baseImagesToLayersMap))
	}
	
	if len(baseImagesToLayersMap[0]) != 1 || baseImagesToLayersMap[0][0].Name != "alpine:3.14" {
		t.Error("first base image map should contain alpine:3.14")
	}
	
	if len(baseImagesToLayersMap[1]) != 1 || baseImagesToLayersMap[1][0].Name != "ubuntu:20.04" {
		t.Error("second base image map should contain ubuntu:20.04")
	}
}
