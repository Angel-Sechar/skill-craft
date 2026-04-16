# skill-craft

> Craft your development skills.

Senior-level development skills for AI coding agents. Stop getting incongruent, generic code from your agent.

---

## What it does

skill-craft installs curated `SKILL.md` files into your AI coding agents.
Each skill teaches your agent to code the way a senior software engineer would —
enforcing architecture rules, flagging anti-patterns, and applying best practices
on every task, every file, every suggestion. Permanently.

No more repeated instructions. No more incongruent output.

---

## 3 operating systems. 3 agents. 1 command.

|         | Claude Code | Codex | OpenCode |
| ------- | ----------- | ----- | -------- |
| Windows | ✓           | ✓     | ✓        |
| Linux   | ✓           | ✓     | ✓        |
| macOS   | ✓           | ✓     | ✓        |

---

## Install

**macOS / Linux**

```bash
brew install Angel-Sechar/tap/skill-craft
```

**Windows**

```bash
scoop bucket add angel https://github.com/Angel-Sechar/scoop-bucket
scoop install skill-craft
```

**Go**

```bash
go install github.com/Angel-Sechar/skill-craft/cmd/skill-craft@latest
```

---

## Usage

Just run it:

```bash
skill-craft
```

The interactive TUI guides you through 5 steps:

```
1. Framework      .NET / Spring Boot / Angular / SQL
2. Stack          specific version for your framework
3. Architecture   Clean / Hexagonal / Onion / Microservices
4. Driven         TDD / Event-Driven Architecture / DDD
5. Practices      SOLID / Dependency Injection / Clean Code / OOP
```

Skills are installed into every agent detected on your machine:

```
~/.claude/skills/            Claude Code
~/.config/opencode/skills/   OpenCode
~/.codex/skills/             Codex
```

---

## Skills catalog

### Framework

| Skill                   | Language   |
| ----------------------- | ---------- |
| .NET Framework 4.5      | C#         |
| .NET Core 8             | C#         |
| ASP.NET Core            | C#         |
| Spring Boot 2 — Java 17 | Java       |
| Spring Boot 3 — Java 21 | Java       |
| Angular 14              | TypeScript |
| Angular 17              | TypeScript |

### Database

| Skill              |       |
| ------------------ | ----- |
| MS SQL Server 2019 | T-SQL |

### Architecture

| Skill                  | Compatible with                 |
| ---------------------- | ------------------------------- |
| Clean Architecture     | Onion, Microservices            |
| Hexagonal Architecture | Onion, Microservices            |
| Onion Architecture     | Clean, Hexagonal, Microservices |
| Microservices          | Clean, Hexagonal, Onion         |

> Hexagonal and Clean Architecture are incompatible — the TUI enforces this automatically.

### Driven

| Skill                           |
| ------------------------------- |
| Test-Driven Development (TDD)   |
| Event-Driven Architecture (EDA) |
| Domain-Driven Design (DDD)      |

### Practices

| Skill                |
| -------------------- |
| SOLID Principles     |
| Dependency Injection |
| Clean Code           |
| OOP                  |

---

## What a skill does

Each `SKILL.md` file contains:

- **Real code examples** in the right language for your stack
- **Architecture rules** the agent enforces on every suggestion
- **Anti-patterns to reject** — the agent flags and warns before writing incongruent code
- **Senior-level decisions** baked in — the kind you learn from years of experience, not documentation

Your agent won't just know Spring Boot. It will know Spring Boot 3 with Java 21 the way a senior Java engineer knows it — virtual threads, Jakarta namespace, constructor injection, proper exception handling, no field injection, no `javax.*`.

---

## Why not just prompt your agent?

Prompting works for one conversation.
A skill works on every task, every file, every suggestion — forever.

| Prompting                                  | skill-craft                         |
| ------------------------------------------ | ----------------------------------- |
| You repeat yourself every session          | Set once, works always              |
| Agent forgets between conversations        | Skills persist across all sessions  |
| Incongruent output with no stack awareness | Opinionated for your exact stack    |
| You fix the same mistakes repeatedly       | Agent flags them before they happen |

---

## License

MIT — by [Angel Sechar](https://github.com/Angel-Sechar)
