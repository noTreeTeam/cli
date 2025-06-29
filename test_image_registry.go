package main

import (
	"fmt"
	"strings"
	"os"
)

// Simulate the fixed GetRegistryImageUrl function
func GetRegistryImageUrl(imageName string) string {
	registry := "ghcr.io" // Simulate the ghcr.io registry
	if registry == "docker.io" {
		return imageName
	}
	// Configure mirror registry
	parts := strings.Split(imageName, "/")
	imageNameOnly := parts[len(parts)-1]
	
	// Only replace Postgres images with notreeteam registry, leave all others as upstream
	if registry == "ghcr.io" && strings.HasPrefix(imageNameOnly, "postgres:") {
		return registry + "/notreeteam/" + imageNameOnly
	}
	
	// For all other images when using ghcr.io, keep the original upstream source
	if registry == "ghcr.io" {
		return imageName
	}
	
	return registry + "/supabase/" + imageNameOnly
}

func main() {
	testCases := []struct {
		input    string
		expected string
	}{
		// Postgres images should be replaced with notreeteam when using short names
		{"postgres:latest", "ghcr.io/notreeteam/postgres:latest"},
		{"postgres:15.8.1.085", "ghcr.io/notreeteam/postgres:15.8.1.085"},
		{"postgres:14.1.0.89", "ghcr.io/notreeteam/postgres:14.1.0.89"},
		
		// Already correctly prefixed postgres images should be kept as-is
		{"ghcr.io/notreeteam/postgres:latest", "ghcr.io/notreeteam/postgres:latest"},
		{"ghcr.io/notreeteam/postgres:15.8.1.085", "ghcr.io/notreeteam/postgres:15.8.1.085"},
		
		// Non-Postgres images should NOT be replaced and should keep original upstream source
		{"supabase/logflare:1.14.2", "supabase/logflare:1.14.2"},
		{"supabase/studio:2025.06.16-sha-c4316c3", "supabase/studio:2025.06.16-sha-c4316c3"},
		{"supabase/gotrue:v2.176.1", "supabase/gotrue:v2.176.1"},
		{"supabase/realtime:v2.36.18", "supabase/realtime:v2.36.18"},
		{"kong:2.8.1", "kong:2.8.1"},
		{"library/kong:2.8.1", "library/kong:2.8.1"},
		{"logflare:1.14.2", "logflare:1.14.2"},
		{"studio:2025.06.16-sha-c4316c3", "studio:2025.06.16-sha-c4316c3"},
	}

	fmt.Println("Testing GetRegistryImageUrl function:")
	fmt.Println("=====================================")

	allPassed := true
	for _, tc := range testCases {
		result := GetRegistryImageUrl(tc.input)
		passed := result == tc.expected
		if !passed {
			allPassed = false
		}
		
		status := "✓ PASS"
		if !passed {
			status = "✗ FAIL"
		}
		
		fmt.Printf("%s Input: %s\n", status, tc.input)
		fmt.Printf("    Expected: %s\n", tc.expected)
		fmt.Printf("    Got:      %s\n", result)
		fmt.Println()
	}

	if allPassed {
		fmt.Println("🎉 All tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some tests failed!")
		os.Exit(1)
	}
}
