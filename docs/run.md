# Run

> Back-up is uploaded/restored by streaming to/from the Storj network.

The following flags can be used with the `store` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `share` - Generates a restricted shareable serialized access with the restrictions specified in the Storj configuration file.

The following flags  can be used with the `restore` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `match` - Matches to regular expression with the databases whose back-up(s) are uplaoded to Storj network and restores the latest back-up of all the matching databases. It only works with the `latest` flag.
* `latest` - Restores the latest back-up of the specified MongoDB database.
* `database` - Storj path of the database back-up to be restored. Takes only database name if used with `latest` flag.

Once you have built the project you can run the following:

## Get help

```
$ ./connector-mongodb --help
```

## Check version

```
$ ./connector-mongodb --version
```

## Upload back-up data to Storj

```
$ ./connector-mongodb store --local <path_to_mongodb_config_file> --storj <path_to_storj_config_file>
```

## Upload back-up data to Storj bucket using Access Key

```
$ ./connector-mongodb store --accesskey
```

## Upload back-up data to Storj and generate a Shareable Access Key based on restrictions in `storj_config.json`

```
$ ./connector-mongodb store --share
```

## Restore the specified back-up of a database

```
$ ./connector-mongodb restore --path <database_backup_name>
```

> Example: `./connector-mongodb restore --path bucket/uploadPath/db/dbYYYY-MM-DD_HH_MM_SS/`. Here, `bucket/uploadPath/db/dbYYYY-MM-DD_HH_MM_SS/` is the path of the back-up to be restored.

## Restore the latest back-up of the specified tatabase

```
$ ./connector-mongodb restore --path <database_name> --latest
```

> Example: `./connector-mongodb restore --path bucket/uploadPath/db --latest`. Here, `bucket/uploadPath/db` is the path of the database whose latest back-up is to be retored.

## Restore the lastest back-up of the database(s) matching with the regular expression

```
$ ./connector-mongodb restore --match <regex> --path <bucket/uploadPath> --latest
```

> Example: `./connector-mongodb restore --match db.* --path bucket/uploadPath --latest`. Here, `db.*` is the regular expression which is matched with the databases inside `bucket/uploadPath` on storj network.