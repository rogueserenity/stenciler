import os
import re
import subprocess
from queue import Queue, Empty
from threading import Thread

from behave.model import when
from behave.runner import Context


def enqueue_output(out, queue):
    while True:
        buf = ""
        c = out.read(1)
        while c != ":":
            buf += c
            c = out.read(1)
        buf = re.sub(r"\[.*\]", "", buf)
        queue.put(buf.rstrip())
        c = out.read(1)  # skip the trailing space


@when("I run stenciler init with the repository URL in an empty directory")
def step_impl(
    context: Context,
):
    stenciler = os.path.join(os.getcwd(), "stenciler")
    command = [stenciler, "init"]
    assert context.repository_url is not None, "context.repository_url must be provided"
    command.append(context.repository_url)

    if context.auth_token is not None:
        command.append("-t")
        command.append(context.auth_token)

    if context.input_dir is not None:
        command.append("-r")
        command.append(context.input_dir.name)

    init = subprocess.Popen(  # pylint: disable=R1732
        command,
        cwd=context.output_dir.name,
        text=True,
        bufsize=0,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )

    out_queue = Queue()
    out_thread = Thread(target=enqueue_output, args=(init.stdout, out_queue))
    out_thread.daemon = True
    out_thread.start()

    while True:
        try:
            line = out_queue.get_nowait()
            print(line)
            value = context.prompts[line]
            print(value)
            init.stdin.write(value + "\n")
        except Empty:
            if init.poll() is not None:
                break

    print("STDOUT: ", init.stdout.read())
    print("STDERR: ", init.stderr.read())
    assert init.returncode == 0


@when("I run stenciler update in the current directory")
def step_impl(
    context: Context,
):
    stenciler = os.path.join(os.getcwd(), "stenciler")
    command = [stenciler, "update"]

    if context.auth_token is not None:
        command.append("-t")
        command.append(context.auth_token)

    if context.input_dir is not None:
        command.append("-r")
        command.append(context.input_dir.name)

    update = subprocess.Popen(  # pylint: disable=R1732
        command,
        cwd=context.output_dir.name,
        text=True,
        bufsize=0,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )

    out_queue = Queue()
    out_thread = Thread(target=enqueue_output, args=(update.stdout, out_queue))
    out_thread.daemon = True
    out_thread.start()

    while True:
        try:
            line = out_queue.get_nowait()
            print(line)
            value = context.prompts[line]
            print(value)
            update.stdin.write(value + "\n")
        except Empty:
            if update.poll() is not None:
                break

    print("STDOUT: ", update.stdout.read())
    print("STDERR: ", update.stderr.read())
    assert update.returncode == 0
