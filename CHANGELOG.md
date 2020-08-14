# Connector-MongoDB Changelog

## [1.0.5] - 11-08-2020
### Changelog:
* Optimized performance by incresing the buffer size to 1 MB.
* Added `restore` command to restore the back-up.

## [1.0.5] - 25-05-2020
### Changelog:
* Made changes according to latest uplink RC.
* Added cobra cli for user interface.
* Restructured project based on the requirements for cobra cli.
* Changed arguments to optional flags.
* Added `--accesskey` and `--share` flags and removed `key`, `test`, `parse` and `debug` flags.
* Added `--storj` flag to set storj config file path and `--mongo` to set maria mongo file path.

## [1.0.2] - 13-12-2019
### Changelog:
* Changess made according to latest libuplink v0.27.1
* Made changes for uploading larger file.

## [1.0.1] - 06-12-2019
### Changelog:
* Changes made according to updated cli package.
* Added Macroon functionality.
* Added option to access storj using Serialized Scope Key. 
* Added keyword `key` to access storj using API key rather than Serialized Scope Key (defalt).
* Added keyword `restrict` to apply restrictions on API key and provide shareable Serialized Scope Key for users.
* Error handling for various events.