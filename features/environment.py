import tempfile


def before_scenario(context, scenario):  # pylint: disable=W0613
    context.local_repo_dir = tempfile.TemporaryDirectory(  # pylint: disable=R1732
        prefix="stenciler-local-"
    )


def after_scenario(context, scenario):  # pylint: disable=W0613
    context.local_repo_dir.cleanup()
