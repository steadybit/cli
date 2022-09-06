# Steadybit CLI &nbsp;

**[Installation](#installation) |**
**[Authorization](#authorization) |**
**[Usage](#usage) |**
**[Changelog](CHANGELOG.md)**

---

The Steadybit CLI enables you to define resilience expectations and resilience policies that your infrastructure needs to comply with through configuration files that live next to your code. This allows you to implement gitops practices easily.

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
### Description
```bash
➜ steadybit policy-binding init
? Name: gateway
? Desired Resilience Level: A - challenges for loose coupling

In order to understand which infrastructure you are describing, we need
to be able to map it to data collected by the Steadybit agents.
Right now, we only support mapping to Kubernetes monitored deployments.

? Kubernetes cluster name: demo-dev
? Kubernetes namespace: steadybit-demo
? Kubernetes deployment: gateway

We need to ensure that the infrastructure is still operating as expected when
verifying compliance with the desired resilience policy. To do so require
a load-balanced HTTP endpoint that can be called during task
execution.

? URL: http://k8s.demo.steadybit.io/products
...
```

The CLI created a .yml file containing your policy binding in your current working directory. From now on you are able to fine tune your settings either via subsequent CLI calls or by editing the .yml once you understand the format.

```yml
id: 9cb2888d-3abf-4ac3-bf8b-40b7eb61d271
name: gateway
policies:
  - name: 'steadybit/definitions/policies/level-b'
    version: 0.1.5
  - name: 'steadybit/definitions/policies/level-e'
    version: 0.1.5
mapping:
  kubernetes:
    cluster: demo-dev
    namespace: steadybit-demo
    deployment: gateway
parameters:
  httpEndpoint: http://k8s.demo.steadybit.io/products
```

The default file name is `.steadybit.yml`. Most CLI commands assume that the input file is named this way. In case you prefer an alternative name, you can specify the alternative name via `-f {your name}.yml`.

### Sync the policy binding with the Steadybit platform
You can now upload the policy binding by executing
```bash
steadybit policy-binding apply
```

or even better, you can automatically sync the policy binding with every commit using our [Github Action](https://github.com/steadybit/define-service)

### See your Policy Binding in Action
You can now use the cli to open steadybit and watch your current state and work on the next steps.
```bash
steadybit policy-binding open
```

which opens the Steadybit platform in your standard web browser. From here, you can inspect the policy's current state, trigger challenges and much more.

As an alternative you can verify your current state directly via the CLI
```bash
steadybit policy-binding verify
```

## Container Image

You can also use the cli via our container image:

```sh
docker run -e"STEADYBIT_TOKEN=****" steadybit/cli:latest policy-binding verify
```
