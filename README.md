# Scryer Java Engine

Independent Go C ABI library and Kotlin/JVM Java fact producer for
[Scryer](https://github.com/scryer-labs/scryer), licensed under [Apache-2.0](LICENSE).
The engine returns JSON only; terminal, Markdown and HTML reports belong to Core.
It imports no Core internal Go packages and uses no Core report assets.

## Install and scan

With a current Scryer CLI:

```sh
scryer engine install java
scryer engine upgrade java
scryer project scan /path/to/java-project --facts dependencies --json
```

Releases for macOS arm64, Windows amd64 and Linux amd64 include Java 21.
The project's own build and tests may require a different toolchain.

## Build from source

Prerequisites: Go 1.24+, a C compiler for cgo, Java 21 and Git. The included Gradle Wrapper builds Kotlin.

```sh
./scripts/build
/path/to/scryer engine install java --from "$PWD/build/install/scryer-java-engine"
/path/to/scryer project scan /path/to/java-project --facts dependencies --json
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
operation no longer depends on the source checkout or a path override. Use `scryer engine upgrade java` for official release updates.
Install/list/uninstall are persistent.
The build generates `engine.json` from the engine's identity and capabilities.

## Known gaps and deferred work

- Selectable scan facts currently expose `dependencies`. Architecture, files,
  call-graph and IPC do not have the modern filtered-fact contract in this engine.
- Java retains the legacy combined `analyze` operation. Separate `coverage` and
  `impact` operations, the new native multi-engine analysis flow and its
  `--resolve` lifecycle have not been ported to Java.
- Existing analysis consumes configured JUnit/JaCoCo evidence; it does not inject
  coverage instrumentation into an unconfigured project. Aggregate evidence does
  not establish individual-test attribution or assertion effectiveness.
- Missing classpaths/types, generated code, reflection/DI wiring and dynamic
  dispatch can leave source boundaries. Java workspace dependency matching and
  test-effectiveness/code-smell rules are deferred.
- Windows release checks cover installation and static scanning, not full native
  target builds/tests or process-tree cancellation acceptance.

## Verification and boundaries

The engine CI and Orrery project plan run Go race tests/vet, Kotlin fact/JSON tests,
build and both package checks. Core and engine plans run before their workspace
integration plan. See [AGENTS.md](AGENTS.md) and [protocol](docs/PROTOCOL.md).
Kotlin contains no legacy command/rendering fixtures or MD/HTML report resources.

Core calls the engine through the C ABI. The library starts the JVM as a child
process. On Unix it cancels the JVM process group when Core requests cancellation.
