# Contributing to stenciler

The development environment is in a [dev container](https://containers.dev). The container is preconfigured with aliases
for running common commands. This also allows all developers to share a common environment and allows the environment to
be used in [GitHub Codespaces](https://docs.github.com/en/codespaces).

## Commiting

All commit messages must meet the requirements of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
A [VS Code plugin](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits) is included
in the Dev Container to make it easier to complete this action.

## Building

```shell
gob
```

## Linting

[golangci-lint](https://golangci-lint.run/) is used for linting. In order to run the linting, run:

```shell
gol
```

## Testing

### Unit Tests

Unit testing is developed using [Testify](https://github.com/stretchr/testify) test suites and assertions.

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

Functional tests are developed in python using [behave](https://github.com/behave/behave). The alias below supports
passing arguments to behave like `-m` and `-w`.

```shell
ftest
```
