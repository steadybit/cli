# Changelog

## Unreleased
 - **Breaking:** The service definition yaml format was changed to directly specify policies and tasks instead of desired resilience levels.

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
