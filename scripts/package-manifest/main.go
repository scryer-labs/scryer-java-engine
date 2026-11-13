// The build invokes this helper to write engine.json.
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/scryer-labs/scryer-java-engine/internal/engine"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(engine.PackageManifest()); err != nil {
		log.Fatal(err)
	}
}
