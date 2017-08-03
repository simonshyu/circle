import requests

r = requests.post("http://localhost:8000/scm/2/repo/4/build")

print(r.content)
