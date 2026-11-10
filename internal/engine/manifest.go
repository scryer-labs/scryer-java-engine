package engine

import "runtime"

type Identity struct {
	ProtocolVersion int    `json:"protocolVersion"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Source          string `json:"source"`
	Mock            bool   `json:"mock"`
}

type Capabilities struct {
	ResolveDependencies bool     `json:"resolveDependencies,omitempty"`
	Languages           []string `json:"languages"`
	Operations          []string `json:"operations"`
	ScanFacts           []string `json:"scanFacts"`
}

type Manifest struct {
	Capabilities  Capabilities `json:"capabilities"`
	SchemaVersion int          `json:"schemaVersion"`
	Engine        Identity     `json:"engine"`
	Library       string       `json:"library,omitempty"`
	Executable    string       `json:"executable"`
	OS            string       `json:"os"`
	Arch          string       `json:"arch"`
}

func capabilities() Capabilities {
	return Capabilities{
		ResolveDependencies: true,
		Languages:           []string{"Java"},
		Operations:          []string{"scan", "analyze"},
		ScanFacts:           []string{"dependencies"},
	}
}

func Describe() Identity {
	return Identity{
		ProtocolVersion: 1,
		ID:              "java",
		Name:            "Scryer Java Engine",
		Version:         "0.0.0-dev",
		Source:          "github.com/scryer-labs/scryer-java-engine",
	}
}

func PackageManifest() Manifest {
	return Manifest{
		Capabilities:  capabilities(),
		SchemaVersion: 1,
		Engine:        Describe(),
		Library:       libraryPath(),
		Executable:    "bin/scryer-java-engine",
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
	}
}

func libraryPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "bin/libscryer-java-engine.dylib"
	case "linux":
		return "bin/libscryer-java-engine.so"
	default:
		return ""
	}
}
