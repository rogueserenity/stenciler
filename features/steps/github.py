import os

from behave.model import given, then
from behave.runner import Context


@given('I have a "{visibility}" repository on GitHub')
def step_impl(
    context: Context,
    visibility: str,
):
    if visibility == "public":
        context.repository = "rogueserenity/stenciler-tests"
    elif visibility == "private":
        context.repository = "rogueserenity/stenciler-tests-private"


@given('I have the "{protocol}" URL of the repository')
def step_impl(
    context: Context,
    protocol: str,
):
    if protocol == "HTTPS":
        context.repository_url = f"https://github.com/{context.repository}.git"
    elif protocol == "SSH":
        context.repository_url = f"git@github.com:{context.repository}.git"


@given("I have a valid GitHub Authentication Token")
def step_impl(
    context: Context,
):
    context.auth_token = os.environ["TEST_REPO_TOKEN"]


@then("I see the current directory initialized with the repo template data")
def step_impl(
    context: Context,
):
    assert os.path.exists(context.output_config_file)
    assert os.path.exists(
        os.path.join(context.output_dir.name, "foo", "bar", "baz.txt")
    )
