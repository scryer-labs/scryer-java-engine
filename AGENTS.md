# Scryer Java Engine contributor instructions

Keep the Go executable thin, the CLI/process contract separate from language
implementation, and all output facts explicit about mock/unsupported behavior.
Do not import Scryer's internal packages. Kotlin/JVM extraction is later work.
Use clear names/branches, English comments and tests named Test_should_xxx_when_xxx.
Commits use feat/fix/refactor/tests/docs/chore/revert followed by a concise summary.
Do not commit or push without authorization or rewrite existing history.

Use Orrery plans for acceptance, not separately executed test/build commands.
Workspace root: /Users/edwardnoaland/Desktop/experiment/scryer-lab.
Engine project plan: scryer-java-engine--daily-development.
Workspace linkage plan: scryer-lab--scryer-java-engine-integration.
Read the current manifest revision and invoke the plan through Orrery's API.
For this mock slice, acceptance is Scryer's plan, the engine project plan, then
its workspace linkage plan. No Orrery/provider or test_workspace run is required.
