import filecmp
import os

from behave import then
from behave.runner import Context


def cat(file):
    with open(file, "r", encoding="utf-8") as f:
        print(f.read())


def verify_same(dcmp):
    assert len(dcmp.left_only) == 0
    assert len(dcmp.right_only) == 0
    if len(dcmp.diff_files) > 0:
        for file in dcmp.diff_files:
            print(f"Diff file: {file}")
            print("----------------")
            cat(os.path.join(dcmp.left, file))
            print("----------------")
            cat(os.path.join(dcmp.right, file))
            print("----------------")
    assert len(dcmp.diff_files) == 0
    for sub_dcmp in dcmp.subdirs.values():
        verify_same(sub_dcmp)


@then("I see the current directory initialized with the template data")
@then("I see the current directory updated with the template data")
def step_impl(
    context: Context,
):
    assert os.path.exists(context.output_config_file)
    dcmp = filecmp.dircmp(
        context.expected_dir.name,
        context.output_dir.name,
        ignore=[context.yaml_file_name],
    )
    dcmp.report_full_closure()
    verify_same(dcmp)
