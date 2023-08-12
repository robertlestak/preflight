# preflight

`preflight` is a suite of tools to help test assertions about running infrastructure. It's designed to be used as part of a CI/CD pipeline to ensure that the infrastructure is in a known state before deploying new code or migrating to a new environment.

## Preflight Checks

- [preflight-dns](https://github.com/robertlestak/preflight-dns)
- [preflight-env](https://github.com/robertlestak/preflight-env)
- [preflight-id](https://github.com/robertlestak/preflight-id)
- [preflight-netpath](https://github.com/robertlestak/preflight-netpath)

## Install

NOTE: you will need `curl`, `bash`, and `jq` installed for the install script to work. It will attempt to install the binary in `/usr/local/bin` and will require `sudo` access. You can override the install directory by setting the `INSTALL_DIR` environment variable.

```bash
curl -sSL https://raw.githubusercontent.com/robertlestak/preflight/main/scripts/install.sh | bash
```

## Configuration

`preflight` is configured using a YAML or JSON file. The default location is a `preflight.yaml` file in the current working directory, but you can specify a different location using the `-config` flag.

```yaml
dns:
- endpoint: https://example.com
  new: new-example.us-east-1.elb.amazonaws.com
- endpoint: https://example.net
  new: new-example.us-east-1.elb.amazonaws.com

env:
  HELLO: world
  ANOTHER: value
  FOO: "" # an empty value will check for the existence of the variable

id:
- provider: kube
  kube:
    serviceAccount: hello-world
- provider: aws
  aws:
    arn: arn:aws:iam::123456789012:role/hello-world

netpath:
- endpoint: my-db:3306
- endpoint: my-redis:6379
```

See the documentation for each check for more information about the configuration options.

## Usage

```bash
$ preflight -h
Usage of preflight:
  -concurrency int
        number of concurrent checks to run (default 1)
  -config string
        path to config file (default "preflight.yaml")
  -log-level string
        log level (default "debug")
```

```bash
$ preflight
```
