## <img src="docs/diagrams/liaison-logo.svg" height="32" alt="" /> Liaison

[English](./README.md) | [简体中文](./README_zh.md) | 日本語 | [한국어](./README_ko.md) | [Español](./README_es.md) | [Français](./README_fr.md) | [Deutsch](./README_de.md)

[![Go](https://github.com/liaisonio/liaison/actions/workflows/go.yml/badge.svg)](https://github.com/liaisonio/liaison/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/liaisonio/liaison)](https://goreportcard.com/report/github.com/liaisonio/liaison)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Tech](https://img.shields.io/badge/Tech-Go%20%7C%20TypeScript%20%7C%20React-blue)](#)
[![Version](https://img.shields.io/badge/Version-v1.5.0-green)](#)

> **NAT の内側にあるデバイスとアプリケーションへ、コネクター経由でアクセス**

![Dashboard](docs/pages/home_en.png)

| Jellyfin(どこからでも家の映画をストリーミング) | OpenClaw(どこからでも家の AI を利用) |
|:---:|:---:|
| ![Jellyfin](docs/pages/jellyfin-ss.png) | ![OpenClaw](docs/pages/openclaw-ss.png) |

[クイックスタート](#クイックスタート) • [概要](#概要) • [ドキュメント](#ドキュメント) • [コントリビュート](#コントリビュート)

---

## 概要

Liaison はエンタープライズ向けのアプリケーション接続ソリューションで、いつでも有効化・無効化でき、LAN や家庭ネットワークのポートを公開する必要がありません。接続済みデバイスのアプリを自動検出し、リアルタイムのトラフィックメトリクスを取得し、TLS 暗号化された安全な通信を提供します。

本プロジェクトが解決する課題:

- **プライベートネットワークへのアクセス** — NAT 配下のデバイスやサービスに、最小限の設定でインターネット経由から到達
- **マルチデバイス管理** — Linux/macOS/Windows をまたいで複数拠点のデバイスを一元管理
- **セキュアな接続** — LAN や家庭ネットワークのポートを晒さずに TLS 暗号化で通信
- **エントリー単位のファイアウォール** — TCP / HTTP 各エントリーに送信元 IP CIDR の許可リストを設定でき、接続受け入れ段階で遮断
- **トラフィックモニタリング** — 運用・キャパシティプランニング向けのリアルタイムなデバイス状態とトラフィック統計
- **アプリケーションプロキシ** — TCP、HTTP/HTTPS、WebSocket などのプロトコル対応
- **API の自動化** — CLI / スクリプト用の Personal Access Token (PAT)、`/cli-auth` でブラウザ経由のサインインフロー

ユースケース:

<div align="center">

| **💼 リモートワーク & 開発** | **🧑‍💻 個人スタジオ** | **🏠 ホームネット / NAS** | **🌐 マルチリージョン** | **⚡ エッジ & 運用** |
|:---:|:---:|:---:|:---:|:---:|
| オフィスと自宅のデバイスを接続し、リモート開発・デバッグ | ワークステーションと私的な環境をセキュアに接続し、機材を統合管理 | 家庭の NAS やスマートホームをインターネットから利用 | 複数リージョン・DC をまたぐサーバーとアプリケーションの一元接続 | エッジアプリのリモート健全性確認とトラフィック監視 |

</div>

---

## クイックスタート

2 つのサーバー導入方式からいずれかを選び、その後にコネクターを導入します。

### サーバー導入 — オプション 1: バイナリ + systemd

**1. ダウンロード**

```bash
wget https://github.com/liaisonio/liaison/releases/download/v1.5.0/liaison-1.5.0-linux-amd64.tar.gz
tar -xzf liaison-1.5.0-linux-amd64.tar.gz
cd liaison-1.5.0-linux-amd64
```

**2. インストールスクリプトを実行**

```bash
sudo ./install.sh
```

公開 IP またはドメインの入力が求められます。30 秒以内に入力しない場合は検出された公開 IP が自動で使用されます。

**3. Web コンソールを開く**

`https://公開IP` にアクセスして Web コンソールへ。

> **ヒント:** 初期管理者の認証情報は install.sh の出力、あるいは設定ファイルに記載されます。

### サーバー導入 — オプション 2: Docker Compose

Docker 20.10+ と `docker compose` プラグインが必要です。`liaison`(Web コンソール + API)と `frontier`(コネクターゲートウェイ)の 2 コンテナを提供し、イメージは同梱済みのため、レジストリからの pull やソースの取得は不要です。

```bash
wget https://github.com/liaisonio/liaison/releases/download/v1.5.0/liaison-1.5.0-docker-amd64.tar.gz
tar -xzf liaison-1.5.0-docker-amd64.tar.gz
cd liaison-1.5.0-docker-amd64
./load.sh
```

`load.sh` は公開 IP を自動検出(30 秒カウントダウンの確認付き)し、イメージをロードしてスタックを起動、liaison の準備ができ次第、一度限りの管理者パスワードを表示します。パスワードを控えた上で `https://<公開IP>` からログインしてください。

データ(`data/` の SQLite)、TLS 証明書(`certs/`)、ログ(`logs/`)は永続化のため `docker-compose.yaml` と同階層に bind mount されます。ソースビルドやアップグレード / 初期化 / リバースプロキシ / カスタム証明書などの詳細は [`deploy/docker/README.md`](deploy/docker/README.md) を参照してください。

### コネクター導入

Web コンソールで **コネクターを新規作成** し、UI から対象プラットフォーム用のインストールコマンドをコピーして対象デバイスで実行すると、コネクターがコンソールに自動で現れます。

---

## 動作要件

| コンポーネント | 要件 |
|:---|:---|
| **サーバー** | Linux(Ubuntu 20.04+ または CentOS 7+ 推奨) |
| **コネクター** | Linux / macOS / Windows(x86_64 および ARM64) |
| **ブラウザ** | Chrome 90+、Firefox 88+、Safari 14+、Edge 90+ |

---

## アーキテクチャ

<img src="./docs/diagrams/liaison.png" width="80%">

Liaison は Frontier ですべてのコネクターを管理する中央集権型アーキテクチャを採用しています。

**コンポーネント**

- **Liaison** — Web UI と API、アプリケーションのエントリーポイント
- **Frontier** — コネクターの接続とトラフィックルーティングを担うゲートウェイ
- **Edge** — 対象デバイス上で動作するコネクタークライアント

---

## 機能紹介

| 機能 | スクリーンショット |
|:---:|:---:|
| デバイス管理 | ![Device](docs/pages/device_en.png) |
| アプリケーション管理 | ![Application](docs/pages/application_en.png) |
| プロキシ設定 | ![Proxy](docs/pages/proxy_en.png) |
| コネクター管理 | ![Edge](docs/pages/edge_en.png) |

---

## ドキュメント

- [ビジネスフロー](./docs/biz_sequence.md)
- [API](./docs/swagger/)

---

## コントリビュート

コントリビュート歓迎です。

- [バグ報告](https://github.com/liaisonio/liaison/issues/new?template=bug_report.md)
- [機能提案](https://github.com/liaisonio/liaison/issues/new?template=feature_request.md)
- [PR を開く](https://github.com/liaisonio/liaison/pulls)
- [ドキュメント改善](https://github.com/liaisonio/liaison/issues/new?template=documentation.md)

1. リポジトリを Fork
2. ブランチを作成(`git checkout -b feature/AmazingFeature`)
3. コミット(`git commit -m 'Add some AmazingFeature'`)
4. Push(`git push origin feature/AmazingFeature`)
5. Pull Request を開く

---

## ライセンス

[Apache License 2.0](LICENSE)

---

<div align="center">

**プロジェクトが役立ったら ⭐ Star をお願いします!**

Made with ❤️ by [Liaison Contributors](https://github.com/liaisonio/liaison/graphs/contributors)

[GitHub](https://github.com/liaisonio/liaison) • [Issues](https://github.com/liaisonio/liaison/issues) • [Discussions](https://github.com/liaisonio/liaison/discussions)

</div>
