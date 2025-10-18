package main

import (
	"testing"
)

// TestCreateSplashScreen tests that the splash screen can be created
func TestCreateSplashScreen(t *testing.T) {
	// Create dummy callbacks
	onStart := func() {}
	onQuit := func() {}
	
	// Create splash screen
	splash := createSplashScreen(onStart, onQuit)
	
	if splash == nil {
		t.Error("Expected splash screen to be created, got nil")
	}
}
