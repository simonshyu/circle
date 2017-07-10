import requests

data = {
    "owner": "root",
    "name": "test",
    "clone_url": "http://172.24.6.123/root/test.git",
    "default_branch": "master"
}

r = requests.post("http://localhost:8000/scm/1/repos", json=data)

print(r.content)
