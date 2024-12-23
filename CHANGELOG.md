# Changelog

## v4.2.5 (next release)
- `advice validate-status` - parameter `query` is now optional

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

- Upgrade Node.js version in CLI Docker image to avoid [high severity security issues](https://nodejs.org/en/blog/release/v18.12.1/).

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

- Changed `exec` and `execute` to `run` in all descriptions and commands to align the cli to the platforms wording. Added an alias `exec` for backward
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
- `steadybit service init` gives more meaningful error messages in case the Steadybit servers fails to respond successfully with teams.
- The service definition file name may now be supplied via the `-f` option. By default, all commands assume that the default service definition file name
  is `.steadybit.yml`.
- Print a success message when a service definition was applied.
- Change label from `pending` to `not checked` in `steadybit service verify` output.
- `steadybit def-repo set-version` can be used to modify version numbers in a task and policy definition repository.
- `steadybit def-repo check` can be used to check task and policy definition files in a repository for format, consistency and references.

## v2.1.1

- `steadybit service ls` renamed to `steadybit service show`

## v2.1.0

- `steadybit service ls` can be used to get a (filtered) list of policies and tasks.

## v2.0.1

- `steadybit service verify` now correctly terminates with status code `1` when at least one task is either pending or failing.

## v2.0.0

- **Breaking:** The service definition yaml format was changed to directly specify policies and tasks instead of desired resilience levels.
- **Breaking:** Health definition subsection was removed service definitions.
- **Breaking:** Adapted the `steadybit service verify` behavior so that it is compatible with the reworked API model.
- `steadybit service apply` now helps to resolve conflicts.
- `steadybit service init` now asks for team and environment information.
- `steadybit service apply` now attaches auto-generated tags to the service definition which provide information about the repository, commit and ref.

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
- Engines field defines unsupported Node.js version. Node.js >=14.17.0 is required because of our `AbortController` usage.
- Define correct `Content-Type` header when making API calls with request bodies.

## v0.2.0

- Add validation to the `contract init` command.
- Adapt generated yaml structure to new format.
- Improve documentation.

## v0.1.0

- Initial release for internal testing
