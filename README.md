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

### Install All Driver Binaries

```bash
curl -sSL https://raw.githubusercontent.com/robertlestak/preflight/main/scripts/install_bins.sh | bash
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
- kube:
    serviceAccount: hello-world
- aws:
    arn: arn:aws:iam::123456789012:role/hello-world
- gcp:
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
  -equiv
        print equivalent command
  -log-level string
        log level (default "info")
  -remote string
        remote preflight server to run checks against
  -remote-token string
        token to use when running remote checks
  -server
        run in server mode
  -server-addr string
        server listen address (default ":8090")
  -server-token string
        token to use when running in server mode
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

#### server mode

`preflight` can also be started in server mode, which will then allow you to send requests to the server to run checks from your local `preflight` instance. This way you don't need to keep copying the configuration yaml file into the pod as you make changes.

```bash
# in one terminal, start the server
kubectl debug -n my-namespace -it --image=robertlestak/preflight:latest -c preflight --attach my-pod -- preflight -server

# in another terminal, forward the server port
kubectl -n my-namespace port-forward my-pod 8090:8090

# now on your local machine, run your checks
# against the remote server
preflight -remote http://localhost:8090 -config preflight.yaml
```

You can also specify the remote within the config file itself:

```yaml
remote: http://localhost:8090
dns:
- endpoint: https://example.com
  new: new-example.us-east-1.elb.amazonaws.com
```

```bash
preflight -config preflight.yaml
````
