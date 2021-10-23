# passbasetest

## Running
To run, follow the steps:
- modify the values of the docker-compose to the approprate values
- run `docker compose up`
- Update the value of the FIXER_APIKEY to a proper key


The project is served both via an http and GRPC endpoint, by default the http is served on port `8081` while grpc on port `8080` but that can be changed with 
the `PORT` and `RESTPORT` values respectively

## NOTE
- Fixer only allows conversion from EUR for free accounts, hence validation logic prevents any other coversions
- Generated swagger documentations are in the `proto/v1/<service>` filder and can be used to easily generate a documentation site

## THINGS TO DO
- Abstract the core logic from the service to allow for easy apikey and project manipulation
- APIKEY management
