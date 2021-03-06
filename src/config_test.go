package main

import (
	"os"
	"runtime"
	"testing"
)

// Testing the configuration file is valid or not
func TestFunctionIsValid(t *testing.T) {
	t.Log("Testing configs")
	config := &Config{}
	if config.IsValid() {
		t.Errorf("Config without `shared_secret` and `dialect` can not be valid")
	}
	config = &Config{Dialect: "sth"}
	if config.IsValid() == true {
		t.Errorf("Config without `shared_secret` can not be valid")
	}
	config = &Config{SharedSecret: "ultrasafesecret"}
	if config.IsValid() == true {
		t.Errorf("Config without `dialect` can not be valid")
	}
}

// Testing signature requirement
func TestFunctionIsSignatureRequired(t *testing.T) {
	cases := []struct {
		Config         *Config
		ExpectedResult bool
	}{
		{&Config{}, true},
		{&Config{Signature: "required"}, true},
		{&Config{Signature: "not-existing-property"}, true},
		{&Config{Signature: "optional"}, false}}

	for _, c := range cases {
		if r := c.Config.IsSignatureRequired(); r != c.ExpectedResult {
			t.Errorf("Expected signature requirements was %s but it was %s instead", c.ExpectedResult, r)
		}
	}
}

// Testing worker size calculation
func TestFunctionGetMaxWorkerSize(t *testing.T) {
	t.Log("Testing worker size initialization")
	config := &Config{MaxWorkerSize: 0}
	if r := config.GetMaxWorkerSize(); r != runtime.NumCPU()+1 {
		t.Errorf("Expected worker size from default value was %d but it was %d instead", runtime.NumCPU()+1, r)
	}
	config = &Config{MaxWorkerSize: 433}
	if r := config.GetMaxWorkerSize(); r != 433 {
		t.Errorf("Expected worker size from configuration was %d but it was %d instead", 433, r)
	}
	os.Setenv("HAMUSTRO_MAX_WORKER_SIZE", "22")
	defer os.Unsetenv("HAMUSTRO_MAX_WORKER_SIZE")
	if r := config.GetMaxWorkerSize(); r != 22 {
		t.Errorf("Expected worker size from environment variable was %d but it was %d instead", 22, r)
	}
}

// Testing queue size calculation
func TestFunctionGetMaxQueueSize(t *testing.T) {
	t.Log("Testing queue size initialization")
	config := &Config{MaxQueueSize: 0}
	if r := config.GetMaxQueueSize(); r != (runtime.NumCPU()+1)*20 {
		t.Errorf("Expected queue size from default value was %d but it was %d instead", (runtime.NumCPU()+1)*20, r)
	}
	config = &Config{MaxQueueSize: 433}
	if r := config.GetMaxQueueSize(); r != 433 {
		t.Errorf("Expected queue size from configuration was %d but it was %d instead", 433, r)
	}
	os.Setenv("HAMUSTRO_MAX_QUEUE_SIZE", "22")
	defer os.Unsetenv("HAMUSTRO_MAX_QUEUE_SIZE")
	if r := config.GetMaxQueueSize(); r != 22 {
		t.Errorf("Expected queue size from environment variable was %d but it was %d instead", 22, r)
	}
}

// Testing port determination
func TestFunctionGetPort(t *testing.T) {
	t.Log("Testing port initialization")
	config := &Config{}
	if r := config.GetPort(); r != "8080" {
		t.Errorf("Expected port was %s but it was %s instead", "8080", r)
	}
	os.Setenv("HAMUSTRO_PORT", "8000")
	defer os.Unsetenv("HAMUSTRO_PORT")
	if r := config.GetPort(); r != "8000" {
		t.Errorf("Expected port was %s but it was %s instead", "8000", r)
	}
}

// Testing host determination
func TestFunctionGetHost(t *testing.T) {
	t.Log("Testing host initialization")
	config := &Config{}
	if r := config.GetHost(); r != "localhost" {
		t.Errorf("Expected host was %s but it was %s instead", "localhost", r)
	}
	os.Setenv("HAMUSTRO_HOST", "127.0.0.1")
	defer os.Unsetenv("HAMUSTRO_HOST")
	if r := config.GetHost(); r != "127.0.0.1" {
		t.Errorf("Expected host was %s but it was %s instead", "127.0.0.1", r)
	}
}

// Testing address determination
func TestFunctionGetAddress(t *testing.T) {
	t.Log("Testing address initialization")
	config := &Config{}
	if r := config.GetAddress(); r != ":8080" {
		t.Errorf("Expected address was %s but it was %s instead", ":8080", r)
	}

	os.Setenv("HAMUSTRO_PORT", "8000")
	defer os.Unsetenv("HAMUSTRO_PORT")
	if r := config.GetAddress(); r != ":8000" {
		t.Errorf("Expected address was %s but it was %s instead", ":8000", r)
	}

	os.Setenv("HAMUSTRO_HOST", "127.0.0.1")
	defer os.Unsetenv("HAMUSTRO_HOST")
	if r := config.GetAddress(); r != "127.0.0.1:8000" {
		t.Errorf("Expected address was %s but it was %s instead", "127.0.0.1:8000", r)
	}
}

