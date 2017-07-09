import requests

data = {
    "host": "http://172.24.6.123",
    "login": "root",
    "password": "12345678",
    "scm_type": "gitlab",
    "private_token": "sWxpWA4Lxtz8KGixW2uy"
}

r = requests.post("http://localhost:8000/scm", json=data)

print(r.content)
