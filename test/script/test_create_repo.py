import requests

data = {
    "owner": "root",
    "name": "docker-demo",
    # "name": "dcos-docker-demo"
}

r = requests.post("http://localhost:8000/scm/1/repo", json=data)

print(r.content)
