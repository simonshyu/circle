import requests

data = {
    "host": "http://172.24.6.123",
    "login": "root",
    "password": "12345678",
    "scm_type": "gitlab"
}

r = requests.post("http://localhost:8000/scm_account", json=data)

print(r.content)
