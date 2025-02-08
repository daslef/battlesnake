from handlers import move


def test_avoid_hitting_the_north_wall():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 1, "y": 10},
            "body": [{"x": 1, "y": 10}, {"x": 1, "y": 9}, {"x": 1, "y": 8}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ("left", "right")


def test_avoid_hitting_the_west_wall():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 0, "y": 5},
            "body": [{"x": 0, "y": 5}, {"x": 1, "y": 5}, {"x": 2, "y": 5}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ("up", "down")


def test_avoid_hitting_the_northwest_corner():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 0, "y": 10},
            "body": [{"x": 0, "y": 10}, {"x": 1, "y": 10}, {"x": 2, "y": 10}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "down"


def test_avoid_hitting_the_southwest_corner_southbound():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 0, "y": 0},
            "body": [{"x": 0, "y": 0}, {"x": 0, "y": 1}, {"x": 0, "y": 2}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "right"


def test_avoid_hitting_the_south_wall():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 5, "y": 0},
            "body": [{"x": 5, "y": 0}, {"x": 5, "y": 1}, {"x": 5, "y": 2}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ("left", "right")


def test_avoid_hitting_the_southeast_corner_eastbound():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 10, "y": 0},
            "body": [{"x": 10, "y": 0}, {"x": 9, "y": 0}, {"x": 8, "y": 0}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "up"


def test_survive_the_northeast_corner_northbound():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 10, "y": 10},
            "body": [{"x": 10, "y": 10}, {"x": 10, "y": 9}, {"x": 10, "y": 8}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "left"


def test_avoid_hitting_the_northeast_corner_eastbound():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 10, "y": 10},
            "body": [{"x": 10, "y": 10}, {"x": 9, "y": 10}, {"x": 8, "y": 10}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] == "down"


def test_avoid_hitting_the_east_wall_southbound():
    game = {
        "board": {"height": 11, "width": 11, "snakes": []},
        "you": {
            "head": {"x": 10, "y": 9},
            "body": [{"x": 10, "y": 9}, {"x": 10, "y": 10}, {"x": 9, "y": 10}],
            "length": 3,
        },
    }

    result = move.handle_move(game)
    assert result["move"] in ("left", "down")
