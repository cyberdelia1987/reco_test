# Test project

## Setup:

- copy `config.yaml.sample` into `config.yaml` file, or into any other desired name
- build application: `go build -o test_app ./main.go`
- run application: `./test_app -config="config.yaml"`

## Configuration

The `config.yaml` file contains configuration options for different application parts.

Note the following:

- `asana.access_token` - contains personal access token for Asana SaaS requests, that can be obtained
  at https://app.asana.com/0/my-apps
- `data_dumper.path` - path to a directory where dumped data from Asana users/projects endpoints will be saved

