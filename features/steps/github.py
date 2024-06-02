import os
import subprocess

from behave import given, when, then
from behave.runner import Context


@given('I have a "{visibility}" repository on GitHub')
def step_impl(
    context: Context,
    visibility: str,
):
    context.visibility = visibility
    if visibility == "public":
        context.repository = "rogueserenity/stenciler-tests"
    elif visibility == "private":
        context.repository = "rogueserenity/stenciler-tests-private"


@given('I have the "{protocol}" URL of the repository')
def step_impl(
    context: Context,
    protocol: str,
):
    context.protocol = protocol
    if protocol == "HTTPS":
        context.repository_url = f"https://github.com/{context.repository}.git"
    elif protocol == "SSH":
        context.repository_url = f"git@github.com:{context.repository}.git"


@when("I run stenciler init with the repository URL in an empty directory")
def step_impl(
    context: Context,
):
    command = ["/workspaces/stenciler/stenciler", "init", context.repository_url]
    if context.visibility == "private" and context.protocol == "HTTPS":
        command.append("-t")
        command.append(os.environ["TEST_REPO_TOKEN"])
    stenciler_init = subprocess.run(
        command,
        check=False,
        cwd=context.local_repo_dir.name,
    )
    assert stenciler_init.returncode == 0


@then("I see the current directory initialized with the repo template data")
def step_impl(
    context: Context,
):
    assert os.path.exists(os.path.join(context.local_repo_dir.name, ".stenciler.yaml"))
    assert os.path.exists(os.path.join(context.local_repo_dir.name, "foo"))
    assert os.path.exists(os.path.join(context.local_repo_dir.name, "foo", "bar"))
    assert os.path.exists(
        os.path.join(context.local_repo_dir.name, "foo", "bar", "baz.txt")
    )