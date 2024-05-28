# Contributing to stenciler

The development environment is in a [dev container](https://containers.dev). The container is preconfigured with aliases
for running common commands. This also allows all developers to share a common environment and allows the environment to
be used in [GitHub Codespaces](https://docs.github.com/en/codespaces).

## Commiting

All commit messages should meet the requirements of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
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

## NOTES

Use [doublestar](https://github.com/bmatcuk/doublestar) for globbing

### init

```pseudo
if repo template dir not supplied, check out repo into temp dir
read template config
if more than one, prompt user to select
validate that every hook file exists for template, exit with error if any missing or not executable
run through prompts for the template, building local config
  run hooks on inputs as values are entered
write local config
run pre-init hooks in order
copy over all raw copy files
copy over all templated files, processing templates as we go
run post-init hooks in order
```

### update

```pseudo
if repo template dir not supplied, check out repo into temp dir
read local template config
validate that every hook file exists for template, exit with error if any missing or not executable
run pre-update hooks in order
copy over all raw copy files, exclude init-only
copy over all templated files, processing templates as we go, exclude init-only
run post-update hooks in order
```
