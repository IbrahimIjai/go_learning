# Go Learning

Personal Go learning repo, structured so each topic is independently runnable.

## How to Run

From the project root (`go_learning/`), run any topic with:

```bash
go run ./<folder>/
```

### Examples

```bash
go run ./001_hello/
go run ./002_pointers/
go run ./003_structs/
go run ./004_goroutines/basic/
go run ./004_goroutines/channels/
go run ./005-price-feed-aggregator/
```

## Adding a New Topic

```bash
mkdir 005_maps
# create 005_maps/main.go starting with:
#   package main
#   func main() { ... }
go run ./005_maps/
```

Each folder is its own Go package — one `main.go` with `package main` and `func main()`.

---

## Go vs JavaScript Mental Map

| JavaScript | Go | Purpose |
|---|---|---|
| `package.json` | `go.mod` | Project name, language version, dependencies |
| `package-lock.json` | `go.sum` | Locked dependency checksums |
| `node_modules/` | `~/go/pkg/mod/` | Downloaded packages cache (global, not per-project) |
| `npm install <pkg>` | `go get <pkg>` | Add a dependency |
| `npm install` | `go mod download` | Install deps from lock file |
| `npm run start` | `go run ./` | Run the project |
| `npm run build` | `go build ./` | Compile to a binary |
| `npx` | none needed | Go runs tools directly |
| `import x from 'y'` | `import "github.com/y"` | Import a package — Go uses full URL paths |
| `export default` | capitalized name | `func Foo()` is exported, `func foo()` is private |
| `async/await` | goroutines + channels | Concurrency model |
| `Promise` | channel / `sync.WaitGroup` | Waiting on async results |
| `try/catch` | multiple return values | `result, err := doThing()` — errors are values, not exceptions |
| `console.log()` | `fmt.Println()` | Print to stdout |
| `typeof` / `instanceof` | `reflect` package | Runtime type inspection |
| `interface` (TS) | `interface` | Structural typing — Go interfaces are implicit (no `implements`) |
| `class` | `struct` + methods | No classes in Go; attach methods to structs |
| `extends` / inheritance | composition (embedding) | Go favors embedding structs over inheritance |
| `.env` / `dotenv` | `os.Getenv()` | Reading environment variables |

### Key Differences to Internalize

1. **No local `node_modules`** — packages download to a global cache (`~/go/pkg/mod/`), not inside your project.

2. **Imports use URLs** — `import "github.com/gin-gonic/gin"` — the module path *is* the source URL.

3. **Errors are values, not exceptions** — instead of `try/catch`, every fallible function returns `(result, error)` and you check it explicitly.

4. **Interfaces are implicit** — a struct satisfies an interface just by having the right methods. No `implements` keyword needed.

5. **Goroutines, not async/await** — concurrency is built into the language via `go func()` and channels, not a library or syntax layer on top.

6. **Stdlib is huge** — HTTP servers, JSON, crypto, file I/O are all in the standard library. You reach for third-party packages far less often than in Node.
