import fire
import run


def cli(*args):
    process = run.main(browser=True, bots=args)()
    while text := process.stdout.decode("utf-8"):
        print(text, end="", flush=True)


if __name__ == "__main__":
    fire.Fire(cli)
