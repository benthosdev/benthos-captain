# -*- mode: Python -*-

# build benthos operator
def deploy_benthos_operator():
    docker_build(
      "ghcr.io/benthosdev/benthos-captain:latest",
      ".",
      ignore=[
        ".git",
        ".github", 
        "*.md",
        "LICENSE",
        ]
    )

    k8s_yaml(
        kustomize('./config/default')
    )

deploy_benthos_operator()
