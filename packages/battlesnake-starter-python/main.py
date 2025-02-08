import logging
import bootstrap
import fire
from config import get_config


def serve(port):
    config = get_config()
    app = bootstrap.configure_app()

    logging.getLogger("werkzeug").setLevel(logging.ERROR)

    print(f"\Running on http://{config['host']}:{port or config['port']}")
    app.run(host=config["host"], port=port or config["port"], debug=True)


if __name__ == "__main__":
    fire.Fire(serve)
