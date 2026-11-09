# ECC — Agent Instructions & Orchestration Guide

This project is powered by **ECC (Everything Claude Code / Engine for Coordinated Coding)** tailored for **Google Antigravity CLI (`agy`)**.

**ECC Version:** 2.2.3  
**Harness Target:** Antigravity Project (`.agents/`)

---

## 1. Core Principles

1. **Agent-First** — Delegate to specialized subagents for distinct domain tasks (`planner`, `architect`, `code-reviewer`, `security-reviewer`, `tdd-guide`, etc.).
2. **Test-Driven Development (TDD)** — Always write tests first (RED), then minimal code to pass (GREEN), then refactor (REFACTOR). Require $\ge 80\%$ test coverage.
3. **Security-First** — Never hardcode secrets, API keys, or tokens. Validate all user inputs and sanitize data at boundaries.
4. **Immutability** — Favor immutable data structures and pure functions; avoid modifying state in-place.
5. **Plan Before Execute** — Complex tasks or multi-file refactorings must be planned first with explicit risk analysis and phase breakdown.
6. **Clean Architecture & Small Files** — Keep functions focused (<50 lines) and files compact (typically 200–400 lines, max 800 lines).

---

## 2. Directory Layout & Customizations

Antigravity CLI discovers workspace customizations located under `.agents/`:

```text
.
├── .agents/
│   ├── agents/                   # 70+ Specialized Subagents (.md)
│   ├── rules/                    # Normative Rules and Scoped Globs (.md)
│   ├── skills/                   # Modular Runbooks and Procedures (SKILL.md)
│   ├── workflows/                # Multi-step Workflow Definitions
│   └── ecc-install-state.json    # Integrity Manifest for ECC Doctor/Repair
├── AGENTS.md                     # Root Orchestration Instructions (this file)
└── ... (Source Code)
```

---

## 3. Specialized Agents Catalog

All agents are accessible via `invoke_subagent` in fresh, isolated contexts or via the CLI (`agy agents`):

| Agent | Purpose | When to Delegate |
| :--- | :--- | :--- |
| [`planner`](.agents/agents/planner.md) | Feature decomposition & planning | Complex tasks, architecture redesigns |
| [`architect`](.agents/agents/architect.md) | Distributed system design & scalability | High-level system decisions, tradeoffs |
| [`tdd-guide`](.agents/agents/tdd-guide.md) | Test-driven development enforcement | Writing new features, fixing regressions |
| [`code-reviewer`](.agents/agents/code-reviewer.md) | Static analysis & code quality | Immediately after any code edit |
| [`security-reviewer`](.agents/agents/security-reviewer.md) | OWASP & vulnerability scanning | Auth logic, secrets handling, input validation |
| [`build-error-resolver`](.agents/agents/build-error-resolver.md) | Build & compiler error resolution | Fix compilation, TypeScript or bundle failures |
| [`refactor-cleaner`](.agents/agents/refactor-cleaner.md) | Dead code cleanup & consolidation | Removing unused modules, safe refactoring |
| [`performance-optimizer`](.agents/agents/performance-optimizer.md) | Runtime bottleneck & memory profiling | Slow operations, high resource usage |
| [`database-reviewer`](.agents/agents/database-reviewer.md) | PostgreSQL / SQL query optimization | Schema migrations, indexing, slow queries |
| [`e2e-runner`](.agents/agents/e2e-runner.md) | End-to-end integration tests | Testing critical user journeys |
| [`doc-updater`](.agents/agents/doc-updater.md) | Architecture docs & codemaps | Syncing documentation with implementation |
| [`python-reviewer`](.agents/agents/python-reviewer.md) | Python PEP 8 & typing review | Python code changes |
| [`typescript-reviewer`](.agents/agents/typescript-reviewer.md) | TypeScript & Node.js code review | TypeScript / JavaScript changes |
| [`go-reviewer`](.agents/agents/go-reviewer.md) | Idiomatic Go patterns & concurrency | Go code changes |
| [`rust-reviewer`](.agents/agents/rust-reviewer.md) | Rust ownership & safety review | Rust code changes |
| [`react-reviewer`](.agents/agents/react-reviewer.md) | React hooks & render performance | React components & JSX |
| [`vue-reviewer`](.agents/agents/vue-reviewer.md) | Vue Composition API & reactivity | Vue single-file components |

---

## 4. Proactive Agent Orchestration Map

The orchestrator should invoke specialists proactively without waiting for explicit prompt commands:

* **Complex Feature Request** $\rightarrow$ Delegate to [`planner`](.agents/agents/planner.md).
* **Writing / Modifying Code** $\rightarrow$ Delegate to [`tdd-guide`](.agents/agents/tdd-guide.md) then [`code-reviewer`](.agents/agents/code-reviewer.md).
* **Compilation / Type Error** $\rightarrow$ Delegate to [`build-error-resolver`](.agents/agents/build-error-resolver.md) or language-specific resolver.
* **Sensitive Code / Auth / Input Handling** $\rightarrow$ Delegate to [`security-reviewer`](.agents/agents/security-reviewer.md).
* **Architecture / Database Changes** $\rightarrow$ Delegate to [`architect`](.agents/agents/architect.md) or [`database-reviewer`](.agents/agents/database-reviewer.md).

---

## 5. Development Workflow (TDD)

1. **RED**: Write a failing behavioral test verifying the desired functionality.
2. **GREEN**: Implement minimal code necessary to satisfy the test.
3. **REFACTOR**: Clean up code, enforce readability and modularity, and verify test coverage $\ge 80\%$.
4. **REVIEW**: Proactively run `code-reviewer` and `security-reviewer`.

---

## 6. ECC Maintenance Commands

Execute via terminal inside this repository:

* **Health Check**: `npx ecc-universal doctor --target antigravity`
* **Restore Drifted Files**: `npx ecc-universal repair --target antigravity`
* **Memory Management**: `npx ecc-universal memory init`
* **Agent List**: `agy agents`
