import requests

config = """
workspace:
  base: /go
  path: src/github.com/drone/envsubst

clone:
  git:
    image: plugins/git:0.5
    depth: 50

pipeline:
  test:
    image: golang:1.7
    commands:
      - ls /go/src/github.com/drone/envsubst
      - go version
  build:
    image: golang:1.7
    commands:
      - ls /go/src/github.com/drone/envsubst
      - cat ~/.netrc
"""

data = {
    "data": config,
}

r = requests.post("http://192.168.1.141:8000/scm/1/repo/1/config", json=data)

print(r.content)
