# Steadybit CLI &nbsp;

**[Installation](#installation) |**
**[Authorization](#authorization) |**
**[Usage](#usage) |**
**[Changelog](CHANGELOG.md)**

---

The Steadybit CLI enables you to use the Steadybit platform features easier in an automated way and implement e.g. GitOps practices easily.
You can retrieve, create or adjust experiment designs as well as running them straight away.

## Prerequisites

- You need to have a Steadybit account. You can create a free account [via our website](https://www.steadybit.com/get-started/).
- at least Node.js 16 as local runtime

## Installation

Via npm
```sh
npm install -g steadybit
```

## Authorization

You need an API access token. You can grab one via our [platform](https://platform.steadybit.io/settings/api-tokens) through the `Settings -> API Access Tokens` page.

```bash
➜ steadybit config profile add
? Profile name: steadybit
? API access token: [hidden]
? Base URL of the Steadybit server: https://platform.steadybit.io
```

## Usage

Get an existing experiment yaml from Steadybit and write it to file:
```bash
steadybit experiment get -k ADM-1 -f experiment.yml
```

Only apply the experiment:
```bash
steadybit experiment apply -f experiment.yml
```

Apply and run the experiment in one step:
```bash
steadybit experiment run -f experiment.yml
```

Run existing experiment:
```bash
➜ steadybit experiment run -k ADM-1
```

Dump all experiments and executions from all teams:
```bash
➜ steadybit experiment dump -d ./dump
```

## Container Image

You can also use the cli via our container image:

```sh
docker run -e"STEADYBIT_TOKEN=****" steadybit/cli:latest experiment get -k ADM-1
```
