import requests

r = requests.post("http://192.168.1.141:8000/scm/1/repo/1/build")

print(r.content)
