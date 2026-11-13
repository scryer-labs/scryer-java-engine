# Scryer Java Engine

Independent Go C ABI library and Kotlin/JVM Java fact producer for
[Scryer](https://github.com/scryer-labs/scryer), licensed under [Apache-2.0](LICENSE).
The engine returns JSON only; terminal, Markdown and HTML reports belong to Core.
It imports no Core internal Go packages and uses no Core report assets.

## Build and run

Prerequisites: Go 1.24+, a C compiler for cgo, Java 21 and Git. The included Gradle Wrapper builds Kotlin.

```sh
./scripts/build
/path/to/scryer engine install java --from "$PWD/build/install/scryer-java-engine"
/path/to/scryer project scan /path/to/java-project --static --json
/path/to/scryer project analyze /path/to/java-project --before HEAD~1 --after HEAD
```

Gradle creates the internal JVM-only intermediate distribution at
`build/install/java-worker/`. Its jars are copied into the complete engine; the Go
wrapper starts Java directly and does not invoke the Gradle-generated worker launcher.
The installed engine has no dependency on this intermediate directory.

The development distribution is `build/install/scryer-java-engine/`. Core
loads the library in `bin/`, which starts the JVM worker with the jars in
`lib/java/`. Scan and analyze retain Java build-tool selection, remote metadata,
isolated snapshots, custom test commands, and execution and coverage evidence.

## Runtime and packages

```sh
./scripts/package
./scripts/package --with-runtime
```

Packages contain the Go C ABI library, JVM jars, README and Apache license. The second
form adds the build host's Java 21 jlink runtime and its legal notices. Runtime
selection is SCRYER_ANALYZER_JAVA_HOME, bundled runtime, then JAVA_HOME/PATH.
SCRYER_JAVA_HOME independently selects the target project's JVM. Source builds
require developer toolchains; a bundled package needs no separately installed Go
or Java for the engine itself. Target builds/tests may still need their own tools.
GitHub Release publishing and release/source installation remain future work.

Each built/extracted package now includes engine.json (schemaVersion, identity,
relative library, OS and architecture). Install the directory through Core:

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
The build generates `engine.json` from the engine's identity and capabilities.

## Verification and boundaries

The engine CI and Orrery project plan run Go race tests/vet, Kotlin fact/JSON tests,
build and both package checks. Core and engine plans run before their workspace
integration plan. See [AGENTS.md](AGENTS.md) and [protocol](docs/PROTOCOL.md).
Kotlin contains no legacy command/rendering fixtures or MD/HTML report resources.

Core calls the engine through the C ABI. The library starts the JVM in its
own process group and cancels that group when Core requests cancellation.
