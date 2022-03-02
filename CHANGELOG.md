# Changelog

## Unreleased
 - Node.js v14 support.
 - `steadybit service init` gives more meaningful error messages in case the Steadybit servers fails to respond successfully with teams.
 - The service definition file name may now be supplied via the `-f` option. By default, all commands assume that the default service definition file name is `.steadybit.yml`.
 - Print a success message when a service definition was applied.
 - Change label from `pending` to `not checked` in `steadybit service verify` output.

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
