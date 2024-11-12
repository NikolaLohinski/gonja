package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func main() {
	// Define some custom configuration
	configuration := config.New()
	configuration.StrictUndefined = true

	// Build a custom environment
	environment := &exec.Environment{
		Context: exec.NewContext(map[string]interface{}{
			"example": "default context in the environment",
		}),
		Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
		// For example this disables all methods and filters
		// Methods:           builtins.Methods,
		// Filters:           builtins.Filters,
	}

	// Instantiate a filesystem loader from the current working directory
	loader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		panic(err)
	}

	// Define the root template as a string
	source := `{{ example is string }}`

	// Compute some unique identifier for the root template
	sha := sha256.New()
	if _, err := sha.Write([]byte(source)); err != nil {
		panic(err)
	}
	rootID := fmt.Sprintf("root-%s", string(sha256.New().Sum(nil)))

	// Shift the filesystem loader with regards to the root template
	shiftedLoader, err := loaders.NewShiftedLoader(rootID, bytes.NewReader([]byte(source)), loader)
	if err != nil {
		panic(err)
	}

	// Create the template object
	template, err := exec.NewTemplate(rootID, configuration, shiftedLoader, environment)
	if err != nil {
		panic(err)
	}

	// Render
	if err = template.Execute(os.Stdout, exec.EmptyContext()); err != nil { // Prints True to stdout
		panic(err)
	}
}
