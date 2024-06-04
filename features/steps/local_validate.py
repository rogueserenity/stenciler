import filecmp
import os

from behave import then
from behave.runner import Context


@then("I see the current directory initialized with the template data")
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
    assert not dcmp.left_only
    assert not dcmp.right_only
    assert not dcmp.diff_files
