# vinyl_webscraper

`vinyl_webscraper-go` is a Golang re-implementation of another one of my [repos](https://github.com/1602077/vinyl_pricechecker): it uses Go's Colly API to concurrently scrape Amazon from an input file of urls and write the pricing data to a postgres db.

## Managing PostgreSQL Container
- Deployment of pg container will automatically create the required tables by running `data/schema.sql`.
- This can be done manually using `docker exec -i  pg psql -d webscraper -U root < data/schema.sql`.
- Run `psql` inside of postgres container using `docker exec -it pg psql -d webscraper -U root`.

