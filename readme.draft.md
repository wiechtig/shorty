<div align="center">
  <img src="./ghost.svg" width="180" alt="shorty logo showing a gopher ghost"/>
  <h1 align="center">shorty</h1>
</div>

<p align="center">
  The no fluff, straight-forward cloud native url-shortener.
</p>

## Features

- create shortened links
  - with random or custom names
  - with expiration date (shortened links accessible for given time-frames)
- regular checks if links are still valid.
  - if invalid, it throws an alert or get notified via chat (Slack, Microsoft Teams, Google Chat, etc.)
- Analytics to track clicks, referrers, location, browser, etc.
  - pushed to prometheus or long term storage in db and visible via UI
- UI secured with OIDC
- cloud-native docker deployment
  - docker-compose or kubernetes

## Features

### URL Shortening
- **Custom & Random Links**: Create shortened URLs with custom names or auto-generated random identifiers
- **Expiration Control**: Set time-based expiration for links
- **Health Monitoring**: Automated checks to verify target URLs are still accessible

### Analytics & Monitoring
- **Real-time Metrics**: Export data to Prometheus for monitoring and alerting
- **Historical Data**: Long-term storage in database with web-based analytics dashboard
- **Performance Insights**: Detailed reports on link usage patterns and trends

### Security & Authentication
- **OIDC Integration**: Secure web interface with OpenID Connect authentication
- **Access Control**: Role-based permissions for link management
- **Audit Logging**: Track all administrative actions and changes

### Deployment & Infrastructure
- **Cloud-Native**: Designed for containerized environments
- **Flexible Deployment**: Support for Docker Compose and Kubernetes
- **Scalable Architecture**: Horizontal scaling with load balancing support

## Getting Started

```sh
docker compose up -d
```

Open your browser and go to `http://localhost:4444`.

## Alternatives

There are many alternatives out there, feel free to use them if they fit your needs better:
- [awesome-url-shortener](https://github.com/738/awesome-url-shortener)
- [flink](https://gitlab.com/rtraceio/web/flink)

## Contributing

Prerequisites:
- Go 1.24+
- Node.js 20+
- OCI Runtime (like Docker and Docker Compose)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/)

Setup database for local development:
```sh
make start-dev
```

Run go server:
```sh
make run
```

Run tests:
```sh
make test
```

### Load Tests

To run the entire monitoring suite (stop development setup before):
```sh
make run-app
```

Run k6 load tests:
```sh
make k6
```

## Credits

- [Gopher](https://github.com/MariaLetta/free-gophers-pack) by Maria Letta
