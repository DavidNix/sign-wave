# Sign Wave

Concurrently signs messages with RSA keys. Prevents double sign and concurrent use of the same key.

## Assumptions
* We have control over the all databases.
* Supports moderate scale, but not large scale.
* Storing private keys in plaintext is acceptable for this experiment, but not in production.

## Architecture

TODO: Diagram

I chose SQLite for the database because it is easy to setup and use. In production, 
I would use a more scalable database like Postgres.

The emit, ingest, and even worker could be combined into a single process leveraging Go's concurrency. 
But to demonstrate a distributed system, I have separated them into separate processes.

Design Highlights:
* The `store` package uses a service object pattern so the methods do not have `*sql.DB` in its signature. This decouples the sql database in case we want to migrate to another database.