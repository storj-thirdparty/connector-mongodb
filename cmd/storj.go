// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	progressbar "github.com/cheggaaa/pb/v3"
	"storj.io/uplink"
)

// ConfigStorj depicts keys to search for within the stroj_config.json file.
type ConfigStorj struct {
	APIKey               string `json:"apikey"`
	Satellite            string `json:"satellite"`
	Bucket               string `json:"bucket"`
	UploadPath           string `json:"uploadPath"`
	EncryptionPassphrase string `json:"encryptionpassphrase"`
	SerializedAccess     string `json:"serializedAccess"`
	AllowDownload        string `json:"allowDownload"`
	AllowUpload          string `json:"allowUpload"`
	AllowList            string `json:"allowList"`
	AllowDelete          string `json:"allowDelete"`
	NotBefore            string `json:"notBefore"`
	NotAfter             string `json:"notAfter"`
}

// LoadStorjConfiguration reads and parses the JSON file that contain Storj configuration information.
func LoadStorjConfiguration(fullFileName string) ConfigStorj {

	var configStorj ConfigStorj
	fileHandle, err := os.Open(filepath.Clean(fullFileName))
	if err != nil {
		log.Fatal("Could not load storj config file: ", err)
	}

	jsonParser := json.NewDecoder(fileHandle)
	if err = jsonParser.Decode(&configStorj); err != nil {
		log.Fatal(err)
	}

	// Close the file handle after reading from it.
	if err = fileHandle.Close(); err != nil {
		log.Fatal(err)
	}

	// Display storj configuration read from file.
	fmt.Println("\nRead Storj configuration from the ", fullFileName, " file")
	fmt.Println("\nAPI Key\t\t: ", configStorj.APIKey)
	fmt.Println("Satellite	: ", configStorj.Satellite)
	fmt.Println("Bucket		: ", configStorj.Bucket)

	// Convert the upload path to standard form.
	if configStorj.UploadPath != "" {
		if configStorj.UploadPath == "/" {
			configStorj.UploadPath = ""
		} else {
			checkSlash := configStorj.UploadPath[len(configStorj.UploadPath)-1:]
			if checkSlash != "/" {
				configStorj.UploadPath = configStorj.UploadPath + "/"
			}
		}
	}

	fmt.Println("Upload Path\t: ", configStorj.UploadPath)
	fmt.Println("Serialized Access Key\t: ", configStorj.SerializedAccess)
	return configStorj
}

// ShareAccess generates and prints the shareable serialized access
// as per the restrictions provided by the user.
func ShareAccess(access *uplink.Access, configStorj ConfigStorj) {

	allowDownload, _ := strconv.ParseBool(configStorj.AllowDownload)
	allowUpload, _ := strconv.ParseBool(configStorj.AllowUpload)
	allowList, _ := strconv.ParseBool(configStorj.AllowList)
	allowDelete, _ := strconv.ParseBool(configStorj.AllowDelete)
	notBefore, _ := time.Parse("2006-01-02_15:04:05", configStorj.NotBefore)
	notAfter, _ := time.Parse("2006-01-02_15:04:05", configStorj.NotAfter)

	permission := uplink.Permission{
		AllowDownload: allowDownload,
		AllowUpload:   allowUpload,
		AllowList:     allowList,
		AllowDelete:   allowDelete,
		NotBefore:     notBefore,
		NotAfter:      notAfter,
	}

	// Create shared access.
	sharedAccess, err := access.Share(permission)
	if err != nil {
		log.Fatal("Could not generate shared access: ", err)
	}

	// Generate restricted serialized access.
	serializedAccess, err := sharedAccess.Serialize()
	if err != nil {
		log.Fatal("Could not serialize shared access: ", err)
	}
	fmt.Println("Shareable sererialized access: ", serializedAccess)
}

