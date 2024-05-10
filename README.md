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

- `ZADARA_LISTEN_ADDR`: The addr to expose the metrics on. (default: `:9090`)
- `ZADARA_LISTN_PATH`: The path to expose the metrics on. (default: `/metrics`)

### Config File
> [!IMPORTANT]  
> Target configuration is required for the exporter to work. The exporter will not start if there are no targets configured.

The exporter can also be configured using a configuration file. The configuration file is a YAML file with the following structure:

```yaml
listen_address: :9090
listen_path: /metrics
targets:
  - name: London
    url: https://command-center-1.zadarastorage.com
    token: "<TOKEN HERE>"
    cloud_name: cc1
  # - name: New York
  #   url: https://command-center-2.zadarastorage.com
  #   token: "<TOKEN HERE>"
  #   cloud_name: cc2

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
  -h, --help                    help for server
      --listen_address string   The address to listen on for the metrics server (default ":9090")
      --listen_path string      The path to expose the metrics on (default "/metrics")

Global Flags:
      --config string   The path to the configuration file
```