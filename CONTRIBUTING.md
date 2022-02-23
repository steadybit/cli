# Contributing

## Working Locally

### Initial Setup

```
git clone git@github.com:steadybit/cli.git
```

### Local CLI Execution

```sh
# Define environment variables
export STEADYBIT_TOKEN="..."
export STEADYBIT_URL="http://localhost:8080"

# Build the CLI locally
npm run build

# Run some CLI commands
./dist/cli/steadybit.js service init
./dist/cli/steadybit.js service apply .steadybit.yml
```

### Local CLI installation
```sh
# Build the CLI locally
npm run build
# Package the CLI locally
npm pack
# Install the local package
npm i -g steadybit-1.0.0.tgz
```

## Releasing

```
# Prepare the release

# Edit CHANGELOG.md by renaming the `## Unreleased` header to `## v{next release number}`
git add CHANGELOG.md
git commit -m 'chore: prepare release'

# Make the release
npm run ci
npm version {major|minor|patch}
git push --tags origin main
```
