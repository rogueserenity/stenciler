import filecmp
import os
import subprocess

import yaml
from behave import given, when, then
from behave.runner import Context


@given("I have template with only raw copy files")
def step_impl(
    context: Context,
):
    os.makedirs(context.template_repo_dir.name + "/foo/bar/baz", exist_ok=True)
    with open(context.template_repo_dir.name + "/foo/bar/baz/file.txt", "w") as f:
        f.write("Hello, World!")
    os.makedirs(context.template_repo_dir.name + "/foo/a/b/c", exist_ok=True)
    with open(context.template_repo_dir.name + "/foo/a/b/c/file.txt", "w") as f:
        f.write("Hello, World!")
    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "raw-copy": [
                    "**/file.txt",
                ],
            },
        ],
    }
    with open(context.template_repo_dir.name + "/.stenciler.yaml", "w") as f:
        yaml.dump(yaml_data, f)


@when("I run stenciler init with the template in an empty directory")
def step_impl(
    context: Context,
):
    command = [
        "/workspaces/stenciler/stenciler",
        "init",
        "https://github.com/foo/bar.git",
        "-r",
        context.template_repo_dir.name,
    ]
    stenciler_init = subprocess.run(
        command,
        check=False,
        cwd=context.local_repo_dir.name,
        capture_output=True,
    )
    print(stenciler_init.stdout)
    print(stenciler_init.stderr)
    assert stenciler_init.returncode == 0


@then("I see the current directory initialized with the template data")
def step_impl(
    context: Context,
):
    dcmp = filecmp.dircmp(
        context.template_repo_dir.name + "/foo",
        context.local_repo_dir.name,
        ignore=[".stenciler.yaml"],
    )
    dcmp.report_full_closure()
    assert not dcmp.left_only
    assert not dcmp.right_only
    assert not dcmp.diff_files
