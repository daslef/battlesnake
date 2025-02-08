import os
import pathlib
import platform
import subprocess


def get_bot_url(bot_name):
    bots = {
        "dummy": "http://localhost:7001/dummy/",
        "simple": "http://localhost:7001/simple/",
    }
    return bots[bot_name]


def get_executable_path():
    def get_os_part(os_type):
        match os_type:
            case "Windows":
                return "windows"
            case "Linux":
                return "linux"
            case "Darwin":
                return "osx_arm"

    def get_executable_format(os_type):
        return ".exe" if os_type == "Windows" else ""

    current_os = platform.system()
    current_directory = pathlib.Path(__file__).parent
    binary_name = f"battlesnake{get_executable_format(current_os)}"

    return str(current_directory / get_os_part(current_os) / binary_name)


def main(browser=False, bots=[]):
    log_path = pathlib.Path(__file__).parent.parent / "tests" / "output.log"

    host = os.environ.get("HOST", "localhost")
    port = int(os.environ.get("PORT", "4000"))
    url = f"http://{host}:{port}"

    cli_path = get_executable_path()
    cli_params = ["play", "--name", "Snake", "--url", url]

    for bot in bots:
        cli_params.extend(["--name", bot, "--url", get_bot_url(bot)])

    if browser:
        cli_params.append("--browser")
    else:
        cli_params.extend(("-o", str(log_path)))

    return lambda: subprocess.run(
        [cli_path, *cli_params],
        shell=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
