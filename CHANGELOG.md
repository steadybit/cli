# Changelog

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
