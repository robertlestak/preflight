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

Whereas each preflight driver accepts a single operational input, `preflight` accepts a list of inputs for each driver. This allows you to check multiple endpoints, environment variables, etc. in a single run.

```yaml
dns:
- endpoint: https://example.com
  new: new-example.us-east-1.elb.amazonaws.com
- endpoint: https://example.net
  new: new-example.us-east-1.elb.amazonaws.com
  timeout: 10s

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
- provider: gcp
  gcp:
    email: example@my-project.google.com

netpath:
- endpoint: my-db:3306
- endpoint: my-redis:6379
  timeout: 10s
```

See the documentation for each driver for more information about the configuration options.

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
$ preflight -config preflight.yaml
```

## Docker

`preflight` is also available as a Docker image. There are two image types available:

- slim: a minimal image with only the `preflight` binary
- full: an image with both the `preflight` binary and all of the preflight driver binaries

Each individual driver also implements its own Docker image, which can be used to run the driver as a standalone container.

### Building

#### Slim Image

```bash
make docker-slim
```

#### Full Image

```bash
make docker-full
```

### Running

#### kubectl debug

Note: this requires Kubernetes 1.25+

```bash
kubectl debug -n my-namespace -it --image=robertlestak/preflight:latest -c preflight --attach my-pod -- sh
```

```bash
kubectl -n my-namespace -c preflight cp preflight.yaml my-pod:/preflight.yaml
```

```bash
# now in your debug session, you can run
preflight -config /preflight.yaml
```

#### kubectl exec

If you don't have access to `kubectl debug`, you can also use `kubectl exec` to install and run `preflight` in a pod.

```bash
kubectl -n my-namespace exec -it my-pod -- bash -c "curl -sSL https://raw.githubusercontent.com/robertlestak/preflight/main/scripts/install.sh | bash"

# now you can run
kubectl -n my-namespace exec -it my-pod -- preflight -config preflight.yaml

# to install all the individual drivers, you can run
kubectl -n my-namespace exec -it my-pod -- bash -c "curl -sSL https://raw.githubusercontent.com/robertlestak/preflight/main/scripts/install_bins.sh | bash"


# now you can run each driver individually

# preflight-dns
kubectl -n my-namespace exec -it my-pod -- preflight-dns -endpoint https://example.com -new new-example.us-east-1.elb.amazonaws.com

# preflight-env
kubectl -n my-namespace exec -it my-pod -- preflight-env -e HELLO=world -e ANOTHER=value -e FOO

# preflight-id
kubectl -n my-namespace exec -it my-pod -- preflight-id -kube-service-account hello-world

# preflight-netpath
kubectl -n my-namespace exec -it my-pod -- preflight-netpath -endpoint my-db:3306
```