import requests

data = {
    "host": "http://192.168.1.141:10080",
    "login": "root",
    "password": "12345678",
    "scm_type": "gitlab",
    "private_token": "F9VGC5d1uT2ee42ZGgne"
}

data2 = {
    "host": "http://172.24.6.123",
    "login": "root",
    "password": "12345678",
    "scm_type": "gitlab",
    "private_token": "sWxpWA4Lxtz8KGixW2uy"
}

r = requests.post("http://localhost:8000/scm", json=data)

print(r.content)
