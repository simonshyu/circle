import requests

data = {
    "data": "test_config_data",
}

r = requests.post("http://localhost:8000/scm/1/repos/1/config", json=data)

print(r.content)
