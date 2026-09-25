# Steadybit CLI &nbsp;

**[Installation](#installation) |**
**[Authorization](#authorization) |**
**[Usage](#usage) |**
**[Changelog](CHANGELOG.md)**

---

The Steadybit CLI enables you to use the Steadybit platform features easier in an automated way and implement e.g. GitOps practices easily.
You can retrieve, create or adjust experiment designs as well as running them straight away.

## Prerequisites

You need a Steadybit account. You can create a free account [via our website](https://www.steadybit.com/get-started/).

## Installation

The CLI is a single binary for Linux, macOS and Windows, on amd64 and arm64.

Via npm, which installs the binary for your platform (any Node.js from 18 on):

```sh
npm install -g steadybit
```

Or download the archive for your platform from the
[releases](https://github.com/steadybit/cli/releases) (`checksums.txt` lists their SHA-256)
and put `steadybit` on your `PATH`:

```sh
curl -sL https://github.com/steadybit/cli/releases/latest/download/steadybit_linux_amd64.tar.gz | tar -xz steadybit
sudo mv steadybit /usr/local/bin/
```

Shell completion is available for bash, zsh, fish and PowerShell, see `steadybit completion --help`.

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

### Services

Keep services in Git like experiments. `apply` writes the id of a new service back into
its file, so that applying it again updates it:

```bash
steadybit service list --team ADM
steadybit service get -i <service-id> -f service.yml
steadybit service apply -f ./services -R
```

Gate a pipeline on the risk of a service:

```bash
steadybit service risk -i <service-id> --fail-above 50
```

Manage the experiments and variables of a service:

```bash
steadybit service experiment list -i <service-id>
steadybit service experiment provide -i <service-id> --template <template-id> -p REPLICAS=3
steadybit service experiment link -i <service-id> -k ADM-1 --category Redundancy
steadybit service variable set -i <service-id> endpoint=http://shop.internal region=eu
```

Service profiles work the same way:

```bash
steadybit service-profile list --origin custom
steadybit service-profile apply -f profile.yml
```

## Everyday use

```bash
steadybit experiment init                     # create an experiment from a template, answering its placeholders
steadybit execution watch -k ADM-1            # follow the latest run of an experiment live
steadybit experiment get -k ADM-1 --profile prod   # use another configured profile for one command
```

Shell completion (`steadybit completion --help`) completes experiment keys, team keys and
the ids of templates, schedules, services and service profiles from the platform.

## GitOps

Keep a team's experiments, schedules, services and custom service profiles in Git:

```bash
steadybit export --team ADM -d ./chaos    # write them as files
steadybit diff -d ./chaos                 # what differs from the platform; exits with 2 if anything does
steadybit apply -d ./chaos --dry-run      # what an apply would create or update
steadybit apply -d ./chaos                # profiles, services, experiments, then schedules
```

Each kind also has its own `diff`, and its `apply` a `--dry-run`, e.g.
`steadybit experiment diff -f ./experiments -R`. Fields the platform fills in with defaults
are not reported as differences.

## In CI

`experiment run --wait` fails the job when a run fails, and a few options make it fit
pipelines:

| Option                        | Does                                                                    |
| ----------------------------- | ----------------------------------------------------------------------- |
| `--report steadybit.xml`      | A JUnit report, one test case per step; `.json` for JSON                |
| `--timeout 30m`               | Cancels the run and fails when it has not ended in time                 |
| `--show-steps`                | Prints each step's state as it changes                                  |
| `--keep-running-on-interrupt` | Leaves the run going when the job is cancelled; by default it is stopped |

In GitHub Actions a summary of every run is added to the job summary.

### GitHub Actions

```yaml
- uses: steadybit/cli@v6
- run: steadybit experiment run -f ./experiments -R --yes --report steadybit.xml
  env:
    STEADYBIT_TOKEN: ${{ secrets.STEADYBIT_TOKEN }}
- uses: mikepenz/action-junit-report@v5
  if: always()
  with:
    report_paths: steadybit.xml
```

### GitLab CI

```yaml
chaos:
  image:
    name: steadybit/cli:6
    entrypoint: ['']
  script:
    - steadybit experiment run -f ./experiments -R --yes --report steadybit.xml
  artifacts:
    when: always
    reports:
      junit: steadybit.xml
```

Every listing prints the platform's items with `-t json` or `-t yaml`, and `--jq` filters
whatever JSON a command prints, without jq installed:

```bash
steadybit service list --team ADM --jq '.[] | "\(.id) \(.name)"'
```

## Container Image

You can also use the cli via our container image:

```sh
docker run -e"STEADYBIT_TOKEN=****" steadybit/cli:latest experiment get -k ADM-1
```
