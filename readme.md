<div align="center">
  <img src="/ghost.svg" width="180" alt="shorty logo showing a gopher ghost"/>
  <h1 align="center">shorty</h1>
</div>

<p align="center">
  The no fluff, straight-forward cloud native url-shortener.
</p>


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
