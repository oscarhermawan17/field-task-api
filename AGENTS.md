# field-task-api — Go + Gin backend

Learning project backend for the `field-task` Expo app (sibling folder). The user (Oscar) is
learning Go here, so keep guidance educational and contrast with Node.js when helpful.

# Multi-device workflow — this is YOUR concern, not the user's

The user develops on **three machines** and switches often. Local setup does NOT sync between them.

| Device            | CPU                     | OS            |
|-------------------|-------------------------|---------------|
| Work laptop (ASUS)| x86_64                  | Windows 11    |
| MacBook Air       | Apple Silicon M1 (arm64)| macOS (Tahoe) |
| Home PC           | AMD Ryzen 5             | Windows 11    |

On a new thread, assume the local environment on THIS machine may be out of sync. Proactively
check setup before assuming a bug, so the user can stay focused on learning Go.

# How to run

From the `field-task-api/` folder:

```
go run ./cmd/api
```

Server listens on `PORT` (default 8080). Health check: `GET /health` -> `{"status":"ok"}`.

Build commands (from NOTED.md):
- `go run ./cmd/api` — dev run, compiles to cache (no .exe left behind).
- `go build ./cmd/api` — builds `api.exe` into the folder root (gitignored).
- `go build -o bin/field-task-api.exe ./cmd/api` — custom output path/name.

# Dependencies (Go vs Node)

Go auto-downloads dependencies on `go run` / `go build` — there is **no `node_modules` and no
mandatory install step** like `npm install`. If needed, run `go mod download` or `go mod tidy`.
So a "package not installed" theory is rarely the cause of a runtime error; if the program reached
its own log/validation line, compilation already succeeded and deps are present.

# Configuration / .env (IMPORTANT per-machine gap)

`.env` is **gitignored**, so it exists per-machine. If it is missing on this machine, copy it:

```
Copy-Item .env.example .env   # PowerShell
cp .env.example .env          # macOS
```

`AUTH_JWT_SECRET` must be **at least 32 characters** or the app exits at startup with
`AUTH_JWT_SECRET must be at least 32 characters`.

The app loads `.env` automatically via `github.com/joho/godotenv` (called at the top of `main()` in
`cmd/api/main.go`). `godotenv.Load()` is intentionally optional: it logs and continues if no `.env`
exists, because in real deployments values come from the OS environment instead.

# Database / seed — CURRENT STATUS (verified)

**As of now there is NO database in use. Data is in-memory (RAM) only.**

- `internal/auth/service.go` holds a hardcoded in-memory user store (see its own comment:
  "temporary in-memory user store ... Replace it with a repository when the database task begins").
- No SQL/pgx/postgres/gorm driver exists anywhere in the Go code yet.
- The single seed user is created in `NewService()` on every startup and disappears when the
  server stops. Password comes from `AUTH_SEED_PASSWORD` in `.env`, hashed with bcrypt.
  - Login: `worker@fieldtask.com` / `password123`
- `compose.yaml` provides Postgres (`docker compose up -d`) but the Go API does NOT connect to it
  yet — it is prepared for a later database task.

Implication for the multi-device concern: because the seed is hardcoded in the Go code (synced via
Git), the seeded data is IDENTICAL on every machine as long as the code matches. There is no local
DB state to diverge yet, so "seeded on Mac but not here" does NOT apply to the current stage.

Once the database task lands, update this section: from then on Postgres state is per-machine and
seeding must be checked per machine.

# Tool versions per machine

Go and Node versions can differ between the user's machines. Verify on THIS machine with
`go version` and `node --version` before assuming a version-related bug. Known values are tracked
in `field-task/AGENTS.md`. (This machine at last check: Go 1.26.5, Node 22.12.0 on Windows 11.)
