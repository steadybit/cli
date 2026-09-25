# Changelog

## v5.0.0

- **The CLI is now a single binary written in Go.** It runs without Node.js, and is installed
  the same ways: `npm install -g steadybit`, which now installs the binary for your platform
  and works with any Node.js from 18 on, or the `steadybit/cli` container image, now 18 MB
  instead of 249 MB. It can also be downloaded directly from the GitHub releases. Commands, flags, messages, exit
  codes, profiles in `~/.steadybit` and the `STEADYBIT_*` variables are unchanged, and
  experiment, schedule and service files are written byte for byte as before.
- Writing a new experiment's key back into a YAML file no longer rewrites the file: the key
  is added at the top and comments, anchors and formatting are kept.
- Shell completion: `steadybit completion bash|zsh|fish|powershell`, which completes
  experiment keys, team keys and the ids of templates, schedules, services and profiles
  from the platform.
- **Interrupting `experiment run --wait` now cancels the run it started**, before exiting
  with 130 (Ctrl-C) or 143 (SIGTERM, as CI runners send when a job is cancelled), so an
  aborted pipeline no longer leaves an attack running. `--keep-running-on-interrupt` keeps
  the previous behaviour.
- `experiment run --wait` gains `--timeout` (cancel and fail a run that takes too long),
  `--report` (a JUnit report with a test case per step, or JSON) and `--show-steps`. In
  GitHub Actions a summary of every run is added to the job summary.
- `diff` for experiments, schedules, services and service profiles shows how files differ
  from the platform, and exits with 2 when they do; `apply --dry-run` reports what an
  apply would change.
- `export --team X -d dir`, `apply -d dir` and `diff -d dir` keep a team's experiments,
  schedules, services and custom service profiles in Git as one project.
- Every listing prints the platform's items with `-t json` or `-t yaml`, and `--jq` filters
  the JSON any command prints, without jq installed.
- `execution watch` follows a run live; `experiment init` creates an experiment from a
  template by asking for its placeholders.
- `--profile <name>` uses a configured profile for one command.
- A GitHub Action, `uses: steadybit/cli@v5`, installs the CLI on a runner.
- `experiment apply --template <id>` creates an experiment from an experiment template, or
  updates the one created before with the same `--external-id`. With `-k` it re-renders an
  existing experiment with new placeholder values. Placeholders are given with
  `-p KEY=VALUE`, from a file with `--placeholders`, or both.
- `experiment run --template <id>` creates and runs an experiment from a template in one
  step, with `--execution-variable` for values that apply to that run only. `--wait`,
  `--retries` and `--allowParallel` work as for any other run.
- `template list` and `template get` find templates and their placeholders;
  `template get --placeholders` writes a placeholders file to fill in.
- `execution` commands for experiment runs: `get`, `cancel`, `property set` and
  `property add` to annotate a run, and `artifact list` and `artifact download` for the
  files its actions attached.
- `schedule` commands to `list`, `get`, `create`, `update`, `enable`, `disable` and `delete`
  experiment schedules, and `apply` to manage them as files in Git, like experiments.
- `service` commands to `list`, `get`, `apply` and `delete` services, managing them as files
  in Git like experiments. `service risk` shows a service's risk, and with `--fail-above`
  fails a pipeline when it is too high. `service experiment` lists, links, unlinks and
  provides (from a profile template) a service's experiments, and `service variable` gets,
  merges or replaces its variables.
- `service-profile` commands to `list`, `get`, `apply` and `delete` service profiles.
- Every command now shows examples in its `--help`.
- Fixed the table printed by `advice validate-status` containing colour escape codes when
  piped. Tables now follow the same terminal check as the rest of the output.
- The CLI's API client is generated from the platform's OpenAPI spec. CI builds it against
  the live platform daily and before every release, so a breaking API change is caught
  before it reaches a pipeline.
- **Security:** `-v, --verbose` no longer prints the API access token, which CI jobs using
  it had in their logs.
- **Security:** requests are only ever sent to the configured platform. An absolute URL in a
  platform response, such as the `Location` header of a started run, now has its origin
  replaced by the configured one so that the access token cannot be sent elsewhere.
- Fixed experiment files being corrupted on `apply` and `get`: merge keys (`<<:`) parsed into
  a literal `"<<"` property and timestamps into strings. `js-yaml` 5 narrowed its default
  schema, and the tag set the CLI relies on is now requested explicitly.
