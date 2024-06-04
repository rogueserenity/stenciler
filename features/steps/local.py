import errno
import os
import shutil

import yaml
from behave import given
from behave.runner import Context


def copy(src, dest):
    try:
        shutil.copytree(src, dest, dirs_exist_ok=True)
    except OSError as e:
        # If the error was caused because the source wasn't a directory
        if e.errno == errno.ENOTDIR:
            shutil.copy(src, dest)
        else:
            assert False, e


@given("I have a local template with only raw copy files")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    leaf_dir1 = os.path.join(root, "bar", "baz")
    leaf_dir2 = os.path.join(root, "a", "b", "c")
    os.makedirs(leaf_dir1, exist_ok=True)
    os.makedirs(leaf_dir2, exist_ok=True)
    with open(os.path.join(leaf_dir1, "file.txt"), "w", encoding="utf-8") as f:
        f.write("Rogue")
    with open(os.path.join(leaf_dir2, "file.txt"), "w", encoding="utf-8") as f:
        f.write("Serenity")

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
    with open(context.input_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)

    copy(root, context.expected_dir.name)


@given("I have a local template with a templated file and no prompts")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    leaf_dir1 = os.path.join(root, "bar", "baz")
    os.makedirs(leaf_dir1, exist_ok=True)
    with open(os.path.join(leaf_dir1, "file.txt"), "w", encoding="utf-8") as f:
        f.write("Rogue{{.ship}}\n")

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "value": "Serenity",
                    },
                ],
            },
        ],
    }
    with open(context.input_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)

    leaf_dir1 = os.path.join(context.expected_dir.name, "bar", "baz")
    os.makedirs(leaf_dir1, exist_ok=True)
    with open(os.path.join(leaf_dir1, "file.txt"), "w", encoding="utf-8") as f:
        f.write("RogueSerenity\n")
