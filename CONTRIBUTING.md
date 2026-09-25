# Contributing

## Working Locally

```sh
git clone git@github.com:steadybit/cli.git
cd cli
go test ./...
```

Go 1.26 or later, as `go.mod` states. `./cli` runs the CLI from the working tree:

```sh
export STEADYBIT_TOKEN="..."
export STEADYBIT_URL="http://localhost:8080"
./cli experiment get -k ADM-1
```

Before pushing, run what CI runs:

```sh
gofmt -l api cmd internal     # prints nothing
go vet ./...
go test ./...
```

## Layout

| Path                     | Holds                                                                  |
| ------------------------ | ---------------------------------------------------------------------- |
| `cmd/steadybit`          | The entry point                                                        |
| `internal/cli`           | The commands and their flags, wired with cobra                         |
| `internal/<area>`        | What each command group does: `experiment`, `schedule`, `service`, ... |
| `internal/platform`      | The HTTP client: authentication, retries, rate limiting, errors        |
| `internal/jsyaml`        | YAML and JSON output byte-compatible with the former TypeScript CLI    |
| `api`                    | The platform client, generated from `openapi/platform-api.json`        |
| `internal/tools`         | The spec fetcher and the npm package builder, used by CI               |
| `npm/steadybit`          | The launcher the `steadybit` npm package runs                          |

## Tests

Tests sit at three levels. Put a test at the lowest one that can hold it; the levels get
slower and harder to debug as you go down this list.

| Level     | Tool                               | Covers                                                  |
| --------- | ---------------------------------- | ------------------------------------------------------- |
| Unit      | `go test`                          | A single function, no I/O                               |
| Command   | `go test` + `internal/platformtest` | A command end to end against a fake platform           |
| Container | `e2e/run.sh` + expect              | Only what needs a real process                          |

`internal/platformtest` starts an `httptest` server whose endpoints a test declares,
records every request, and points the configuration at it:

```go
p := platformtest.New(t)
p.Reply("GET /api/experiments/TST-1", platformtest.Reply{Body: design})
out, err := platformtest.Stdout(t, func() error { return experiment.Get(ctx, p.Client, opts) })
```

The container tests are deliberately thin. They exist for what no in-process test can
reach: real exit codes, a real terminal, and the packaged artifact. They assert exit status
and a line of output, never content.

```sh
docker build -t steadybit/cli:under-test .
docker run --rm -v "$PWD/e2e:/e2e" --entrypoint sh steadybit/cli:under-test /e2e/run.sh
```

### Output compatibility

Users keep the files `get` writes in Git, so the YAML and JSON the CLI writes must stay
byte for byte what the TypeScript CLI wrote with js-yaml and `JSON.stringify`.
`internal/jsyaml/testdata/cases.json` holds values and what js-yaml rendered for them,
and the tests compare against it. To add a case, add it to `generate.mjs` and regenerate:

```sh
cd internal/jsyaml/testdata && npm ci && node generate.mjs
```

## Platform API Client

The client in `api/platform.gen.go` is generated from the platform's OpenAPI spec,
committed as `openapi/platform-api.json`. Refresh both together:

```sh
go run ./internal/tools/spec fetch   # the live spec into openapi/platform-api.json
go generate ./api                    # the client from it
go build ./...                       # a breaking change fails here
```

CI fails when the generated client does not match the committed spec, and builds against
the live spec daily and before each release. When the API is versioned, use the latest
version only. Experiment designs and other files users keep pass through as documents,
not generated structs, so that fields the spec does not know yet are never dropped.

## Releasing

Releases are published by CI, not from a workstation: pushing a `v*` tag builds the
binaries with goreleaser, creates the GitHub release, publishes the npm packages (one per
platform, and `steadybit`, which installs the right one), and pushes the Docker image.

```sh
# 1. Head the CHANGELOG.md entry with the version being released
git commit -am 'chore: prepare release'

# 2. Tag and push
git tag v5.0.0
git push origin main v5.0.0
```

Use a major version for breaking changes to commands, flags, output or exit codes.

## Contributor License Agreement (CLA)

In order to accept your pull request, we need you to submit a CLA. You only need to do this once. If you are submitting a pull request for the first time, just submit a Pull Request and our CLA Bot will give you instructions on how to sign the CLA before merging your Pull Request.

All contributors must sign an [Individual Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/individual-cla.md).

If contributing on behalf of your company, your company must sign a [Corporate Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/corporate-cla.md). If so, please contact us via office@steadybit.com.

If for any reason, your first contribution is in a PR created by other contributor, please just add a comment to the PR
with the following text to agree our CLA: "I have read the CLA Document and I hereby sign the CLA".
