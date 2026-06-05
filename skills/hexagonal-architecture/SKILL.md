---
name: hexagonal-architecture
description: Use when authoring, reviewing, or designing Go or Python service code that follows ports-and-adapters / hexagonal architecture — features with use cases, driven ports (repositories, gateways), driving adapters (HTTP, CLI, queue consumers), or anything calling itself "ports", "adapters", "use case", or "application service". Also triggers on stage-fit questions ("is this over-engineered for a POC?"), promotion decisions (introducing a Clock port, in-memory fakes, contract tests, per-feature sqlc, UnitOfWork), and PRs that touch repository/adapter naming, dependency direction, or cross-feature transactions.
---

# Hexagonal Architecture

Port-based design for Go and Python services. The full reference is `reference.md` in this folder — it is the source of truth. This file is the discovery layer and the must-knows that prevent the most common violations.

## Always do this first

**Load `reference.md` before writing or reviewing code.** Do not work from this SKILL.md alone — it omits the running example, the layer-by-layer rules, the testing patterns, the operational concerns, and the running example in §10. The inline content below is triage only.

## The two rules that decide most calls

1. **Dependencies point inward.** Domain depends on stdlib only. Application services depend on domain + ports. Adapters depend on ports. The composition root (`main.go` / `bootstrap.py`) is the only place that knows every concrete type. An `import` violating this is a defect — read §4.
2. **Stage gates every prescription.** A pattern that is mandatory at scaling is ceremony at POC. Identify the project's stage (§3) before recommending any pattern. If the user has not stated stage, ask once.

## Stage gradient — quick pick

| Stage | When | What's mandatory beyond §6/§7 port/adapter rules |
| --- | --- | --- |
| **POC** | Solo/pair, no real users | Migrations + typed config + one integration test per feature |
| **Early production** | Real users ≤ thousands, team ≤ 4 | + structured Logger port, healthz/readyz, idempotency keys, slow-query tripwire, adapter integration tests per PR |
| **Scaling** | Multi-team or imminent extraction | + per-feature sqlc, Tracer/Metrics, latency budgets in CI, contract tests for ports with ≥2 adapters, documented cross-feature tx pattern |

Each stage **inherits** the prior. Anything not mandatory is deferred — only add it when a §3.4 promotion trigger fires.

## Promotion triggers — do not add these without the trigger

| Pattern | Trigger |
| --- | --- |
| Named driving-port interface (`UserRegistrar`) | A second driving adapter (CLI, queue, second HTTP route group) needs the use case |
| In-memory adapter fake | Integration test latency >100ms/case AND use case doesn't need Postgres-specific behavior |
| Cross-port contract test suite | A second production adapter implements the same port |
| `Clock` / `IDGenerator` port | First test depends on deterministic time/IDs |
| Per-feature sqlc output | ≥3 features OR a feature is within one quarter of extraction |
| Cross-feature transaction pattern | First use case must atomically write across two features |
| Read-model feature | ≥2 endpoints share the join OR join touches ≥3 features |
| Caching / read-replica / circuit-breaker / Tracer ports | Measured signal — never speculative |

Adding any of these without a trigger documented in the PR description is **rejection-on-sight**.

## Anti-patterns — reject on sight

- Port references a framework or generated type (`*http.Request`, `*sql.Rows`, `*pb.User`, sqlc-generated `db.User`).
- Adapter calls another adapter, or reads env vars, or builds its own DB pool.
- Use case references a concrete adapter type, retries external I/O directly, or threads `*sql.Tx` as an argument.
- Domain calls `time.Now()`, reads env, or imports a logger.
- Naming: `IUserRepository`, `UserRepositoryInterface`, `UserDAO`. The port IS the canonical name; adapters carry tech prefix (`PostgresUserRepository`).
- `CreateAndNotify` / `SaveAndSendEmail` port methods — compose at the application-service layer.
- Returning `pgx.ErrNoRows` / `404 Not Found` to the application service — map at the adapter boundary.
- Feature package imports another feature's internals (model, service, adapter, `db` sub-package) instead of its public ports.
- Mixing migration tools (`golang-migrate` + `alembic` + `atlas`) in one project.
- Reading `os.Getenv` anywhere outside `config/`.

Full table in §12. When in doubt, cite the section and the dimension violated.

## Naming — three lines

- **Port:** business name, no prefix/suffix (`UserRepository`, `EmailSender`).
- **Adapter:** technology prefix or suffix (`PostgresUserRepository`, `InMemoryUserRepository`, `RedisUserCache`).
- **Use case:** one verb-phrase per business operation (`RegisterUser`, `CancelOrder`) — never `UserService` with 12 methods.

## Section index — jump into `reference.md`

| Question | Section |
| --- | --- |
| Is this dependency legal? | §4 (rule), §4.1 (lint enforcement) |
| Do I need this pattern yet? | §3 stage gradient, §3.4 triggers |
| How do I structure a feature package? | §5 layers, §10 running example |
| Port shape / error contract / nil convention | §6 |
| Adapter shape / error mapping / retry placement | §7 |
| Cross-cutting (auth, tenancy, logger, clock) | §8 |
| Migrations, idempotency, config, slow queries | §9 |
| Cross-feature reads and writes | §10.4, Appendix A |
| Tests — domain, app service, contract, integration | §11 |
| Code-review checklist | §13 (per stage) |

## Asking before answering

Before reviewing or generating code, confirm: **stage**, **language** (Go or Python), and any **decisions not up for debate** (e.g., tx pattern already chosen). Recommendations made without these will be miscalibrated. The architecture-review skill enforces the same discipline.
