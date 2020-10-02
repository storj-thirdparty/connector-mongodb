package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	"testing"

	"github.com/storj-thirdparty/connector-mongodb/cmd"
	"storj.io/uplink"
)

func TestMongoStore(t *testing.T) {

	storjConfig := cmd.LoadStorjConfiguration("../config/storj_config_test.json")
	_, project := cmd.ConnectToStorj(storjConfig, false)

	// Converting JSON data to bson data.  TODO: convert to BSON using call to mongo library
	bsonData, _ := json.Marshal("{'testKey': 'testValue'}")

	// Create a buffer as an io.Reader implementor.
	buf1 := bytes.NewBuffer(bsonData)

	fmt.Printf("Initiating back-up.\n")
	uploadFileName := path.Join("testdb", "testdb"+time.Now().Format("2006-01-02_15_04_05"))
	cmd.UploadData(project, storjConfig, uploadFileName, buf1, "testdb")
	fmt.Printf("Back-up complete.\n\n")

}

func TestMongoReStore(t *testing.T) {

	storjConfig := cmd.LoadStorjConfiguration("../config/storj_config_test.json")
	_, project := cmd.ConnectToStorj(storjConfig, false)

	fmt.Printf("Initiating Restore.")
	cmd.RestoreData(project, "connectortest/testdb", true, false)

	fmt.Printf("\nDeleting the test back-up.\n")
	ctx := context.Background()
	backups := project.ListObjects(ctx, storjConfig.Bucket, &uplink.ListObjectsOptions{Prefix: "testdb/"})
	// Loop to find the latest back-up of all the back-ups.
	for backups.Next() {
		item := backups.Item()
		collections := project.ListObjects(ctx, storjConfig.Bucket, &uplink.ListObjectsOptions{Prefix: item.Key})
		for collections.Next() {
			item := collections.Item()
			_, err := project.DeleteObject(ctx, storjConfig.Bucket, item.Key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("Deleted the test back-up.\n\n")
}
