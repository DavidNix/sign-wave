# Sign Wave

An experiment. Not for production use!

Concurrently signs messages with RSA keys. Prevents double sign and concurrent use of the same key.

## Assumptions
* We have control over the all databases.
* Supports moderate scale, but not large scale.
* Storing private keys in plaintext is acceptable for this experiment, but not in production.
* If private keys are stored elsewhere, we can create private_key records using the public key hash or otherwise as the identifier. Then the signing process looks up the private key from elsewhere.

## Architecture

TODO: Diagram

I chose SQLite for the database because it is easy to setup and use. In production, 
I would use a more scalable database like Postgres.

The emit, ingest, and even worker could be combined into a single process leveraging Go's concurrency. 
But to demonstrate a distributed system, I have separated them.

### Design Highlights:
* The `store` package uses a service object pattern so the methods do not have `*sql.DB` in its signature. This decouples the sql database in case we want to migrate to another database.

### Areas of Improvement (aside from this being an experiment)
* Pass contexts and listen for cancellation.
* The `store` package methods take an `int64` id as arguments. This is error prone as it's easy to mix up ids. Instead, use a custom type for each id, e.g. `SignatureID`, `PrivateKeyID`, etc.
* Monitor query performance and add missing indexes.
* Use prepared statements.