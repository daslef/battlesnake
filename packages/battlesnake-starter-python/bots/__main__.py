import json
from flask import Flask, request
from snakes import get_snake


def main():
    app = Flask(__name__)

    @app.route("/<snake_name>/")
    def index(snake_name):
        color = "#ee0000" if snake_name == "dummy" else "#ef30be"
        return {
            "apiversion": "1",
            "author": "Battlesnake",
            "color": color,
            "head": "default",
            "tail": "default",
        }

    @app.route("/<snake_name>/start", methods=["POST"])
    def start():
        return "ok"

    @app.route("/<snake_name>/move", methods=["POST"])
    def move(snake_name):
        snake = get_snake(snake_name)()
        data = request.get_json()
        gamestate = snake.payload_to_game_state(data)
        move = snake.move(gamestate)
        if not move:
            return json.dumps({"move": "up", "taunt": "hello"})

        if type(move) is tuple:
            move, taunt = move
            return json.dumps({"move": move.direction(), "taunt": taunt})

        return json.dumps({"taunt": "hello", "move": move.direction()})

    @app.route("/<snake_name>/end", methods=["GET", "POST"])
    def end():
        return "ok"

    app.run(host="localhost", port=7001, debug=True)


if __name__ == "__main__":
    main()