// ConnectToStorj reads Storj configuration from given file
// and connects to the desired Storj network.
// It then reads data property from an external file.
func ConnectToStorj(fullFileName string, configStorj ConfigStorj, accesskey bool) (*uplink.Access, *uplink.Project) {

	var access *uplink.Access
	var cfg uplink.Config

	// Configure the UserAgent
	cfg.UserAgent = "MongoDB"
	ctx := context.Background()
	var err error

	if accesskey {
		fmt.Println("\nConnecting to Storj network using Serialized access.")
		// Generate access handle using serialized access.
		access, err = uplink.ParseAccess(configStorj.SerializedAccess)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("\nConnecting to Storj network.")
		// Generate access handle using API key, satellite url and encryption passphrase.
		access, err = cfg.RequestAccessWithPassphrase(ctx, configStorj.Satellite, configStorj.APIKey, configStorj.EncryptionPassphrase)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open a new porject.
	project, err := cfg.OpenProject(ctx, access)
	if err != nil {
		log.Fatal(err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project
	_, err = project.EnsureBucket(ctx, configStorj.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to Storj network.")
	return access, project
}

// UploadData uploads the backup file to storj network.
func UploadData(project *uplink.Project, configStorj ConfigStorj, uploadFileName string, dbReader io.Reader, firstCollection string) {

	ctx := context.Background()

	// Create an upload handle for the first collection.
	upload, err := project.UploadObject(ctx, configStorj.Bucket, configStorj.UploadPath+uploadFileName+"/"+firstCollection+".bson", nil)
	if err != nil {
		log.Fatal("Could not initiate upload : ", err)
	}
	fmt.Printf("Uploading %s to %s...\n", configStorj.UploadPath+uploadFileName+"/"+firstCollection+".bson", configStorj.Bucket)

	buf := make([]byte, 10485760)
	var err1 = io.ErrShortBuffer
	// Loop to upload and commit each collection one by one.
	for err1 != nil {
		_, err1 = io.CopyBuffer(upload, dbReader, buf)
		if err1 != nil && err1 != io.ErrShortBuffer {
			// Commit the current copied collection.
			err = upload.Commit()
			if err != nil {
				log.Fatal("Could not commit object upload : ", err)
			}
			// Create upload handle for the next collection to be uploaded.
			upload, err = project.UploadObject(ctx, configStorj.Bucket, configStorj.UploadPath+uploadFileName+"/"+currentCollection+".bson", nil)
			if err != nil {
				log.Fatal("Could not initiate upload : ", err)
			}
			fmt.Printf("Uploading %s to %s...\n", configStorj.UploadPath+uploadFileName+"/"+currentCollection+".bson", configStorj.Bucket)
		}
	}

	// Commit the upload after copying the last collection.
	err = upload.Commit()
	if err != nil {
		log.Fatal("Could not commit object upload : ", err)
	}
}

func findLatestBackup(project *uplink.Project, backupPath string) string {

	ctx := context.Background()
	keys := strings.Split(backupPath, "/")
	// Object iterator to traverse all the back-ups of the specified database.
	objects := project.ListObjects(ctx, keys[0], &uplink.ListObjectsOptions{Prefix: backupPath[len(keys[0])+1:] + "/"})
	var backups []string
	// Loop to find the latest back-up of all the back-ups.
	for objects.Next() {
		item := objects.Item()
		backups = append(backups, item.Key)
	}
	sort.Strings(backups)
	if len(backups) == 0 {
		log.Fatal("Error: No back-up to restore!")
	}
	return backups[len(backups)-1]
}

// RestoreData restores the latest backup correspoinding to the path provided
func RestoreData(project *uplink.Project, backupPath string, latest bool, showProgress bool) {

	ctx := context.Background()
	var collections *uplink.ObjectIterator
	keys := strings.Split(backupPath, "/")
	if latest {
		fmt.Printf("Restoring the latest backup of %s...\n", backupPath)
		checkSlash := backupPath[len(backupPath)-1:]
		if checkSlash == "/" {
			backupPath = backupPath[:len(backupPath)-1]
		}
		pathTokens := strings.Split(backupPath, "/")
		if len(pathTokens) > 3 {
			log.Fatal("Error: Invalid regular expression! It should only contain the pattern of database name.\n")
		}
		latestBackup := findLatestBackup(project, backupPath)
		collections = project.ListObjects(ctx, keys[0], &uplink.ListObjectsOptions{Prefix: latestBackup})
	} else {
		fmt.Printf("Restoring the backup of %s...\n", backupPath)
		// Convert the backup path to standard form
		checkSlash := backupPath[len(backupPath)-1:]
		if checkSlash != "/" {
			backupPath = backupPath + "/"
		}
		collections = project.ListObjects(ctx, keys[0], &uplink.ListObjectsOptions{Prefix: backupPath[len(keys[0])+1:]})
	}

	var restored []*uplink.Object
	// Download all the collection back-up files corresponding to the back-up inside the ./dump folder.
	for collections.Next() {
		item := collections.Item()
		download, err := project.DownloadObject(ctx, keys[0], item.Key, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n")
		var bar *progressbar.ProgressBar
		var reader io.ReadCloser
		if showProgress {
			info := download.Info()
			bar = progressbar.New64(info.System.ContentLength)
			reader = bar.NewProxyReader(download)
			bar.Start()
		} else {
			reader = download
		}

		// Read everything from the download stream
		receivedContents, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}

		downloadFileName := filepath.Join("dump", filepath.Base(filepath.Dir(item.Key)), filepath.Base(item.Key))
		_ = os.MkdirAll(filepath.Dir(downloadFileName), 0750)
		err = ioutil.WriteFile(downloadFileName, receivedContents, 0600)
		if err != nil {
			log.Fatal(err)
		}
		restored = append(restored, item)
	}
	if len(restored) == 0 {
		log.Fatal("Error: Nothing to restore as the given path.")
	}
	fmt.Printf("\nBackup of %s restored.\n", keys[len(keys)-2])
}

// MatchAndRestore finds the databases corresponding the pattern entered by the user
// and restores the latest backup of each matching database.
func MatchAndRestore(project *uplink.Project, matchPattern string, backupPath string, latest bool, showProgress bool) {

	keys := strings.Split(backupPath, "/")
	if !latest {
		log.Fatal("Error: match used without `latest` flag!")
	}
	ctx := context.Background()
	if len(keys) == 1 {
		if _, err := project.StatBucket(ctx, backupPath); err != nil {
			log.Fatal(err)
		}
		databases := project.ListObjects(ctx, backupPath, nil)
		for databases.Next() {
			item := databases.Item()
			matched, err := regexp.MatchString(matchPattern, filepath.Base(item.Key))
			if err != nil {
				log.Fatal(err)
			}
			if matched {
				fmt.Println("Matching database: ", backupPath+"/"+item.Key)
				RestoreData(project, backupPath+"/"+item.Key, latest, showProgress)
			}
		}
	} else {
		if _, err := project.StatBucket(ctx, keys[0]); err != nil {
			log.Fatal(err)
		}
		databases := project.ListObjects(ctx, keys[0], &uplink.ListObjectsOptions{Prefix: backupPath[len(keys[0])+1:] + "/"})
		for databases.Next() {
			item := databases.Item()
			matched, err := regexp.MatchString(matchPattern, filepath.Base(item.Key))
			if err != nil {
				log.Fatal(err)
			}
			if matched {
				fmt.Println("Matching database: ", keys[0]+"/"+item.Key)
				RestoreData(project, keys[0]+"/"+item.Key, latest, showProgress)
			}
		}
	}
}
