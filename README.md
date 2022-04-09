# vinyl_webscraper

`vinyl_webscraper-go` is a Golang re-implementation of another one of my [repos](https://github.com/1602077/vinyl_pricechecker): it uses Go's Colly API to concurrently scrape Amazon from an input file of urls and write the pricing data to a postgres db.

## Managing PG Container
1. After building containers (`docker compose build`), launch pg container and then run `docker exec -i  pg psql -d webscraper -U root < data/schema.sql` to create tables for database.
2. Confirm these tables have been created by launching `psql` inside of container: `docker exec -it pg psql -d webscraper -U root`.
3. Services can now be deployed with `docker compose up`.

