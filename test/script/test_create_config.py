import requests

config = """
workspace:
  base: /go
  path: src/github.com/drone/envsubst

clone:
  git:
    image: plugins/git
    depth: 50

pipeline:
  build:
    image: alpine:3.2
    commands:
      - ls /go/src/github.com/drone/envsubst
      - cat ~/.netrc
"""

data = {
    "data": config,
}

r = requests.post("http://localhost:8000/scm/1/repo/1/config", json=data)

print(r.content)
