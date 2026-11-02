# Scryer Java Engine

Independent Go wrapper and Kotlin/JVM Java fact producer for
[Scryer](https://github.com/scryer-labs/scryer), licensed under [Apache-2.0](LICENSE).
The engine returns JSON only; terminal, Markdown and HTML reports belong to Core.
It imports no Core internal Go packages and uses no Core report assets.

## Build and run

Prerequisites: Go 1.24+, Java 21 and Git. The included Gradle Wrapper builds Kotlin.

```sh
./scripts/build
./build/bin/scryer-java-engine describe
export SCRYER_ENGINE_JAVA="$PWD/build/bin/scryer-java-engine"
/path/to/scryer project scan /path/to/java-project --static --json
/path/to/scryer project analyze /path/to/java-project --before HEAD~1 --after HEAD
```

Gradle creates the internal JVM-only intermediate distribution at
`build/install/java-worker/`. Its jars are copied into the complete engine; the Go
wrapper starts Java directly and does not invoke the Gradle-generated worker launcher.
The installed engine has no dependency on this intermediate directory.

`build/bin/scryer-java-engine` points into the complete development distribution at
`build/install/scryer-java-engine/`. Keep its bin/ and lib/java/ directories together.
Alternatively put the distribution's bin/ on PATH. A standalone `go build` builds
only the wrapper; jars must be supplied via SCRYER_JAVA_CLASSPATH in that case.

`describe` writes protocolVersion 1, id java, name/version/source and mock false.
`execute` forwards one v1 JSON request from stdin to the JVM and its response back
to stdout. Diagnostics use stderr. Both scan and analyze preserve existing Java
behavior, including build-tool selection, remote metadata, isolated snapshots,
custom test commands, execution/coverage evidence and explicit unknowns.

## Runtime and packages

```sh
./scripts/package
./scripts/package --with-runtime
```

Packages contain the Go wrapper, JVM jars, README and Apache license. The second
form adds the build host's Java 21 jlink runtime and its legal notices. Runtime
selection is SCRYER_ANALYZER_JAVA_HOME, bundled runtime, then JAVA_HOME/PATH.
SCRYER_JAVA_HOME independently selects the target project's JVM. Source builds
require developer toolchains; a bundled package needs no separately installed Go
or Java for the engine itself. Target builds/tests may still need their own tools.
GitHub Release publishing and release/source installation remain future work.

Each built/extracted package now includes engine.json (schemaVersion, identity,
relative executable, OS and architecture). Install the directory through Core:

```sh
scryer engine install java --from /path/to/extracted/package
scryer engine list --json
scryer project scan /path/to/project --static --json
scryer engine uninstall java
```

Core copies the complete package into ~/.config/scryer by default (including on
macOS), with an optional SCRYER_HOME override; installed
operation no longer depends on the source checkout or a path override. Existing
installs require uninstall before reinstalling. Install/list/uninstall are persistent.
The wrapper's package-manifest operation generates metadata during scripts/build;
it is separate from describe and scan/analyze fact data. Remote/source installation
and GitHub Release publication are still deferred.

## Verification and boundaries

The engine CI and Orrery project plan run Go race tests/vet, Kotlin fact/JSON tests,
build and both package checks. Core and engine plans run before their workspace
integration plan. See [AGENTS.md](AGENTS.md) and [protocol](docs/PROTOCOL.md).
Kotlin contains no legacy command/rendering fixtures or MD/HTML report resources.

Current execution is Core → Go wrapper → JVM. The JVM inherits the wrapper's
process group so Core can cancel the complete operation on macOS. A future pair
experiment will explore a C ABI shared-library wrapper; it is not implemented.