// Testing the buffer size calculation for buffered storage
func TestFunctionGetBufferSize(t *testing.T) {
	t.Log("Testing the buffer size calculations")
	config := &Config{}
	if exp := ((runtime.NumCPU() + 1) * (runtime.NumCPU() + 1) * 20) * 10; config.GetBufferSize() != exp {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, config.GetBufferSize())
	}
	config = &Config{BufferSize: 100000}
	if exp := 100000; config.GetBufferSize() != exp {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, config.GetBufferSize())
	}
}

// Testing the spreading property
func TestFunctionIsSpreadBuffer(t *testing.T) {
	cases := []struct {
		Config         *Config
		ExpectedResult bool
	}{
		{&Config{}, false},
		{&Config{SpreadBufferSize: false}, false},
		{&Config{SpreadBufferSize: true}, true}}

	for _, c := range cases {
		if r := c.Config.IsSpreadBuffer(); r != c.ExpectedResult {
			t.Errorf("Expected spread buffer was %s but it was %s instead", c.ExpectedResult, r)
		}
	}
}

// Testing the retry attempt
func TestFunctionGetRetryAttempt(t *testing.T) {
	t.Log("Testing the retry attempt property")
	config := &Config{}
	if exp := 3; config.GetRetryAttempt() != exp {
		t.Errorf("Expected retry attempt was %d but it was %d instead", exp, config.GetRetryAttempt())
	}
	config = &Config{RetryAttempt: 8}
	if exp := 8; config.GetRetryAttempt() != exp {
		t.Errorf("Expected retry attempt was %d but it was %d instead", exp, config.GetRetryAttempt())
	}
}

// Testing the dialect determination
func TestFunctionDialectConfig(t *testing.T) {
	t.Log("Testing dialect selector when dialect is not exists")
	config := &Config{Dialect: "hohoho"}
	if _, err := config.DialectConfig(); err == nil {
		t.Errorf("Not existing dialect should raise an error")
	}

	t.Log("Testing dialect selector when existing dialect is lowercase")
	config = &Config{Dialect: "aqs"}
	if _, err := config.DialectConfig(); err != nil {
		t.Errorf("Existing dialect must be found when lowercase")
	}

	t.Log("Testing dialect selector when existing dialect is uppercase")
	config = &Config{Dialect: "AQS"}
	if _, err := config.DialectConfig(); err != nil {
		t.Errorf("Existing dialect must be found when uppercase")
	}
}

// Testing the truncate ip functionality setting
func TestFunctionIsMaskedIP(t *testing.T) {
	t.Log("Testing the truncate_ip when not defined")
	config := &Config{}
	if exp := false; config.IsMaskedIP() != exp {
		t.Errorf("Expected masked IP setting is %s but it was %s instead", exp, config.IsMaskedIP())
	}
	config = &Config{MaskedIP: false}
	if exp := false; config.IsMaskedIP() != exp {
		t.Errorf("Expected masked IP setting is %s but it was %s instead", exp, config.IsMaskedIP())
	}
	config = &Config{MaskedIP: true}
	if exp := true; config.IsMaskedIP() != exp {
		t.Errorf("Expected masked IP setting is %s but it was %s instead", exp, config.IsMaskedIP())
	}
}

// Test the maintance key is empty
func TestFunctionMaintanceKeyIsEmpty(t *testing.T) {
	t.Log("Testing the maintance key when not defined")
	config := &Config{}
	if exp := ""; config.MaintenanceKey != exp {
		t.Errorf("Expected maintenance key setting is %s but it was %s instead", exp, config.MaintenanceKey)
	}
}

// Testing the auto flush interval update function
func TestFunctionUpdateAutoFlushIntervalToSeconds(t *testing.T) {
	t.Log("Testing the auto flush interval property")
	config := &Config{}
	if exp := 0; config.AutoFlushInterval != exp {
		t.Errorf("Expected auto flush interval was %d but it was %d instead", exp, config.AutoFlushInterval)
	}
	config = &Config{AutoFlushInterval: 30}
	if exp := 30; config.AutoFlushInterval != exp {
		t.Errorf("Expected auto flush interval was %d but it was %d instead", exp, config.AutoFlushInterval)
	}
	config = &Config{AutoFlushInterval: 60}
	config.UpdateAutoFlushIntervalToSeconds()
	if exp := 3600; config.AutoFlushInterval != exp {
		t.Errorf("Expected auto flush interval was %d but it was %d instead", exp, config.AutoFlushInterval)
	}
}
