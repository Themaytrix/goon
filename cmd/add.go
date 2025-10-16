/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to staging area",
	Long:  `Adds files and folders to the index`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")

		if len(args) > 0 {
			for _, arg := range args {
				entries, _ := os.ReadDir(".")
				fmt.Println(arg)

				for _, entry := range entries {
					// handle folders
					if entry.IsDir() {
						continue
					}
					// extract the file extention
          base := entry.Name()
          ext := filepath.Ext(base)

          match := base[:len(base) - len(ext)]
          if match == arg{
    // get entire file path
            absPath, err := filepath.Abs(base)
            if err != nil{
              fmt.Println(err)
              return
            }
            fmt.Println(absPath)
          }
				}
			}
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
