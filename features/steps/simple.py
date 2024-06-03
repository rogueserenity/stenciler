import filecmp
import os

import yaml
from behave import given, then
from behave.runner import Context


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


@then("I see the current directory initialized with the template data")
def step_impl(
    context: Context,
):
    assert os.path.exists(context.output_config_file)
    dcmp = filecmp.dircmp(
        os.path.join(context.input_dir.name, context.template_root_dir),
        context.output_dir.name,
        ignore=[context.yaml_file_name],
    )
    dcmp.report_full_closure()
    assert not dcmp.left_only
    assert not dcmp.right_only
    assert not dcmp.diff_files
