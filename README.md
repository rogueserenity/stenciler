# stenciler
 Software repository template manager

The development environment is in a [dev container](https://containers.dev). The container is preconfigured with aliases
for running common commands. This also allows all developers to share a common environment and allows the environment to
be used in [GitHub Codespaces](https://docs.github.com/en/codespaces).

## Building

```shell
docker compose build
```

## Linting

[golangci-lint](https://golangci-lint.run/) is used for linting. The default configuration is preferred, but if there
are specific exceptiions that are required, the configuration manual can be found
[here](https://golangci-lint.run/usage/configuration/). In order to run the linting, run:

```shell
gol
```

## Testing

### Unit Tests

Unit testing is done using the [Ginkgo](https://onsi.github.io/ginkgo/) Framework along with
[Gomega](https://onsi.github.io/gomega/) matchers.

#### Generating Mocks

[mockery](https://vektra.github.io/mockery/latest/) is used to generate mocks.
[This guide](https://vektra.github.io/mockery/latest/features/#packages-configuration) explains how to configure
`mockery` with it's `packages` configuration feature. Regenerate mocks by running:

```shell
gom
```

#### Running Unit Tests

Execute the unit tests by running:

```shell
gou
```

### Functional Tests

```shell
```
