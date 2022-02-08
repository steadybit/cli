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

## Releasing

```
npm run ci
npm version {major|minor|patch}
git push --tags origin main
```
