<div align="center">
<a href="https://github.com/nightnoryu/cadence-backend" target="blank">
<img src="./assets/logo.png" alt="Logo" width="90" />
</a>

<h2>Cadence</h2>

</div>

## 💡 Overview

Cadence is a simple repertoire manager for actively gigging bands. Key features include:

- **Bands Participation:** Be a member of several bands and keep them all at a glance.
- **Per-band Repertoire:** Maintain a repertoire of your bands with keys, tempos and other useful information.
- **Per-band Setlists:** Plan and collaborate on setlists for your gigs.

This repo contains only the core backend. Other repositories for this project:

- [cadence-client](https://github.com/nightnoryu/cadence-client) - cross-platform (web & mobile) client for API.
- [cadence-platform](https://github.com/nightnoryu/cadence-platform) - k8s settings for deployment.

## 🛠️ Local development

To get a local copy of the project up and running, follow these steps.

### Prerequisites

- [mise](https://mise.jdx.dev)
- Docker with docker-compose-plugin

### First launch

1. Copy the override file and set the binary volume mapping according to project folder location on your machine

    ```shell
    cp compose.override.example.yml compose.override.yml
    $EDITOR compose.override.yml
    ```

2. Build the binary with mise

    ```shell
    mise run
    ```

3. Start the application with `docker compose`

    ```shell
    docker compose up -d
    ```

### Managing the runtime

Use `docker compose` to manage the application:

```shell
# Build & restart to apply changes
mise run && docker restart cadence

# Stop the application
docker compose down
```

## 📜 License

Distributed under the MIT License. See [License](/LICENSE) for more information.
