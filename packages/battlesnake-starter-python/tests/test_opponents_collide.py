from handlers import move


def test_avoid_hitting_the_east_wall_and_another_snake_eastbound():
    game = {
        "board": {
            "height": 11,
            "width": 11,
            "snakes": [
                {
                    "id": "c8cb622e-7bde-4bdf-b570-2303c83777c0",
                    "name": "me",
                    "health": 92,
                    "body": [{"x": 10, "y": 5}, {"x": 9, "y": 5}, {"x": 8, "y": 5}],
                    "latency": 0,
                    "head": {"x": 10, "y": 5},
                    "length": 3,
                    "shout": "",
                    "squad": "",
                },
                {
                    "id": "ca66924c-d5fe-4035-9bc1-17f63ba39852",
                    "name": "opp",
                    "health": 96,
                    "body": [{"x": 10, "y": 4}, {"x": 9, "y": 4}, {"x": 8, "y": 4}],
                    "latency": 0,
                    "head": {"x": 5, "y": 9},
                    "length": 4,
                    "shout": "",
                    "squad": "",
                },
            ],
        },
        "you": {
            "head": {"x": 10, "y": 5},
            "body": [{"x": 10, "y": 5}, {"x": 9, "y": 5}, {"x": 8, "y": 5}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "up"


def test_avoid_head_to_head_one():
    game = {
        "board": {
            "height": 11,
            "width": 11,
            "snakes": [
                {
                    "id": "c8cb622e-7bde-4bdf-b570-2303c83777c0",
                    "name": "us",
                    "health": 99,
                    "length": 3,
                    "head": {"x": 5, "y": 5},
                    "body": [{"x": 5, "y": 5}, {"x": 5, "y": 6}, {"x": 4, "y": 6}],
                    "latency": 0,
                    "shout": "",
                    "squad": "",
                },
                {
                    "id": "opponent-123asd",
                    "name": "them",
                    "health": 99,
                    "length": 3,
                    "head": {"x": 6, "y": 4},
                    "body": [{"x": 6, "y": 4}, {"x": 7, "y": 4}, {"x": 7, "y": 3}],
                    "latency": 0,
                    "shout": "",
                    "squad": "",
                },
            ],
        },
        "you": {
            "head": {"x": 5, "y": 5},
            "body": [{"x": 5, "y": 5}, {"x": 5, "y": 6}, {"x": 4, "y": 6}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "left"


def test_avoid_head_to_head_two():
    you = {
        "head": {"x": 9, "y": 2},
        "body": [{"x": 9, "y": 2}, {"x": 9, "y": 1}, {"x": 9, "y": 0}],
        "length": 3,
    }

    game = {
        "board": {
            "height": 11,
            "width": 11,
            "snakes": [
                {
                    "id": "c8cb622e-7bde-4bdf-b570-2303c83777c0",
                    "name": "us",
                    "health": 99,
                    "length": 3,
                    "head": you["head"],
                    "body": you["body"],
                    "latency": 0,
                    "shout": "",
                    "squad": "",
                },
                {
                    "id": "opponent-123asd",
                    "name": "them",
                    "health": 99,
                    "length": 3,
                    "head": {"x": 9, "y": 4},
                    "body": [{"x": 9, "y": 4}, {"x": 9, "y": 5}, {"x": 9, "y": 6}],
                    "latency": 0,
                    "shout": "",
                    "squad": "",
                },
            ],
        },
        "you": you,
    }

    result = move.handle_move(game)
    assert result["move"] in ("left", "right")
