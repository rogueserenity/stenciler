import os
import tempfile


def before_all(context):
    context.yaml_file_name = ".stenciler.yaml"

    # Initialize the context variables for stenciler command line flags
    context.repository_url = None
    context.auth_token = None
    context.input_dir = None
    context.template_root_dir = None


def before_scenario(context, scenario):
    context.output_dir = tempfile.TemporaryDirectory(  # pylint: disable=R1732
        prefix="stenciler-output-"
    )
    context.output_config_file = os.path.join(
        context.output_dir.name, context.yaml_file_name
    )

    if "remote" not in context.feature.tags and "remote" not in scenario.tags:
        context.input_dir = tempfile.TemporaryDirectory(  # pylint: disable=R1732
            prefix="stenciler-input-"
        )
        context.input_config_file = os.path.join(
            context.input_dir.name, context.yaml_file_name
        )


def after_scenario(context, _):
    context.output_dir.cleanup()
    if context.input_dir is not None:
        context.input_dir.cleanup()
