import subprocess

from behave import when
from behave.runner import Context


@when("I run stenciler init with the repository URL in an empty directory")
def step_impl(
    context: Context,
):
    command = ["/workspaces/stenciler/stenciler", "init"]
    assert context.repository_url is not None, "context.repository_url must be provided"
    command.append(context.repository_url)

    if context.auth_token is not None:
        command.append("-t")
        command.append(context.auth_token)

    if context.input_dir is not None:
        command.append("-r")
        command.append(context.input_dir.name)

    stenciler_init = subprocess.run(
        command,
        check=False,
        cwd=context.output_dir.name,
        capture_output=True,
    )

    assert stenciler_init.returncode == 0
