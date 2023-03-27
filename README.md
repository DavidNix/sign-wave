# Sign Wave

An experiment. Not for production use!

Concurrently signs batches of records with RSA keys. Prevents double sign and concurrent use of the same key.

## Assumptions
* We have control over the all databases.
* Does not need to support large scale.
* Storing private keys in plaintext is acceptable for this experiment, but not in production.
* If private keys are stored elsewhere, we can create private_key records using the public key hash or otherwise as the identifier. Then the signing process looks up the private key from elsewhere.

## Run Locally

Pre-requisites: Install Go and openssl.

```shell
rm -f signwave.db* && go run . schema && ./seed.sh

go run . emit

# In a separate terminal
go run . ingest

# In a separate terminal
go run . worker
```

## Architecture

```
┌─────────────────────┐         ┌─────────────────────┐          ┌─────────────────────┐
│                     │         │                     │          │                     │
│                     │         │                     │          │                     │
│        Emit         │────────▶│       Ingest        │─────────▶│       Worker        │
│                     │         │                     │          │                     │
│                     │         │                     │          │                     │
└─────────────────────┘         └─────────────────────┘          └─────────────────────┘
     Find unsigned                   Assign private                   Sign records      
        records                     key and prep for                    with key        
                                       signature                                        
```

I chose SQLite for the database because it is easy to setup and use. In production, 
I would use a more scalable database like Postgres.

The emit, ingest, and even worker could be combined into a single process leveraging Go's concurrency. 
But to demonstrate a distributed system, I have separated them.

We emit and ingest batches as quickly as possible. Thanks to the maturity and constraints of SQL, we can have workers sign records at their leisure.

### Areas of Improvement (aside from this being an experiment)
* It's odd to have a microservice architecture in which all services query the same database. In this case, we are using the database as the synchronization mechanism.
* A message queue (Kafka, RabbitMQ, etc.) would better control throughput. However, the current architecture demonstrates that the mechanism of message passing doesn't matter. 
* Pass contexts and listen for cancellation.
* The `store` package methods take an `int64` id as arguments. This is error prone as it's easy to mix up ids. Instead, use a custom type for each id, e.g. `SignatureID`, `PrivateKeyID`, etc.
* Monitor query performance and add missing indexes.
* Use prepared statements.
* It's easy to lock SQLite databases. Therefore, a potential improvement is to introduce an API service which has sole access to the database. Other services then call this API.