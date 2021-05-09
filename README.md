## Use it
```go install ./cmd/tbb/...```

```tbb help```

## Summary
Blockchain is a distributed database.\
Each time you want to make a transaction (an action which changes the database) we create a new Status.\
This status is based on the database and transactions are ran into it.\
Once the transaction is finished, the current Status is persisted into the database.\
In order to add security, we hash the full database per transaction, this way, each status has a hash asigned and if you want to modify a transaction, the whole hash chain will change and database will be invalidated.
But, this is inefficient because database might be huge.
Thats why we use blocks. They contains information about transaction's batches (we don't process transactions one by one anymore) and a reference to previous block, the hash is calcultaed based on previous block and the transcation batches. This way we don't need to hash the entire database but just our block and the parent information.
