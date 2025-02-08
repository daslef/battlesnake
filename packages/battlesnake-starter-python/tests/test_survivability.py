import pathlib
from cli import run


def test_must_survive_at_least_20_turns():
    run.main(browser=False)()

    log_path = pathlib.Path(__file__).parent / "output.log"
    print(str(log_path))

    with open(str(log_path)) as f:
        moves = len(f.readlines()) - 3

    assert moves >= 20
