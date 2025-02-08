from flask import Flask, request

from handlers.info import handle_info
from handlers.move import handle_move


def configure_app():
    server_header = "battlesnake/replit/starter-snake-python"
    app = Flask(__name__)

    @app.get("/")
    def on_info():
        return handle_info()

    @app.post("/start")
    def on_start():
        return "ok"

    @app.post("/end")
    def on_end():
        return "ok"

    @app.post("/move")
    def on_move():
        game_state = request.get_json()
        return handle_move(game_state)

    @app.after_request
    def identify_server(response):
        response.headers.set("server", server_header)
        return response

    return app
