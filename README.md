# Scryer Java Engine

An independent Go wrapper for Scryer's Java language engine, licensed under
[Apache-2.0](LICENSE). This initial repository is a **mock**, not a Java analyzer.
The existing Kotlin/JVM implementation has not moved here yet.

```sh
go build -trimpath -o build/bin/scryer-java-engine ./cmd/scryer-java-engine
./build/bin/scryer-java-engine describe
```

`describe` writes one JSON identity with protocolVersion 1, id `java`, version
`0.0.0-dev`, repository source and `mock: true`. Unknown operations fail with exit 2.
There is no scanning, JVM dependency, installation state or daemon at this stage.
The independent process boundary has no imports from Scryer's internal packages.

After both repositories have been built, use the mock explicitly:

```sh
scryer engine install java --executable ../scryer-java-engine/build/bin/scryer-java-engine
scryer engine list --executable ../scryer-java-engine/build/bin/scryer-java-engine
scryer engine uninstall java --executable ../scryer-java-engine/build/bin/scryer-java-engine
```

All three Scryer commands probe this engine's identity. They do not persist an
installation, copy/delete packages or create a global Scryer configuration. List
shows the explicitly supplied engine, not an installed inventory. Use `--json` for
the mock receipt. Source installation and Java extraction are later work.

Acceptance is run through the engine's Orrery project plan and the separate
Scryer + Java Engine workspace plan. The GitHub workflow contains the same Go
checks; local plan setup and contributor conventions are in [AGENTS.md](AGENTS.md).
