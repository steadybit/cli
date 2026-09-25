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
- at least Node.js 22.13 as local runtime

## Installation

Via npm

```sh
npm install -g steadybit
```

## Authorization

You need an API access token. You can grab one via our [platform](https://platform.steadybit.com/settings/api-tokens) through the `Settings -> API Access Tokens` page.

```bash
➜ steadybit config profile add
? Profile name: steadybit
? API access token: [hidden]
? Base URL of the Steadybit server: https://platform.steadybit.com
```

## Rate limiting

The platform meters API requests with a token bucket: a burst of 100 requests, refilling
by 25 every 15 seconds. The CLI paces itself to that allowance, so a command that walks a
large tenant &mdash; `experiment dump` above all &mdash; takes a while rather than being
rejected part way through. It warns up front when a dump covers more than 100 experiments.

Override the allowance if your deployment is configured differently:

| Variable                        | Default | Meaning                          |
| ------------------------------- | ------- | -------------------------------- |
| `STEADYBIT_RATE_LIMIT_BURST`    | `100`   | Requests allowed before pacing   |
| `STEADYBIT_RATE_LIMIT_REFILL`   | `25`    | Requests restored each interval  |
| `STEADYBIT_RATE_LIMIT_INTERVAL` | `15`    | Length of that interval, seconds |

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
steadybit experiment run -k ADM-1
```

Dump all experiments and executions from all teams:

```bash
steadybit experiment dump -d ./dump
```

Dump only certain teams, by team key:

```bash
steadybit experiment dump -d ./dump --team ADM WEBHOOK
```

Validate advice status

```bash
steadybit advice validate-status -e "Global" -q "k8s.cluster-name=dev-demo and k8s.namespace=steadybit-demo"
```

Every command shows examples with `--help`, e.g. `steadybit schedule create --help`.

### Experiments from templates

Find a template and the placeholders it asks for:

```bash
steadybit template list --search kubernetes
steadybit template get -i <template-id> --placeholders -f values.yml
```

Create an experiment from it, or update the one created before with the same external id:

```bash
steadybit experiment apply --template <template-id> --team ADM --environment Global \
  --external-id shop-latency --placeholders values.yml -p CLUSTER=prod
```

Create and run it in one step:

```bash
steadybit experiment run --template <template-id> --team ADM --placeholders values.yml
```

### Experiment runs

```bash
steadybit execution get -i 1234 -t json
steadybit execution cancel -i 1234
steadybit execution property set -i 1234 -k approvedBy --value "Jane Doe"
steadybit execution artifact list -i 1234
steadybit execution artifact download -i 1234 -d ./artifacts
```

### Experiment schedules

```bash
steadybit schedule create -k ADM-1 --cron "0 0 9 ? * MON-FRI" --timezone Europe/Berlin
steadybit schedule list --team ADM
steadybit schedule disable -i <schedule-id>
```

Schedules can be kept in Git like experiments. `apply` writes the id of a new schedule back
into its file, so that applying it again updates it:

```bash
steadybit schedule get -i <schedule-id> -f schedule.yml
steadybit schedule apply -f ./schedules -R
```

## Container Image

You can also use the cli via our container image:

```sh
docker run -e"STEADYBIT_TOKEN=****" steadybit/cli:latest experiment get -k ADM-1
```
