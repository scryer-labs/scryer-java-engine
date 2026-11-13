# Scryer Java Engine contributor instructions

Keep the Go ABI adapter thin, the worker protocol separate from language
implementation, and all output facts explicit about mock/unsupported behavior.
Do not import Scryer's internal packages. Kotlin/JVM facts and runtime packaging belong here. Never add Core reporting resources or Kotlin MD/HTML/terminal rendering fixtures.
Use clear names/branches, English comments and tests named Test_should_xxx_when_xxx.
Commits use feat/fix/refactor/tests/docs/chore/revert followed by a concise summary.
Do not commit or push without authorization or rewrite existing history.

Use Orrery plans for acceptance, not separately executed test/build commands.
Workspace root: /Users/edwardnoaland/Desktop/experiment/scryer-lab.
Engine project plan: scryer-java-engine--daily-development.
Workspace linkage plan: scryer-lab--scryer-java-engine-integration.
Read the current manifest revision and invoke the plan through Orrery's API.
Acceptance order is Scryer Core plan, engine project plan, then workspace integration. Integration validates persistent installation/source independence, Core reporting, Java tests, snapshots and cancellation. Seven-project provider scans belong to modernization--scryer-orrery-integration in test_workspace, run after the linkage plan installs the current bundled engine for ordinary discovery. Keep the plan aligned with CI; use only the plan for test/build acceptance.

The installed engine uses a C ABI shared library in Core and launches a JVM worker in its own process group. Keep the header in cmd/java-engine-abi aligned with Core’s canonical header.
