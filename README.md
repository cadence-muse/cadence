# Cadence

Core service of the Cadence app.

## Local Development

### Prerequisites

- [mise](https://mise.jdx.dev)
- Docker with docker-compose-plugin

### First launch

```shell
git clone https://github.com/cadence-muse/cadence
cd cadence

# Set the binary volume mapping according to project folder location on your machine
cp compose.override.example.yml compose.override.yml
$EDITOR compose.override.yml

mise run
docker compose up -d
```

## License

Distributed under the MIT License. See [License](/LICENSE) for more information.
