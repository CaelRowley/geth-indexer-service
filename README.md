# EVM Indexer

A work-in-progress EVM indexer and block explorer written in go. The project aims to sync blockchain data onto a PostgreSQL database and provide search and analytics using Elasticsearch.

## Environment Variables

- **DB_URL**: _The URL connection to a PostgreSQL database server._
  
- **NODE_URL**: _WebSocket connection to Ethereum node._

## Flags

- **port**: _Specifies the port number where the service will run. Default is 8080. Use this flag to define a custom port for the service._
  
- **sync**: _Enables block synchronization with the node's database. By default, synchronization is turned off (false). Use this flag to initiate synchronization of blockchain data into the PostgreSQL database._

## Getting Started

You can run the backend locally with [go](https://go.dev/), [make](https://www.gnu.org/software/make/manual/make.html#Introduction) or [Docker](https://docs.docker.com/).

**Note** The backend runs on port `8080` and the DB on port `5432` so make sure those ports are free if running locally.

### Running with Make

```bash
make dev
```

### Running with Docker Compose

```bash
docker-compose -f docker-compose.yml up
```

## Database Table Structure

`blocks` table:


| Column        | Type      | Key       | Description                                                                            |
|---------------|-----------|-----------|----------------------------------------------------------------------------------------|
| hash          | char(66)  | Primary   | The hash of the block header.                                                          |
| number        | numeric   |           | Numeric identifier of the block within the blockchain.                                 |
| gas_limit     | numeric   |           | Maximum gas allowed for transactions in the block.                                     |
| gas_used      | numeric   |           | Total gas consumed by transactions in the block.                                       |
| difficulty    | varchar   |           | Difficulty level for mining this block.                                                |
| time          | numeric   |           | Timestamp of when the block was mined, in seconds since the epoch.                     |
| parent_hash   | char(66)  |           | The hash of the parent block, the previous block in the blockchain.                    |
| nonce         | varchar   |           | A 64-bit hash used in mining to demonstrate PoW for a block. No longer used for PoS.   |
| miner         | char(42)  |           | Address of the miner who mined the block.                                              |
| size          | numeric   |           | Size of the block in bytes.                                                            |
| root_hash     | char(66)  |           | Root hash of transactions in the block.                                                |
| uncle_hash    | char(66)  |           | Hash of the uncle blocks (or ommer blocks) included in this block.                     |
| tx_hash       | char(66)  |           | Hash of all transaction hashes in this block.                                          |
| receipt_hash  | char(66)  |           | Hash of the receipts of all transactions in this block.                                |
| extra_data    | bytea     |           | Additional binary data associated with the block.                                      |



