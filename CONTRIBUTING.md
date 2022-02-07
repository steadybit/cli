# Contributing

## Working Locally

### Initial Setup

```
git clone git@github.com:steadybit/cli.git
```

### Local CLI Execution

```
npm run build && STEADYBIT_TOKEN="42" ./dist/cli/steadybit.js contract init
```

## Releasing

```
npm run ci
npm version {major|minor|patch}
git push --tags origin main
npm publish
```
