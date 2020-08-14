## Flow Diagram

![](https://github.com/storj-thirdparty/connector-mongodb/blob/master/README.assets/arch.drawio.png)

## Config Files

There are two config files that contain Storj network and mongoDB connection information. The tool is designed so you can specify a config file as part of your tooling/workflow. 

### `db_property.json`

Inside the `./config` directory there is a `db_property.json` file, with following information about your mongoDB instance:

* `hostName`- Host Name connect to mongoDB
* `port` - Port connect to mongoDB
* `username` - User Name of mongoDB
* `password` - Password of mongoDB
* `database` - mongoDB Database Name

### `storj_config.json`

Inside the `./config` directory a `storj_config.json` file, with Storj network configuration information in JSON format:

* `apiKey` - API Key created in Storj Satellite GUI (mandatory)
* `satellite` - Storj Satellite URL (mandatory)
* `encryptionPassphrase` - Storj Encryption Passphrase (mandatory)
* `bucket` - Name of the bucket to upload data into (mandatory)
* `uploadPath` - Path on Storj Bucket to store data (optional) or "" or "/" (mandatory)
* `serializedAccess` - Serialized access shared while uploading data used to access bucket without API Key (mandatory)
* `allowDownload` - Set *true* to create serialized access with restricted download (mandatory while using *share* flag)
* `allowUpload` - Set *true* to create serialized access with restricted upload (mandatory while using *share* flag)
* `allowList` - Set *true* to create serialized access with restricted list access
* `allowDelete` - Set *true* to create serialized access with restricted delete
* `notBefore` - Set time that is always before *notAfter*
* `notAfter` - Set time that is always after *notBefore*

## Run

Once you have built the project run the following commands as per your requirement:

### Get help

```
$ ./connector-mongodb --help
```

### Check version

```
$ ./connector-mongodb --version
```

### Create back-up from mongoDB and upload them to Storj

```
$ ./connector-mongodb store --mongob <path_to_mongo_config_file> --storj <path_to_storj_config_file>
```

### Create back-up from mongoDB and upload it to Storj bucket using Access Key

```
$ ./connector-mongodb store --accesskey
```

### Create back-up from mongoDB and upload it to Storj and generate a Shareable Access Key based on restrictions in `storj_config.json`

```
$ ./connector-mongodb store --share
```

## Restore the specified back-up of a database

```
$ ./connector-mongodb restore --path <database_backup_name>
```

> Example: `./connector-mongodb restore --path <bucket/uploadPath/db/dbYYYY-MM-DD_HH_MM_SS/>`. Here, `bucket/uploadPath/db/dbYYYY-MM-DD_HH_MM_SS/` is the path of the back-up to be restored.

## Restore the latest back-up of the specified tatabase

```
$ ./connector-mongodb restore --path <database_name> --latest
```

> Example: `./connector-mongodb restore --path <bucket/uploadPath/db> --latest`. Here, `bucket/uploadPath/db` is the path of the database whose latest back-up is to be retored.

## Restore the lastest back-up of the database(s) matching with the regular expression

```
$ ./connector-mongodb restore --match <regex> --path <bucket/uploadPath> --latest
```

> Example: `./connector-mongodb restore --match db.* --path bucket/uploadPath --latest`. Here, `db.*` is the regular expression which is matched with the databases inside `bucket/uploadPath` on storj network.

## Testing

The project has been tested on the following operating systems:

```
	* Windows
		* Version: 10 Pro
		* Processor: Intel(R) Core(TM) i3-5005U CPU @ 2.00GHz 2.00GHz

	* macOS Catalina
		* Version: 10.15.4
		* Processor: 2.5 GHz Dual-Core Intel Core i5

	* ubuntu
		* Version: 16.04 LTS
		* Processor: AMD A6-7310 APU with AMD Radeon R4 Graphics × 4
```
