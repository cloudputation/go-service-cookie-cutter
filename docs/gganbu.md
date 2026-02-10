# gganbu

Context orchestrator CLI for Claude Code workflows.

## What It Is

A CLI tool that brings Kubernetes-style architecture to AI context management. It wraps Claude Code to orchestrate CLAUDE.md and CLAUDELET.md files over git.

## Architecture

```
CLAUDE.md (root)        = Control Plane
                          System-wide policies, codebase overview
                          Single source of authority

CLAUDELET.md (nodes)    = Distributed Context
                          Subsystem-specific knowledge
                          Live where the code lives

Code                    = Blackbox
                          gganbu never parses code
                          Only orchestrates context files

Git                     = Source of Truth
                          Discovery via git ls-files
                          Change detection via git diff
                          History as anchor points
```

## Problems It Solves

| Problem | Current State | With gganbu |
|---------|---------------|-------------|
| Context limits | AI hits window ceiling, loses coherence | Load only relevant nodes |
| No discovery | AI doesn't know what context exists | Manifest of all CLAUDELET.md files |
| Stale context | No way to detect outdated docs | Health checks tied to git history |
| Reactive loading | AI loads context after hitting walls | Proactive loading based on file touches |
| Manual scaling | CLAUDE.md grows into monolith | Distributed nodes scale with codebase |

## Core Concepts

**CLAUDELET.md** - Distributed context nodes. Each one provides focused, subsystem-specific information. Named to parallel CLAUDE.md (CLAUDE + let = small/local Claude context).

**Codebase policy enforcement**
Automatically enforces codebase policies in CLAUDE.md file and via system/user prompt to ensure a non probabilitic decision making experience. Steer Claude Code in decisive deterministic coding operation to prevent state deviation.

**Reconciliation Loop** - Like k8s controllers:
1. Discover all CLAUDELET.md files
2. Detect changes via git diff
3. Map file touches to relevant nodes
4. Health check for staleness
5. Pre-load context before AI starts work

**Checkpoints** - Git commits as anchor points. Meaningful commit messages become breadcrumbs for tracing development history.

## CLI Commands (Conceptual)

```bash
gganbu sync              # Reconcile context state
gganbu status            # Show node health and staleness
gganbu nodes             # List all CLAUDELET.md files
gganbu checkpoint        # Create anchor point with context
```

## Key Insight

This tool is FOR Claude Code, not for humans. Humans own the code, but the context orchestration serves the AI's ability to work effectively at scale.

Context should live where code lives - not in a central monolith that grows unbounded.

## Status

Concept stage. No existing tool fills this gap. Claude Code has primitive hierarchical CLAUDE.md support but no orchestration layer.
