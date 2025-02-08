from handlers import move

board = {"height": 11, "width": 11, "snakes": []}


def test_avoid_backwards_self_collide_from_right():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 6},
            "body": [
                {"x": 6, "y": 6},
                {"x": 7, "y": 6},
            ],
            "length": 2,
        },
    }

    result = move.handle_move(game)
    assert result["move"] != "right"


def test_avoid_backwards_self_collide_from_left():
    game = {
        "board": board,
        "you": {
            "head": {"x": 7, "y": 6},
            "body": [
                {"x": 7, "y": 6},
                {"x": 6, "y": 6},
            ],
            "length": 2,
        },
    }

    result = move.handle_move(game)
    assert result["move"] != "left"


def test_avoid_backwards_self_collide_from_top():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 6},
            "body": [
                {"x": 6, "y": 6},
                {"x": 6, "y": 5},
            ],
            "length": 2,
        },
    }

    result = move.handle_move(game)
    assert result["move"] != "down"


def test_avoid_backwards_self_collide_from_bottom():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 5},
            "body": [
                {"x": 6, "y": 5},
                {"x": 6, "y": 6},
            ],
            "length": 2,
        },
    }

    result = move.handle_move(game)
    assert result["move"] != "up"


def test_avoid_self_collide_from_bottom():
    game = {
        "board": board,
        "you": {
            "head": {"x": 7, "y": 5},
            "body": [
                {"x": 7, "y": 5},
                {"x": 8, "y": 5},
                {"x": 8, "y": 6},
                {"x": 7, "y": 6},
                {"x": 6, "y": 6},
            ],
            "length": 5,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ["left", "down"]


def test_avoid_self_collide_from_top():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 6},
            "body": [
                {"x": 6, "y": 6},
                {"x": 7, "y": 6},
                {"x": 8, "y": 6},
                {"x": 8, "y": 5},
                {"x": 7, "y": 5},
                {"x": 6, "y": 5},
            ],
            "length": 6,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ["left", "up"]


def test_avoid_self_collide_from_left():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 6},
            "body": [
                {"x": 6, "y": 6},
                {"x": 6, "y": 7},
                {"x": 8, "y": 6},
                {"x": 7, "y": 7},
                {"x": 7, "y": 6},
                {"x": 7, "y": 5},
            ],
            "length": 6,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ["left", "down"]


def test_avoid_self_collide_from_right():
    game = {
        "board": board,
        "you": {
            "head": {"x": 6, "y": 6},
            "body": [
                {"x": 6, "y": 6},
                {"x": 6, "y": 7},
                {"x": 8, "y": 6},
                {"x": 5, "y": 7},
                {"x": 5, "y": 6},
                {"x": 5, "y": 5},
            ],
            "length": 6,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ["right", "down"]
