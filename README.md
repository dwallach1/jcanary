# jcanary (JSON Canaries)

Verify your API data contract with just a JSON config file. Easily add integration tests
to your CI pipelines with no custom code, just config files.

Take in a JSON config and perform verification of web server functionality.

**Motivation**: when building a new service and constantly making changes, we want to verify integration tests; however, these can be cumbersome to write and maintain. `jcanary` offers a solution to abstract your tests to JSON and have this define how you manage your integration tests.

 s
uses https://github.com/Jeffail/gabs.


TODO:

- [ ] dockerize
https://github.com/peter-evans/docker-compose-actions-workflow

## Action Types

## Variables
* `constant`
* `env`



```yml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build webserver to test
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile
          tags: integrationtests/myservice:latest
      -
        name: Run integration tests against webserver
        uses: dwallach1/jcanary@v1
        with:
          webserver_docker_img: integrationtests/myservice:latest
```
