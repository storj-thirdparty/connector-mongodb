# connector-mongodb (uplink v1.0.5)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d637c9ce9f574fb9a1bab1a8f1960a1d)](https://app.codacy.com/gh/storj-thirdparty/connector-mongodb?utm_source=github.com&utm_medium=referral&utm_content=storj-thirdparty/connector-mongodb&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/storj-thirdparty/connector-mongodb)](https://goreportcard.com/report/github.com/storj-thirdparty/connector-mongodb)
![Cloud Build](https://storage.googleapis.com/storj-utropic-services-badges/builds/connector-mongodb/branches/master.svg)


## Overview

The mongoDB Connector connects to an mongoDB database, takes a backup of the specified database and uploads the backup data on Storj network. It can also restore the latest back-up of a specified database to the local storage.

```bash
Usage:
  connector-mongodb [command] <flags>

Available Commands:
  help        Help about any command
  restore	  Command to restore the latest back-up to the local disk
  store       Command to upload data to a Storj V3 network
  version     Prints the version of the tool

```

`store` - Connect to the specified database (default: `db_property.json`). Back-up of the database is generated using tooling provided by mongoDB and then uploaded to the Storj network. Connect to a Storj v3 network using the access specified in the Storj configuration file (default: `storj_config.json`). 

 Back-up data is iterated through and upload in 1 MB chunks to the Storj network.

The following flags  can be used with the `store` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `share` - Generates a restricted shareable serialized access with the restrictions specified in the Storj configuration file.

`restore` - Connect to a Storj v3 network using the access specified in the Storj configuration file (default: `storj_config.json`). Latest back-up of the particular database is located and downloaded to local storage. 

The following flags  can be used with the `restore` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `match` - Matches to regular expression with the databases whose back-up(s) are uplaoded to Storj network and restores the latest back-up of all the matching databases. It only works with the `latest` flag.
* `latest` - Restores the latest back-up of the specified MongoDB database.
* `database` - Storj path of the database back-up to be restored. Takes only database name if used with `latest` flag.

Sample configuration files are provided in the `./config` folder. 

## Requirements and Install

To build from scratch, [install the latest Go](https://golang.org/doc/install#install).

> Note: Ensure go modules are enabled (GO111MODULE=on)

### Option #1: clone this repo (most common)

To clone the repo

```
git clone https://github.com/storj-thirdparty/connector-mongodb.git
```

Then, build the project using the following:

```
cd connector-mongodb
go build
```

### Option #2:  ``go get`` into your gopath

 To download the project inside your GOPATH use the following command:

```
go get github.com/storj-thirdparty/connector-mongodb
```

## Run (short version)

Once you have built the project run the following commands as per your requirement:

### Get help

```
$ ./connector-mongodb --help
```

### Check version

```
$ ./connector-mongodb --version
```

### Create backup from mongoDB and upload to Storj

```
$ ./connector-mongodb store
```

### Restore latest backup from from Storj and save to local disk

```
$ ./connector-mongodb restore
```

## Documentation

For more information on runtime flags, configuration, testing, and diagrams, check out the [Detail](//github.com/storj-thirdparty/connector-mongodb/wiki) or jump to:

* [Config Files](//github.com/storj-thirdparty/connector-mongodb/wiki/#config-files)
* [Run (long version)](//github.com/storj-thirdparty/connector-mongodb/wiki/#run)
* [Testing](//github.com/storj-thirdparty/connector-mongodb/wiki/#testing)
* [Flow Diagram](//github.com/storj-thirdparty/connector-mongodb/wiki/#flow-diagram)
