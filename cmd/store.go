package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Command to upload data to storj V3 network.",
	Long:  `Command to connect and transfer all tables' data from a desired MongoDB instance to given Storj Bucket.`,
	Run:   mongoStore,
}

func init() {

	// Setup the store command with its flags.
	rootCmd.AddCommand(storeCmd)
	var defaultMongoFile string
	var defaultStorjFile string
	storeCmd.Flags().BoolP("accesskey", "a", false, "Connect to storj using access key(default connection method is by using API Key).")
	storeCmd.Flags().BoolP("share", "s", false, "For generating share access of the uploaded backup file.")
	storeCmd.Flags().StringVarP(&defaultMongoFile, "mongo", "m", "././config/db_property.json", "full filepath contaning MongoDB configuration.")
	storeCmd.Flags().StringVarP(&defaultStorjFile, "storj", "u", "././config/storj_config.json", "full filepath contaning storj V3 configuration.")
}

func mongoStore(cmd *cobra.Command, args []string) {

	// Process arguments from the CLI.
	mongoConfigfilePath, _ := cmd.Flags().GetString("mongo")
	fullFileNameStorj, _ := cmd.Flags().GetString("storj")
	useAccessKey, _ := cmd.Flags().GetBool("accesskey")
	useAccessShare, _ := cmd.Flags().GetBool("share")

	// Read MongoDB instance's configurations from an external file and create an MongoDB configuration object.
	configMongoDB := LoadMongoProperty(mongoConfigfilePath)

	// Read storj network configurations from and external file and create a storj configuration object.
	storjConfig := LoadStorjConfiguration(fullFileNameStorj)

	// Connect to storj network using the specified credentials.
	access, project := ConnectToStorj(fullFileNameStorj, storjConfig, useAccessKey)

	// Establish connection with MongoDB and create the customized reader to implement streaming
	reader := ConnectToDB(configMongoDB)
	// Fetch all backup files from MongoDB instance and simultaneously store them into desired Storj bucket.
	fmt.Printf("Initiating back-up.\n\n")
	uploadFileName := path.Join(configMongoDB.Database, configMongoDB.Database+time.Now().Format("2006-01-02_15_04_05"))
	UploadData(project, storjConfig, uploadFileName, reader, reader.collectionNames[0])
	fmt.Printf("\nBack-up complete.\n\n")

	// Create restricted shareable serialized access if share is provided as argument.
	if useAccessShare {
		ShareAccess(access, storjConfig)
	}
}
