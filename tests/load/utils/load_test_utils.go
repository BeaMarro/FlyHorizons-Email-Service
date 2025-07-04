package load_test_utils

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type LoadTestUtils struct {
}

// Ssave metrics to a JSON file
func (utils *LoadTestUtils) SaveToJSON(t *testing.T, metrics *vegeta.Metrics) {
	// JSON
	results := map[string]interface{}{
		"requests":     metrics.Requests,
		"success_rate": metrics.Success * 100,
		"latency_mean": metrics.Latencies.Mean.String(),
		"throughput":   metrics.Throughput,
		"errors":       metrics.Errors,
		"timestamp":    time.Now().Format(time.RFC3339), // Tracks when the test was ran
	}

	// Create or overwrite the JSON file
	file, err := os.Create("load_testing.json")
	if err != nil {
		t.Fatalf("Failed to create load_testing.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		t.Fatalf("Failed to write metrics to load_testing.json: %v", err)
	}
}

func (utils *LoadTestUtils) LogMetrics(t *testing.T, metrics *vegeta.Metrics) {
	t.Logf("Requests: %d", metrics.Requests)
	t.Logf("Success Rate: %.2f%%", metrics.Success*100)
	t.Logf("Latency (mean): %s", metrics.Latencies.Mean)
	t.Logf("Throughput: %.2f req/s", metrics.Throughput)
	t.Logf("Errors: %v", metrics.Errors)
}

func (utils *LoadTestUtils) EvaluateMetricsSuccess(t *testing.T, metrics *vegeta.Metrics) {
	if metrics.Success < 0.95 { // Fails if success rate is below 95%
		t.Errorf("Success rate too low: %.2f%% (expected >= 95%%)", metrics.Success*100)
	}
	if len(metrics.Errors) > 0 {
		t.Errorf("Encountered errors: %v", metrics.Errors)
	}
}
