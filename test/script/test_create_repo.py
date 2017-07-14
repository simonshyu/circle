import requests

r = requests.post("http://localhost:8000/scm/1/repos/root/dcos-docker-demo")

print(r.content)