- Fixed `advice validate-status` reporting every advice as not matching, whatever its real
  status. The platform reports `IMPLEMENTED` while the default for `--status` is written
  `Implemented`, and the two were compared exactly, so the check failed even when all
  advice was implemented. Case and the separator are now both ignored, so `IMPLEMENTED`,
  `Implemented` and `action needed` all work.
- Fixed `--wait` failing when the platform returns a relative `Location` header.
- Fixed `experiment get -t json` printing to stdout in a format no JSON reader accepts,
  with unquoted keys and single quotes, so piping it into `jq` failed. Writing to a file
  was unaffected.
- Fixed `--retries` and `--retryInterval` being misparsed when given more than once, which
  could turn `--retryInterval 30 --retryInterval 25` into a 65 second wait. Both now reject a
  non-numeric value instead of silently disabling retries.
- Fixed the "an experiment is already running, run in parallel?" recovery never triggering,
  because the response body it inspects had already been consumed.
- Fixed the rate-limit retry loop waiting out a full window after its final attempt, delaying
  a 429 it already had by up to five minutes.
- `experiment dump` now survives a large tenant. It bounds its concurrent requests instead
  of scaling them with the size of the tenant, paces them to the allowance the platform
  meters (a burst of 100, refilling by 25 every 15 seconds), and no longer discards the
  whole run when one experiment or execution cannot be fetched — those are reported and
  counted, and leave a non-zero exit status so an incomplete dump is not mistaken for a
  complete one. It reports how many experiments it is about to walk, and how long that
  will take, before starting.
- `experiment dump --team <keys...>` restricts the walk to the named teams. An unknown key
  is refused and the accessible ones listed, rather than quietly dumping less than asked.
- Requests that fail in transit, such as a flaky DNS lookup, are retried for methods that
  are safe to repeat. A single such failure used to end a command outright.
- `STEADYBIT_RATE_LIMIT_BURST`, `STEADYBIT_RATE_LIMIT_REFILL` and
  `STEADYBIT_RATE_LIMIT_INTERVAL` override the assumed rate limit for deployments
  configured differently. A value that is not a positive whole number is reported and
  ignored rather than silently changing how hard the CLI polls.
- Dependency updates

## v4.3.2

- Dependency updates

## v4.3.1

- Dependency updates
- Switch to npm trusted publishing

## v4.3.0

- Switch to npm trusted publishing

## v4.2.28

- Dependency updates

## v4.2.26

- Dependency updates

## v4.2.25

- Dependency updates

## v4.2.23

- Dependency updates

## v4.2.22

- Add `--retries` and `--retryInterval` options to `experiment run` command to retry on validation errors (e.g., missing targets)

## v4.2.21

- Dependency updates

## v4.2.20

- Dependency updates

## v4.2.19

- Dependency updates

## v4.2.18

- Dependency updates

## v4.2.17

- Dependency updates

## v4.2.16

- Dependency updates

## v4.2.15

- Dependency updates

## v4.2.14

- Dependency updates

## v4.2.13

- Dependency updates

## v4.2.12

- Dependency updates

## v4.2.11

- Dependency updates
- add --verbose option for request logging

## v4.2.10

- Dependency updates

## v4.2.9

- Dependency updates

## v4.2.8

- add option --allowParallel to run an experiment in parallel to other experiments. e.g.
  `steadybit experiment run -k ADM-1 --allowParallel --yes`

## v4.2.7

- add param forcePersist=true to run an experiment to create always an execution

## v4.2.6

- feat: allow parallel experiment executions

## v4.2.5

- `advice validate-status` - parameter `query` is now optional
- Execution of experiments will now show a link to the ui.

## v4.2.4

- Removed unnecessary dependency

## v4.2.3

- Update base image to node 22 / npm 11 (CVE-2024-21538)

## v4.2.2

- Dependency updates

## v4.2.1

- Dependency updates (CVE-2024-21538)

## v4.2.0

- Dependency updates
- Requires Node >= 18
- Renamed `experiment get-all` to `experiment dump` to prevent misuse.
- Added `advice validate-status`

## v4.1.1

- Dependency updates

## v4.1.0

- Add get-all command to fetch all experiments

## v4.0.3

- Experiment execution shows error/failure reason

## v4.0.2

- Fixed experiment execution via key / file

## v4.0.1

- Added support to configure profiles non-interactively

## v4.0.0

- **Breaking:** Policies feature has been removed from CLI and is not longer supported

## v3.0.4

- New `experiment [get|apply|delete]` commands
- New `experiment run -f <file>` to execute experiment from file

## v3.0.3

