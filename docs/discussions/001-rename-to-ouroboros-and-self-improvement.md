# Discussion 001: Rename to Ouroboros & Self-Improvement Architecture

**Date**: 2026-07-21

**Status**: Exploratory

---

## 1. The Rename: thoth-agent → ouroboros

Full identity shift: Go module path (`github.com/DaviMGDev/ouroboros`), all imports, binary name, repo on GitHub, docs.

- **Timing is ideal** — zero external consumers, everything behind `internal/`, empty `go.sum`. Cost only goes up from here.
- "Ouroboros" trades functional descriptiveness for poetic accuracy — the name captures recursive self-consumption, which directly maps to the planned features.
- GitHub handles renames gracefully (redirects); Go modules do not (old module path breaks), but moot right now.

## 2. Why "Ouroboros": The Self-Referential Features

Three planned features that embody the snake eating its tail:

| Feature | Layer | Description |
|---------|-------|-------------|
| **Self-MoA** | LLM/provider | A new `LLM` implementation wrapping multiple models with proposers + aggregator. The MoA configuration *is* a distinct model identity (e.g., `deepseek-v4-flash+moa-2p`). Transparent to the agent — `Chat()` returns a single clean `ChatResponse`. Aggregator breaks ties when proposers disagree on tool calls. |
| **Session analysis** | Agent/meta | A command (not an inline hook) that reads session files and produces improvement suggestions — hooks, skills, system prompt updates. Human-in-the-loop by default, with opt-in auto-approval. |
| **Self-suggestion** | Agent/hooks | The agent suggests its own improvements based on its behavior. Initially triggered by the analysis command; potentially later via `AfterAgent` hook for automatic post-session review. |

**These map to three nested ouroboroi**:
1. **LLM level** — MoA feeding outputs back into the aggregator (which may itself be one of the wrapped models in self-MoA mode).
2. **Agent level** — Tool results feeding back into the LLM loop.
3. **Meta level** — Session traces feeding into configuration changes that alter future agent behavior.

## 3. Hybrid Continuous Learning

Instead of batch-only retrospective analysis, a hybrid memory layer:

- **Automatic event capture during sessions**: the LLM recognizes user complaints and records them; the framework detects "tool failed → later succeeded" patterns and captures the argument delta. Dumb but reliable recording, no interpretation.
- **Per-tool "hints" slot**: each `Tool` gets a `Hints` field that accumulates usage tips (e.g., "quote arguments with spaces" for BashTool). Hints are injected into the tool schema at call time so the LLM sees them automatically.
- **Global lessons**: separate from tool hints — user preferences, project conventions — injected as system prompt material.
- **Retrospective analysis command**: processes accumulated raw events into structured hints/lessons that get installed. Smart layer is batched and reviewable; raw layer is real-time and automatic.
- **Storage open question**: do hints/lessons live in the session file, a separate learning file, or mutate the tool definitions themselves? Defines what "persistent learning across sessions" means.
- **Quality control tension**: fully automatic hint accumulation risks bad hints degrading future performance. The hybrid model gates installation behind analysis, keeping raw capture automatic but installation deliberate.

## 4. Architecture That Makes This Possible (Already Built)

- The `LLM` interface abstraction means MoA is just a new implementation — no changes to the agent loop.
- The hook middleware (6 lifecycle points, forward chaining for "before," reverse for "after") is the natural injection point for self-improvement triggers.
- `internal/` encapsulation means all of this can be built and refactored without external API commitments.
