# dp-legacy-cache-api

REST API for managing cache control information for pages within the legacy CMS

### Getting started

- Ensure Docker is installed on your local machine. Installation steps can be found [here](https://docs.docker.com/desktop/install/mac-install/).
- Run `docker run --name mongo-test -p 27017:27017 -e MONGO_INITDB_DATABASE=cache -v $(pwd)/mongo-init:/docker-entrypoint-initdb.d -d mongo`.
  - This command launches a MongoDB container named `mongo-test`, maps port 27017 from the host to the container, sets `cache` as the default database, runs initialization scripts (located in the `mongo-init` directory), and operates in the background.
- Run `make debug` to run the application on http://localhost:29100.
- By default, the write (PUT) endpoint is disabled. To be able to create or update resources, please follow these steps:
  - Run [Zebedee](https://github.com/ONSdigital/zebedee).
  - Run `IS_PUBLISHING=true make debug`. This will make the PUT endpoint available.
  - Send a valid request to the PUT endpoint. You'll need to set the Bearer token (the `Authorization` header's value should be `Bearer your-token-here`).
    - For local usage, you can use the Service Auth Token specified in the [DP's install guide](https://github.com/ONSdigital/dp/blob/a9ceaa3fb500e5e2850c8b4853bebf922640083b/guides/INSTALLING.md#environment-variables).
    - For Sandbox/Production usage (or to generate a different token), please follow [this guide](https://github.com/ONSdigital/zebedee#service-authentication-with-zebedee).
- Run `make help` to see a full list of make targets.

### Dependencies

- No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                         | Description                                                                                                        |
|------------------------------|---------------------------------|--------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                    | :29100                          | The host and port to bind to                                                                                       |
| MONGODB_BIND_ADDR            | localhost:27017                 | The MongoDB bind address                                                                                           |
| MONGODB_USERNAME             |                                 | The MongoDB Username                                                                                               |
| MONGODB_PASSWORD             |                                 | The MongoDB Password                                                                                               |
| MONGODB_DATABASE             | cache                           | The MongoDB database                                                                                               |
| MONGODB_COLLECTIONS          | CacheTimesCollection:cachetimes | The MongoDB collections                                                                                            |
| MONGODB_REPLICA_SET          |                                 | The name of the MongoDB replica set                                                                                |
| MONGODB_ENABLE_READ_CONCERN  | false                           | Switch to use (or not) majority read concern                                                                       |
| MONGODB_ENABLE_WRITE_CONCERN | true                            | Switch to use (or not) majority write concern                                                                      |
| MONGODB_CONNECT_TIMEOUT      | 5s                              | The timeout when connecting to MongoDB (`time.Duration` format)                                                    |
| MONGODB_QUERY_TIMEOUT        | 15s                             | The timeout for querying MongoDB (`time.Duration` format)                                                          |
| MONGODB_IS_SSL               | false                           | Switch to use (or not) TLS when connecting to mongodb                                                              |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                              | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL         | 30s                             | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                             | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| IS_PUBLISHING                | false                           | Determines if the instance is in publishing or not                                                                 |
| ZEBEDEE_URL                  | http://localhost:8082           | Zebedee host address and port for authentication                                                                   |

Copyright © 2024, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
