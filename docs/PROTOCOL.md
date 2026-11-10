# Engine JSON protocol v1

Core loads the installed shared library and calls `execute_v1` with the JSON
envelope. The library locates its JVM worker relative to its own path. The ABI
header is `cmd/java-engine-abi/scryer_engine_v1.h`.

For standalone use, run `scryer-java-engine execute`, write
one UTF-8 JSON request to stdin and close it. The engine starts its JVM for
supported requests. Stdout contains exactly one JSON response for a valid
envelope, framed by EOF. Progress and diagnostics use stderr. Both entrypoints
limit requests to 1 MiB and request IDs to 256 bytes. Output-file,
color, HTML and Markdown parameters never enter the engine.

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

`describe` is a separate operation without stdin. It returns protocolVersion 1,
id java, name Scryer Java Engine, version, source
 github.com/scryer-labs/scryer-java-engine, and mock false.

Full report shapes and compatibility validation remain in the
[Core contract documentation](https://github.com/scryer-labs/scryer/blob/main/docs/contracts/README.md).
All existing fields, nullable unknowns, evidence datasets and additive extensions
must survive transport. There is no presentation layer in this repository.

Core owns the outer process group and cancellation. The wrapper must start the
JVM inside that group rather than creating a separate group. On direct graceful
termination the wrapper also cancels its worker context. Windows process-tree
behavior is not claimed by the current macOS-first acceptance.

## Local package metadata

package-manifest returns schemaVersion 1, engine (the describe identity), executable
bin/scryer-java-engine, and the build GOOS/GOARCH as os/arch. scripts/build writes it
as engine.json in the distribution. Both package variants include it. Core validates
this document and the copied executable's describe identity before persisting an
installation. It describes a tool package, not Java source facts or a new version
of the scan/analyze report contract.

The package manifest declares Java as its language, scan and analyze as its
operations, and dependencies as its selectable scan fact. A nonempty
`scan.facts` array selects a report with a `facts` object. Java reports
`UNSUPPORTED_FACT` and `UNSUPPORTED_OPERATION` as JSON protocol errors.

The Go engine service owns the capability definitions used to generate
`engine.json`. Both the standalone CLI and the C ABI delegate request validation
to that service before it invokes the JVM worker.
