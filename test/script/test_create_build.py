import requests

r = requests.post("http://localhost:8000/scm/1/repos/1/build")

print(r.content)
