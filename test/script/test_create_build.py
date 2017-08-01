import requests

r = requests.post("http://localhost:8000/scm/2/repo/5/build")

print(r.content)
