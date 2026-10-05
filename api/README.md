# Vintage Jazz Records API

A small in-memory REST API built with Go's standard `net/http` package.

Run it:

```sh
go run ./cmd/server
```

Endpoints:

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/records` | Return all records. |
| `POST` | `/records` | Add a record. |
| `GET` | `/records/{id}` | Return one record. |

Example request:

```sh
curl -i -X POST http://localhost:8080/records \
  -H 'Content-Type: application/json' \
  -d '{"id":"a-love-supreme","title":"A Love Supreme","artist":"John Coltrane","year":1965,"label":"Impulse!"}'
```

The server seeds two records on startup. Data is intentionally in memory, so it resets when the server stops.
