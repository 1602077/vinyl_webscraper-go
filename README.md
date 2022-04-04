# vinyl_webscraper

`vinyl_webscraper-go` is a Golang re-implementation of another one of my [repos](https://github.com/1602077/vinyl_pricechecker): it uses Go's Colly API to concurrently scrape Amazon and write this information to a postgres db.

## Running in Docker

1. `docker compose up`
2. `docker exec -it pg psql -U postgres` - opens psql in container
3. `docker exec -i pg psql -U postgres < data/schema.sql` - create tables in container on first run

