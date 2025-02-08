import os


def get_config():
    host = os.environ.get("HOST", "localhost")
    port = int(os.environ.get("PORT", "4000"))
    url = f"http://{host}:{port}"

    return {"host": host, "port": port, "url": url}
