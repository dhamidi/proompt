package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalEditor := os.Getenv("EDITOR")
	originalPicker := os.Getenv("PROOMPT_PICKER")
	
	// Clean up after test
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
		if originalPicker != "" {
			os.Setenv("PROOMPT_PICKER", originalPicker)
		} else {
			os.Unsetenv("PROOMPT_PICKER")
		}
	}()

	t.Run("default values when env vars not set", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		os.Unsetenv("PROOMPT_PICKER")

		config := Load()

		if config.Editor != "nano" {
			t.Errorf("Load() Editor = %q, want %q", config.Editor, "nano")
		}
		if config.Picker != "fzf" {
			t.Errorf("Load() Picker = %q, want %q", config.Picker, "fzf")
		}
	})

	t.Run("uses environment variables when set", func(t *testing.T) {
		expectedEditor := "vim"
		expectedPicker := "rofi"
		
		os.Setenv("EDITOR", expectedEditor)
		os.Setenv("PROOMPT_PICKER", expectedPicker)

		config := Load()

		if config.Editor != expectedEditor {
			t.Errorf("Load() Editor = %q, want %q", config.Editor, expectedEditor)
		}
		if config.Picker != expectedPicker {
			t.Errorf("Load() Picker = %q, want %q", config.Picker, expectedPicker)
		}
	})

	t.Run("empty environment variables use defaults", func(t *testing.T) {
		os.Setenv("EDITOR", "")
		os.Setenv("PROOMPT_PICKER", "")

		config := Load()

		if config.Editor != "nano" {
			t.Errorf("Load() Editor = %q, want %q", config.Editor, "nano")
		}
		if config.Picker != "fzf" {
			t.Errorf("Load() Picker = %q, want %q", config.Picker, "fzf")
		}
	})
}

func TestGetEnv(t *testing.T) {
	// Save original environment variable
	original := os.Getenv("TEST_VAR")
	
	// Clean up after test
	defer func() {
		if original != "" {
			os.Setenv("TEST_VAR", original)
		} else {
			os.Unsetenv("TEST_VAR")
		}
	}()

	t.Run("returns env var when set", func(t *testing.T) {
		expected := "test_value"
		os.Setenv("TEST_VAR", expected)

		result := getEnv("TEST_VAR", "default")
		if result != expected {
			t.Errorf("getEnv() = %q, want %q", result, expected)
		}
	})

	t.Run("returns default when env var not set", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		expected := "default_value"

		result := getEnv("TEST_VAR", expected)
		if result != expected {
			t.Errorf("getEnv() = %q, want %q", result, expected)
		}
	})

	t.Run("returns default when env var is empty", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		expected := "default_value"

		result := getEnv("TEST_VAR", expected)
		if result != expected {
			t.Errorf("getEnv() = %q, want %q", result, expected)
		}
	})
}
