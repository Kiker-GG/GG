import docker
import os
import pytest

client = docker.from_env()

image, _ = client.images.build(path=".", tag="my_http_server")

container = client.containers.run(image, ports={'8000': '8000'},detach=True)

os.system("pytest tests.py")

#container = client.containers.run(image, ports={'8000': ('127.0.0.1', 8000)})
#container = client.containers.run(image, ports={'8000': '8000'})