- Upgrade Node.js version in CLI Docker image to
  avoid [high severity security issues](https://nodejs.org/en/blog/release/v18.12.1/).

## v3.0.2

- Support new exit state `ERRORED` of experiments

## v3.0.1

- Use new team api format

## v3.0.0

- **Breaking:** This version contains breaking changes within the CLI interface. All commands previously found
  under `steadybit service` and `steadybit service-definition` can now be found under `steadybit policy-binding`.
  This was done to remove the problematic term `service` and `service-definition`.
- Update the Node.js version to Node.js 18 within the Docker image.

## v2.10.1

- Adding a new service will yield the services `id` so the API server can figure out a proper one.

## v2.10.0

- Adding `http-client-fault-tolerance` policy to the list of available policies when running `steadybit service init`.

## v2.9.1

- Changed `exec` and `execute` to `run` in all descriptions and commands to align the cli to the platforms wording.
  Added an alias `exec` for backward
  compatibility.

## v2.9.0

- Publish Docker image additionally under the major version, e.g., `steadybit/cli:2`

## v2.8.2

- Add missing `git` dependency to the CLI Docker image.

## v2.8.1

- Maintenance release to retag the Docker image.

## v2.8.0

- `steadybit experiment exec` can be used to execute single experiments.

## v2.7.1

- initial version of the CLI Docker image.

## v2.7.0

- Correct license file header.
- Execute tasks via `steadybit service exec`.

## v2.6.1

- Added a confirmation dialog if the user did not select any policies.

## v2.6.0

- Improve first-user experience when using the CLI.

## v2.5.0

- Fail hard when an incompatible Node.js version is used.
- Upgrade versions of policies suggested via `steadybit service init` to `0.2.2`.
- Provide additional help when a conflict occurs while running `steadybit service apply`.

## v2.4.0

- `steadybit config profile add` now provides guidance w.r.t. the API access token creation.

## v2.3.0

- `steadybit service init` now has an improved policy selection process.

## v2.2.0

- Node.js v14 support.
- `steadybit service init` gives more meaningful error messages in case the Steadybit servers fails to respond
  successfully with teams.
- The service definition file name may now be supplied via the `-f` option. By default, all commands assume that the
  default service definition file name
  is `.steadybit.yml`.
- Print a success message when a service definition was applied.
- Change label from `pending` to `not checked` in `steadybit service verify` output.
- `steadybit def-repo set-version` can be used to modify version numbers in a task and policy definition repository.
- `steadybit def-repo check` can be used to check task and policy definition files in a repository for format,
  consistency and references.

## v2.1.1

- `steadybit service ls` renamed to `steadybit service show`

## v2.1.0

- `steadybit service ls` can be used to get a (filtered) list of policies and tasks.

## v2.0.1

- `steadybit service verify` now correctly terminates with status code `1` when at least one task is either pending or
  failing.

## v2.0.0

- **Breaking:** The service definition yaml format was changed to directly specify policies and tasks instead of desired
  resilience levels.
- **Breaking:** Health definition subsection was removed service definitions.
- **Breaking:** Adapted the `steadybit service verify` behavior so that it is compatible with the reworked API model.
- `steadybit service apply` now helps to resolve conflicts.
- `steadybit service init` now asks for team and environment information.
- `steadybit service apply` now attaches auto-generated tags to the service definition which provide information about
  the repository, commit and ref.

## v1.0.0

- First public release.

## v0.6.1

- `steadybit service verify` exit code is wrong.

## v0.6.0

- List succeeded/failed/pending tasks in `steadybit service verify` output.

## v0.5.0

- Change terminology in `steadybit service init` command to comply with terminology found in the public documentation.
- Add a `steadybit service verify .steadybit.yml` command to compare the actual vs. desired resilience level.

## v0.4.0

- Correct supported Node.js version to `>=16`.
- Send a `User-Agent` header when making Steadybit API calls.
- Add a `steadybit service open .steadybit.yml` command.
- Add the ability to manage configuration profiles via `steadybit config profile`.

## v0.3.0

- Rename `contract` to `service definition`
- Support a `delete` command.
- Rename `establish` to `apply` to be closer to the wording of kubectl.
- Add `steadybit service` as an alias for `steadybit service-definition`.
- Engines field defines unsupported Node.js version. Node.js >=14.17.0 is required because of our `AbortController`
  usage.
- Define correct `Content-Type` header when making API calls with request bodies.

## v0.2.0

- Add validation to the `contract init` command.
- Adapt generated yaml structure to new format.
- Improve documentation.

## v0.1.0

- Initial release for internal testing
