# Engine JSON protocol v1

Core loads the installed shared library and calls `execute_v1` with the JSON
envelope. The library locates its JVM worker relative to its own path. The ABI
header is `cmd/java-engine-abi/scryer_engine_v1.h`.

The library accepts one UTF-8 JSON request of up to 1 MiB, with a request
ID of up to 256 bytes. The library starts its JVM for supported requests.
It forwards JSON over the worker's stdin/stdout, with diagnostics on stderr.
Output-file, color, HTML and Markdown parameters never enter the engine.

Envelope fields are protocolVersion (1), requestId, stack (java), operation
(scan or analyze), and the corresponding scan/analyze request object. Responses
repeat the first four fields and include report plus error (null on success).
Structured errors use code JAVA_ANALYZER_ERROR and a message; report is null.
A failed target test can instead be a valid report with execution.status FAILED.
Protocol version and report schemaVersion (1) are independent.

Scan fields: repository (absolute path), static, remoteList, optional buildTool.
Analyze fields: repository, immutable before/after SHAs, beforeRoot/afterRoot,
optional comparisonRepository, skipTests, optional testCommand and buildTool.
Core creates/cleans snapshots; engine verifies SHAs and owns Java facts, build/test
execution, coverage collection and JSON projection. Engine must not remove
caller-owned snapshots. Logs/evidence are retained under SCRYER_CACHE_HOME or the
existing temporary Scryer cache, outside snapshot directories.

Full report shapes and compatibility validation remain in the
[Core contract documentation](https://github.com/scryer-labs/scryer/blob/main/docs/contracts/README.md).
All existing fields, nullable unknowns, evidence datasets and additive extensions
must survive transport. There is no presentation layer in this repository.

The library starts the JVM in a separate process group. Core requests
cancellation through `cancel_v1`, and the worker adapter stops that group.

## Local package metadata

`scripts/build` generates `engine.json` with the schema version, engine
identity, library path, platform, and capabilities. Core validates the
manifest and loads the ABI symbols before persisting the installation.

The package manifest declares Java as its language, scan and analyze as its
operations, and dependencies as its selectable scan fact. A nonempty
`scan.facts` array selects a report with a `facts` object. Java reports
`UNSUPPORTED_FACT` and `UNSUPPORTED_OPERATION` as JSON protocol errors.

The Go engine service owns the capability definitions used to generate
`engine.json`. The C ABI delegates request validation to that service before it invokes the JVM worker.
