package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Command to restore back-up from storj V3 network.",
	Long:  `Command to connect to storj network and restore latest back-up of the desired MongoDB database from given storj Bucket to local disk.`,
	Run:   mongorestore,
}

func init() {

	// Setup the restore command with its flags.
	rootCmd.AddCommand(restoreCmd)
	var defaultBackupPathStorj string
	var defaultMatchDatabase string
	var defaultStorjFile string
	restoreCmd.Flags().BoolP("progress", "b", true, "if true, show progress.")
	restoreCmd.Flags().StringVarP(&defaultBackupPathStorj, "path", "p", "", "storj path of the back-up to be restored in the format bucket/uploadPath/db/dbYYYY-MM-DD_HH_MM_SS.")
	restoreCmd.Flags().BoolP("latest", "l", false, "to restore the latest back-up.")
	restoreCmd.Flags().StringVarP(&defaultMatchDatabase, "match", "m", "", "pattern to match with the database(s) whose back-up is to be restored.")
	restoreCmd.Flags().StringVarP(&defaultStorjFile, "storj", "s", "././config/storj_config.json", "full filepath contaning storj V3 configuration.")
}

func mongorestore(cmd *cobra.Command, args []string) {

	// Process arguments from the CLI.
	showProgress, _ := cmd.Flags().GetBool("progress")
	matchPattern, _ := cmd.Flags().GetString("match")
	fullFileNameStorj, _ := cmd.Flags().GetString("storj")
	backupPath, _ := cmd.Flags().GetString("path")
	useAccessKey, _ := cmd.Flags().GetBool("accesskey")
	backupLatest, _ := cmd.Flags().GetBool("latest")

	// Read storj network configurations from and external file and create a storj configuration object.
	storjConfig := LoadStorjConfiguration(fullFileNameStorj)

	// Connect to storj network using the specified credentials.
	_, project := ConnectToStorj(fullFileNameStorj, storjConfig, useAccessKey)

	// Restore the backup from specified Storj bucket.
	fmt.Printf("Initiating restore.\n\n")
	if matchPattern != "" {
		pathTokens := strings.Split(matchPattern, "/")
		if len(pathTokens) > 1 {
			log.Fatal("Error: Invalid regular expression! It should only contain the pattern of database name.\n")
		}
		checkSlash := backupPath[len(backupPath)-1:]
		if checkSlash == "/" {
			backupPath = backupPath[:len(backupPath)-1]
		}
		pathTokens = strings.Split(backupPath, "/")
		if len(pathTokens) > 2 {
			log.Fatal("Error: Invalid back-up path!\n")
		}
		MatchAndRestore(project, matchPattern, backupPath, backupLatest, showProgress)
	} else {
		if backupLatest {
			RestoreData(project, backupPath, backupLatest, showProgress)
		} else {
			RestoreData(project, backupPath, backupLatest, showProgress)
		}
	}
}
