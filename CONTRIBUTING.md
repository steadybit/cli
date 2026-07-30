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

## Tests

Tests sit at three levels. Put a test at the lowest one that can hold it &mdash; the
levels get slower and harder to debug as you go down this list.

| Level     | Tool                               | Covers                                                 |
| --------- | ---------------------------------- | ------------------------------------------------------ |
| Unit      | vitest                             | A single function or class, no I/O                     |
| Command   | vitest + msw + `@inquirer/testing` | A command end to end in process, including its prompts |
| Container | `e2e/run.sh` + expect              | Only what needs a real process                         |

Prompts are driven through `@inquirer/testing`. Mock the prompt package with
`wrapPrompt` so the application's own call is intercepted, and use the helpers in
`src/mocks/prompts.ts` rather than writing to the screen directly:

```ts
vi.mock('@inquirer/input', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/input')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

await answerPrompt('Profile name:', 'my-profile');
```

The container tests are deliberately thin. They exist for the four things no in-process
test can reach &mdash; real exit codes, a real terminal, the spawn of a subcommand, and
the packaged artifact &mdash; and they assert exit status and a line of output, never
content. Anything checking structure belongs at the command level.

```sh
docker build -t steadybit/cli:under-test .
docker run --rm -v "$PWD/e2e:/e2e" --entrypoint sh steadybit/cli:under-test /e2e/run.sh
```

The scripts are mounted into the image rather than baked into a derived one, and CI runs
the image by id rather than by name. Both avoid the same mistake: a name is resolved
against a registry when it cannot be found locally, so the suite can end up exercising
the last release while reporting success.

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
