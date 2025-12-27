/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Themaytrix/goon/utils"
	"github.com/spf13/cobra"

	"github.com/Themaytrix/goon/internals/index"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to staging area",
	Long:  `Adds files and folders to the index`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("add called")
		currDir, _ := os.Getwd()
		root, isgoon := utils.IsGoonRepo(currDir, ".goon")
		if isgoon && len(args) > 0 {

			// create slice for files
			var files []string
			for _, arg := range args {
				// check if its file
				if utils.IsFile(arg) {
					files = append(files, arg)
				} else if utils.IsDir(arg) {
					// traverse directory to append the files
					utils.WalkDir(&files, arg)
				}
      }
				// create new index instance
				idx := index.NewIndex()
				// find if goon/index exists
        idxDir := filepath.Join(root, "index")
				if _, err := os.Stat(idxDir); err == nil {
					idx, _ = index.ReadIndex(idxDir)
				}
				// loop through files
        for _,path := range files{
          path = filepath.ToSlash(path)

          // check for deleted files
          if _,err := os.Stat(path); err != nil{
            idx.RemovePath(path)
            continue
          }

          // create new entry
          entry, err := index.NewIndexEntry(path)

          if err != nil {
            panic(err)
          }

          // replace alreading existing entry
          idx.RemovePath(path)
          idx.AddEntry(entry)
        }

      // write to the index file
      idx.WriteIndex(idxDir)

		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
