import requests

r = requests.post("http://localhost:8000/repo", json={"scm": "git", "clone_url": "clone", "default_branch": "master"})

print(r.content)
