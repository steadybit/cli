# Contributing

## Working Locally

### Initial Setup

```sh
git clone git@github.com:steadybit/cli.git
nvm use # Node.js 22.13 or later, per .nvmrc
npm ci
```

Run `npm run ci` before pushing. It type-checks, tests, lints and builds, and is the
same script CI runs.

### Local CLI Execution

```sh
# Define environment variables
export STEADYBIT_TOKEN="..."
export STEADYBIT_URL="http://localhost:8080"

# Build the CLI locally
npm run build

# Run some CLI commands
./cli experiment get -k ADM-1
```

### Local CLI installation

```sh
# Build the CLI locally
npm run build
# Package the CLI locally
npm pack
# Install the local package
npm i -g steadybit-*.tgz
```

## Releasing

Releases are published by CI, not from a workstation: pushing a `v*` tag triggers
[the release workflow](.github/workflows/release.yml), which publishes to npm via
trusted publishing and pushes the Docker image. Never run `npm publish` locally.

```sh
# 1. Head the CHANGELOG.md entry with the version being released
git commit -am 'chore: prepare release'

# 2. Bump package.json and create the matching v<version> tag
npm run ci
npm version {major|minor|patch}

# 3. Push the commit together with the tag
git push --follow-tags origin main
```

Use `major` for breaking changes, which includes raising the Node.js floor, since
that breaks installs for users on older runtimes.

## Contributor License Agreement (CLA)

In order to accept your pull request, we need you to submit a CLA. You only need to do this once. If you are submitting a pull request for the first time, just submit a Pull Request and our CLA Bot will give you instructions on how to sign the CLA before merging your Pull Request.

All contributors must sign an [Individual Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/individual-cla.md).

If contributing on behalf of your company, your company must sign a [Corporate Contributor License Agreement](https://github.com/steadybit/.github/blob/main/.github/cla/corporate-cla.md). If so, please contact us via office@steadybit.com.

If for any reason, your first contribution is in a PR created by other contributor, please just add a comment to the PR
with the following text to agree our CLA: "I have read the CLA Document and I hereby sign the CLA".
