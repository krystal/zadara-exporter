# Zadara Exporter

A Prometheus exporter for Zadara Storage.

## Overview
The Zadara Exporter is a Prometheus exporter designed to collect data from Zadara Storage and expose it as Prometheus metrics. It utilizes otel metrics to gather the data.

## Features
- Collects data from Zadara Storage
- Exposes data as Prometheus metrics
- Configurable through environment variables, a configuration file, or command line flags

## Usage
To start the Zadara Exporter server, run the following command:


```sh
❯ zadara-exporter server 
2024/05/04 14:40:48 INFO starting Metrics server address=:9090 path=/metrics
```

## Configuration

### Environment Variables

The exporter is configured using the following environment variables:

- `ZADARA_API_BASE_URL`: The Zadara API base url to connect to. (default: `https://api.zadara.com`)
- `ZADARA_TOKEN`: The Zadara API token to use for authentication.
- `ZADARA_LISTEN_ADDR`: The addr to expose the metrics on. (default: `:9090`)
- `ZADARA_LISTN_PATH`: The path to expose the metrics on. (default: `/metrics`)

### Config File

The exporter can also be configured using a configuration file. The configuration file is a YAML file with the following structure:

```yaml
api_base_url: 'https://api.zadara.com' # default: 'https://api.zadara.com'
token: '<ZADARA_API_TOKEN>' # required
listen_address: ':9091' # default: ':9090'
listen_path: '/metr' # default: '/metrics'
```

The configuration file can be specified using the `ZADARA_CONFIG_FILE` environment variable.

The exporter will look for the configuration file (`config.yaml`) in the current working directory if the `ZADARA_CONFIG_FILE` environment variable is not set.
It will also look for the configuration file in the following directories if the `ZADARA_CONFIG_FILE` environment variable is not set:
- `/etc/zadara_exporter/config.yaml`
- `$HOME/.zadara-exporter/config.yaml`

### Command Line Flags

The exporter can also be configured using command line flags. The following flags are available:

```sh
Flags:
      --api_base_url string     The base URL of the Zadara Command Centre API (default "https://api.zadara.com")
  -h, --help                    help for server
      --listen_address string   The address to listen on for the metrics server (default ":9090")
      --listen_path string      The path to expose the metrics on (default "/metrics")
      --token string            The API token for the Zadara Command Centre API

Global Flags:
      --config string   The path to the configuration file
```