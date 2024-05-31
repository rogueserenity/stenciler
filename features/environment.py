import tempfile


def before_scenario(context, scenario):
    context.local_repo_dir = tempfile.TemporaryDirectory(prefix="stenciler-local-")


def after_scenario(context, scenario):
    context.local_repo_dir.cleanup()
