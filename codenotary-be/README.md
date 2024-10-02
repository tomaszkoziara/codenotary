# Codenotary-BE

This is the backend of the Codenotary coding challenge.

## Commands

### Build image

In order to build a docker image run `make build --always-make`.

## API

### Create accounting info

Creates a new accounting information record.

```
POST /api/v0/accountinginfo

Body
{
    "accountNumber": string, required
    "accountName": string, required
    "iban": string, required
    "address": string, required
    "amount": number, required
    "type": one of sending receiving, required
}
```

### List accounting info

Lists accounting informations for a given account name.

```
GET /api/v0/accountinginfo

Query
accountName: string, required
page: number, required
pageSize: number, required
```