# BuilderOS

A personal command-line "operating system" for solo builders and indie hackers. BuilderOS helps you capture ideas the moment you have them, get a blunt AI gut-check on whether they're worth pursuing, timebox focused work sessions against a specific idea, keep a running dev journal, and see a report of where your time actually went.

It's a single-user, local-first tool: no server, no database — just a small Go CLI (plus an early desktop GUI) reading and writing flat JSON files in the current directory.

> Source: [github.com/b2netpro/builderos](https://github.com/b2netpro/builderos)

## What it does

BuilderOS ties together four workflows that solo builders usually track in separate apps:

| Workflow | What it's for |
|---|---|
| **Ideas** | Quickly capture an idea from the command line, optionally get instant AI feedback on it. |
| **Time tracking** | Start a stopwatch-style session tied to a specific idea, stop it with Enter, and have the duration logged. |
| **Journal** | Freeform dated notes with optional tags — a lightweight dev diary. |
| **Reporting** | Roll up logged time by idea, so you can see how much effort went into each one. |

## Architecture

```
cmd/builderos/   CLI entry point (main app)
cmd/gui/         Early desktop GUI (Fyne), idea list + add only
internal/ideas/  Idea capture, storage, listing, AI feedback plumbing
internal/journal/ Journal entry storage and listing (plain + markdown output)
internal/timebox/ Time-session start/stop and per-idea reporting
models/          Shared struct definitions (Idea, JournalEntry) — currently unused by the rest of the code
```

- **Persistence:** plain JSON files written to the current working directory — `ideas.json`, `journal.json`, `time_logs.json`. There's no database and no configurable data directory; the app must be run from wherever you want its state to live.
- **AI feedback:** `idea:add "..." --ai` calls the OpenAI Chat Completions API (`gpt-4`) via [`sashabaranov/go-openai`](https://github.com/sashabaranov/go-openai), using a system prompt that casts the model as "Chuck," a blunt AI dev mentor. Requires `OPENAI_API_KEY` in the environment.
- **GUI:** `cmd/gui` is a minimal [Fyne](https://fyne.io/) desktop app — currently just a list of ideas and an "Add Idea" box, sharing the same `internal/ideas` storage as the CLI. It's an early prototype, not a full port of the CLI's functionality.

## Tech stack

- Go 1.22
- [`github.com/sashabaranov/go-openai`](https://github.com/sashabaranov/go-openai) — OpenAI API client
- [`fyne.io/fyne/v2`](https://fyne.io/) — cross-platform GUI toolkit (GUI prototype only)
- No external database — data lives in local JSON files

## Building & running

```bash
git clone https://github.com/b2netpro/builderos.git
cd builderos
go build -o builderos ./cmd/builderos   # CLI
go build -o builderos-gui ./cmd/gui     # optional GUI prototype
```

Set `OPENAI_API_KEY` in your environment before using `--ai` idea feedback.

## Usage

```bash
# Capture an idea
builderos idea:add "Build a CLI for X"

# Capture an idea and get AI feedback on it
builderos idea:add "Build a CLI for X" --ai

# View feedback previously captured for an idea
builderos idea:feedback <id>

# List all captured ideas
builderos idea:list

# Start a timeboxed work session on idea 3 (press Enter to stop)
builderos time:start 3 "optional note"

# Add a journal entry, optionally tagged
builderos journal:add "Shipped the idea CLI today" --tags dev,builderos

# List journal entries (optionally filtered to the last N days, or as markdown)
builderos journal:list --last 7 --md

# Show total time logged per idea
builderos report
```

## Current state & rough edges

This is an early-stage personal tool, not a polished product. Notable gaps observed in the code:

- **Committed secret** — see the security notice above; this is the most urgent item.
- **Compiled binary checked into git** — a ~7.6 MB `builderos` executable is tracked in the repo root alongside the source. There's no `.gitignore`, so build artifacts and local data files (`ideas.json`, `journal.json`, `time_logs.json`) are all being committed rather than ignored.
- **No tests** — the repository has no `_test.go` files.
- **Local-only, single-user** — no multi-project support, no config file; the tool infers everything from the current directory.
- **`models` package is unused** — `Idea` and `JournalEntry` are defined a second time here but nothing in `cmd/` or `internal/` imports them; the real types live in `internal/ideas` and `internal/journal`.
- **GUI prototype is minimal** — it only supports viewing and adding ideas, with no time tracking, journaling, or AI feedback yet.
