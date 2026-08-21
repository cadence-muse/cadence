<div align="center">
<a href="https://github.com/nightnoryu/cadence-backend" target="blank">
<img src="./assets/logo.png" alt="Logo" width="90" />
</a>

<h2>Cadence</h2>

</div>

## Overview

Cadence is a simple repertoire manager for gigging bands. Key features include:

- **Bands Participation:** Be a member of several bands and keep them all at a glance.
- **Per-band Repertoire:** Maintain a repertoire of your bands with keys, tempos and other useful information.
- **Per-band Setlists:** Plan and collaborate on setlists for your gigs.

This repo contains only the core backend. Check out the cross-platform client
in [nightnoryu/cadence-client](https://github.com/nightnoryu/cadence-client).

## Tech Stack

- **Backend**: Go, PostgreSQL, Redis
- **Client**: Dart, Flutter
- **CI/CD**: mise, GitHub Actions, GitHub Container Registry
- **Deployment**: k3s
- **Landing website**: Hugo

## Development

### Prerequisites

- [mise](https://mise.jdx.dev)
- Docker with docker-compose-plugin

### First launch

```shell
git clone https://github.com/nightnoryu/cadence-backend
cd cadence-backend

# Set the binary volume mapping according to project folder location on your machine
cp compose.override.example.yml compose.override.yml
$EDITOR compose.override.yml

mise run
docker compose up -d
```

## License

Distributed under the MIT License. See [License](/LICENSE) for more information.
