# <img src="docs/diagrams/liaison-logo.svg" height="40" align="absmiddle" alt="" /> Liaison

> **Connector-powered access to devices and apps behind NAT**

[![Go](https://github.com/liaisonio/liaison/actions/workflows/go.yml/badge.svg)](https://github.com/liaisonio/liaison/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/liaisonio/liaison)](https://goreportcard.com/report/github.com/liaisonio/liaison)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Tech](https://img.shields.io/badge/Tech-Go%20%7C%20TypeScript%20%7C%20React-blue)](#)
[![Version](https://img.shields.io/badge/Version-v1.5.0-green)](#)

English | [简体中文](./README_zh.md) | [日本語](./README_ja.md) | [한국어](./README_ko.md) | [Español](./README_es.md) | [Français](./README_fr.md) | [Deutsch](./README_de.md)

![Dashboard](docs/pages/home_en.png)

| Jellyfin (Stream Home Movies Anywhere) | OpenClaw (Use Home AI Anywhere) |
|:---:|:---:|
| ![Jellyfin](docs/pages/jellyfin-ss.png) | ![OpenClaw](docs/pages/openclaw-ss.png) |

[Quick Start](#quick-start) • [Introduction](#introduction) • [Documentation](#documentation) • [Contributing](#contributing)

---

## Introduction

Liaison is an enterprise-grade application access solution that can be enabled or disabled at any time, without exposing ports on your LAN or home network. It provides a complete feature set: automatic app discovery on connected devices, real-time traffic metrics, and secure TLS-encrypted transport.

This project addresses:

- **Private network access** — Reach devices and services behind NAT from the public internet with minimal setup
- **Multi-device management** — Manage devices across locations with Linux/macOS/Windows support
- **Secure connectivity** — TLS-encrypted transport without exposing ports on your LAN or home network
- **Per-entry firewall** — Source-IP CIDR allowlist on each TCP or HTTP entry, enforced at connection accept
- **Traffic monitoring** — Real-time device status and traffic metrics for operations and capacity planning
- **Application proxy** — TCP, HTTP/HTTPS, WebSocket and other protocols
- **API automation** — Personal Access Tokens (PAT) for CLI/scripts with a browser-mediated sign-in flow at `/cli-auth`

Use cases:

<div align="center">

| **💼 Remote Work & Dev** | **🧑‍💻 Personal Studio** | **🏠 Home Network / NAS** | **🌐 Multi-datacenter / Multi-region** | **⚡ Edge & Ops** |
|:---:|:---:|:---:|:---:|:---:|
| Connect office and home devices for remote development and debugging | Securely connect workstations and private environments with unified device management | Access home NAS and smart-home services from the public internet | Unified connectivity for servers and applications across regions and datacenters | Connect and monitor edge applications with remote health and traffic checks |

</div>

---

## Quick Start

Pick one of the two server deployment options, then install a connector.

### Install Server — Option 1: Binary + systemd

**1. Download**

```bash
wget https://github.com/liaisonio/liaison/releases/download/v1.5.0/liaison-1.5.0-linux-amd64.tar.gz
tar -xzf liaison-1.5.0-linux-amd64.tar.gz
cd liaison-1.5.0-linux-amd64
```

**2. Run install script**

```bash
sudo ./install.sh
```

You will be prompted for a public IP or domain; if none is entered within 30 seconds, the detected public IP is used.

**3. Open Web console**

Visit `https://your-public-ip` to access the Web console.

> **Tip:** Default admin credentials are shown in the install script output or config.

### Install Server — Option 2: Docker Compose

Requires Docker 20.10+ with the `docker compose` plugin. The bundle ships `liaison` + `frontier` as two containers; images are pre-built — no registry or source checkout needed.

```bash
wget https://github.com/liaisonio/liaison/releases/download/v1.5.0/liaison-1.5.0-docker-amd64.tar.gz
tar -xzf liaison-1.5.0-docker-amd64.tar.gz
cd liaison-1.5.0-docker-amd64
./load.sh
```

`load.sh` auto-detects your public IP (with a 30-second countdown prompt), loads the bundled images, starts the stack, and prints the one-time admin password when liaison is ready. Save the password and open `https://<public-ip>` to log in.

Data (`data/` SQLite), TLS certs (`certs/`), and logs (`logs/`) are bind-mounted next to `docker-compose.yaml` for persistence. See [`deploy/docker/README.md`](deploy/docker/README.md) for source builds, upgrade / reset / reverse-proxy / custom-cert recipes.

### Install Connector

**Create a new connector** in the Web console, copy the install command for your platform from the UI, and run it on the target device. The connector will appear in the console automatically.

---

## System Requirements

| Component | Requirements |
|:---|:---|
| **Server** | Linux (Ubuntu 20.04+ or CentOS 7+ recommended) |
| **Connector** | Linux / macOS / Windows (x86_64 and ARM64) |
| **Browser** | Chrome 90+, Firefox 88+, Safari 14+, Edge 90+ |

---

## Architecture

<img src="./docs/diagrams/liaison.png" width="80%">

Liaison uses a centralized architecture with Frontier managing all connectors.

**Components**

- **Liaison** — Web UI and API, plus application entry points
- **Frontier** — Connector gateway that handles connector connections and traffic routing
- **Edge** — Connector client on target devices

---

## Feature Showcase

| Feature | Screenshot |
|:---:|:---:|
| Device Management | ![Device](docs/pages/device_en.png) |
| Application Management | ![Application](docs/pages/application_en.png) |
| Proxy Configuration | ![Proxy](docs/pages/proxy_en.png) |
| Edge Management | ![Edge](docs/pages/edge_en.png) |

---

## Documentation

- [Business flow](./docs/biz_sequence.md)
- [API](./docs/swagger/)

---

## Contributing

Contributions are welcome.

- [Report a bug](https://github.com/liaisonio/liaison/issues/new?template=bug_report.md)
- [Suggest a feature](https://github.com/liaisonio/liaison/issues/new?template=feature_request.md)
- [Open a PR](https://github.com/liaisonio/liaison/pulls)
- [Improve docs](https://github.com/liaisonio/liaison/issues/new?template=documentation.md)

1. Fork the repo  
2. Create a branch (`git checkout -b feature/AmazingFeature`)  
3. Commit (`git commit -m 'Add some AmazingFeature'`)  
4. Push (`git push origin feature/AmazingFeature`)  
5. Open a Pull Request  

---

## License

[Apache License 2.0](LICENSE).

---

<div align="center">

**If this project helps you, please give it a ⭐ Star!**

Made with ❤️ by [Liaison Contributors](https://github.com/liaisonio/liaison/graphs/contributors)

[GitHub](https://github.com/liaisonio/liaison) • [Issues](https://github.com/liaisonio/liaison/issues) • [Discussions](https://github.com/liaisonio/liaison/discussions)

</div>
