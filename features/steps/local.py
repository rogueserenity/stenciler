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
        f.write("{{Serenity}}")

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


@given("I have a local updated template with only raw copy files")
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

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "raw-copy": [
                    "**/file.txt",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


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


@given("I have a local updated template with a templated file and no prompts")
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

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
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
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given(
    "I have a local template with a templated file that prompts with no default value"
)
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

    context.prompts = {
        "What is the name of your ship?": "Serenity",
        "Who is the captain?": "Malcolm Reynolds",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                    },
                    {
                        "name": "captain",
                        "prompt": "Who is the captain?",
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


@given("I have a local updated template with existing values")
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

    context.prompts = {
        "What is the name of your ship?": "Serenity",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
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

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                        "value": "Serenity",
                    },
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given(
    "I have a local template with a templated file that prompts with a default value"
)
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

    context.prompts = {
        "What is the name of your ship?": "",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                        "default": "Serenity",
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


@given(
    "I have a local template with a templated file that prompts with no default value and a hook"  # pylint: disable=C0301
)
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
    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "validate_ship.py")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("echo 'Serenity'")
    os.chmod(hook, 0o755)

    context.prompts = {
        "What is the name of your ship?": "Alliance",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                        "validation-hook": "hooks/validate_ship.py",
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


@given(
    "I have a local template with a templated file that prompts with a default value and a hook"  # pylint: disable=C0301
)
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
    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "validate_ship.py")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("echo 'Serenity'")
    os.chmod(hook, 0o755)

    context.prompts = {
        "What is the name of your ship?": "",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                        "default": "Alliance",
                        "validation-hook": "hooks/validate_ship.py",
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


@given("I have a local template with a pre-init hook")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "pre-init-hooks": [
                    "hooks/create_ship.sh",
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


@given("I have a local template with a pre-update hook")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "pre-update-hooks": [
                    "hooks/create_ship.sh",
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

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "pre-update-hooks": [
                    "hooks/create_ship.sh",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given("I have a local template with multiple pre-init hooks")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)

    hook1 = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook1, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook1, 0o755)

    hook2 = os.path.join(hooks_dir, "create_foo.sh")
    with open(hook2, "w", encoding="utf-8") as f:
        f.write("touch bar/baz/foo.txt")
    os.chmod(hook2, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "pre-init-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
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

    with open(os.path.join(leaf_dir1, "foo.txt"), "w", encoding="utf-8") as f:
        pass


@given("I have a local template with multiple pre-update hooks")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)

    hook1 = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook1, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook1, 0o755)

    hook2 = os.path.join(hooks_dir, "create_foo.sh")
    with open(hook2, "w", encoding="utf-8") as f:
        f.write("touch bar/baz/foo.txt")
    os.chmod(hook2, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "pre-update-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
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

    with open(os.path.join(leaf_dir1, "foo.txt"), "w", encoding="utf-8") as f:
        pass

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "pre-update-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given("I have a local template with a post-init hook")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "post-init-hooks": [
                    "hooks/create_ship.sh",
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


@given("I have a local template with a post-update hook")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "post-update-hooks": [
                    "hooks/create_ship.sh",
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

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "post-update-hooks": [
                    "hooks/create_ship.sh",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given("I have a local template with multiple post-init hooks")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)

    hook1 = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook1, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook1, 0o755)

    hook2 = os.path.join(hooks_dir, "create_foo.sh")
    with open(hook2, "w", encoding="utf-8") as f:
        f.write("touch bar/baz/foo.txt")
    os.chmod(hook2, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "post-init-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
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

    with open(os.path.join(leaf_dir1, "foo.txt"), "w", encoding="utf-8") as f:
        pass


@given("I have a local template with multiple post-update hooks")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)

    hook1 = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook1, "w", encoding="utf-8") as f:
        f.write("mkdir -p bar/baz; echo 'RogueSerenity' > bar/baz/file.txt")
    os.chmod(hook1, 0o755)

    hook2 = os.path.join(hooks_dir, "create_foo.sh")
    with open(hook2, "w", encoding="utf-8") as f:
        f.write("touch bar/baz/foo.txt")
    os.chmod(hook2, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "post-update-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
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

    with open(os.path.join(leaf_dir1, "foo.txt"), "w", encoding="utf-8") as f:
        pass

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "post-update-hooks": [
                    "hooks/create_ship.sh",
                    "hooks/create_foo.sh",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given("I have a local updated template with a new param")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    leaf_dir1 = os.path.join(root, "bar", "baz")
    os.makedirs(leaf_dir1, exist_ok=True)
    with open(os.path.join(leaf_dir1, "file.txt"), "w", encoding="utf-8") as f:
        f.write("Rogue{{.ship}}-{{.captain}}\n")

    context.prompts = {
        "What is the name of your ship?": "Serenity",
        "Who is the captain?": "Malcolm Reynolds",
    }

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                    },
                    {
                        "name": "captain",
                        "prompt": "Who is the captain?",
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
        f.write("RogueSerenity-Malcolm Reynolds\n")

    yaml_data = {
        "templates": [
            {
                "repository": context.repository_url,
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "prompt": "What is the name of your ship?",
                        "value": "Serenity",
                    },
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)


@given("I have a local template with a pre-init hook that uses parameter variables")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    os.makedirs(root, exist_ok=True)

    hooks_dir = os.path.join(context.input_dir.name, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)
    hook = os.path.join(hooks_dir, "create_ship.sh")
    with open(hook, "w", encoding="utf-8") as f:
        f.write("echo ${STENCILER_CAPTAIN} > ${STENCILER_SHIP}")
    os.chmod(hook, 0o755)

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "params": [
                    {
                        "name": "ship",
                        "value": "Serenity",
                    },
                    {
                        "name": "captain",
                        "value": "Malcolm Reynolds",
                    },
                ],
                "pre-init-hooks": [
                    "hooks/create_ship.sh",
                ],
            },
        ],
    }
    with open(context.input_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)

    file = os.path.join(context.expected_dir.name, "Serenity")
    with open(file, "w", encoding="utf-8") as f:
        f.write("Malcolm Reynolds\n")


@given("I have a local updated template with init-only files")
def step_impl(
    context: Context,
):
    context.repository_url = "https://github.com/local/repo"
    context.template_root_dir = "foo"
    root = os.path.join(context.input_dir.name, context.template_root_dir)
    leaf_dir1 = os.path.join(root, "bar", "baz")
    os.makedirs(leaf_dir1, exist_ok=True)
    with open(os.path.join(leaf_dir1, "file.txt"), "w", encoding="utf-8") as f:
        f.write("Rogue")
    with open(os.path.join(leaf_dir1, "file2.txt"), "w", encoding="utf-8") as f:
        f.write("Rogue2")
    with open(os.path.join(leaf_dir1, "test.txt"), "w", encoding="utf-8") as f:
        f.write("Serenity")

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "init-only": [
                    "**/file.txt",
                    "**/file2.txt",
                ],
                "raw-copy": [
                    "**/file2.txt",
                ],
            },
        ],
    }
    with open(context.input_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)

    leaf_dir1 = os.path.join(context.expected_dir.name, "bar", "baz")
    os.makedirs(leaf_dir1, exist_ok=True)
    with open(os.path.join(leaf_dir1, "test.txt"), "w", encoding="utf-8") as f:
        f.write("Serenity")

    yaml_data = {
        "templates": [
            {
                "directory": "foo",
                "init-only": [
                    "**/file.txt",
                    "**/file2.txt",
                ],
            },
        ],
    }
    with open(context.output_config_file, "w", encoding="utf-8") as f:
        yaml.dump(yaml_data, f)
