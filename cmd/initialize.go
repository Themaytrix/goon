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
)

// initializeCmd represents the initialize command
var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize goon repository",
	Long:  `Inintializes the goon repository by instantiating .goon directory`,
	Run: func(cmd *cobra.Command, args []string) {

      
		fmt.Println("initialize called")
		// fmt.Println("current working directory:" ,currPath)
		currPath, _ := os.Getwd()
		_,found := utils.IsGoonRepo(currPath, ".goon")

		if found {
			fmt.Println("This is already a goon repository")
		} else {
			CreateDirectories()
		}
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateDirectories() {
	currPath, _ := os.Getwd()
	rootPath := currPath + "/.goon"

	dirs := []string{
		".goon/objects",
		".goon/objects/info",
		".goon/objects/pack",
		".goon/refs",
		".goon/info",
		".goon/hooks",
	}

	files := []string{
		".goon/HEAD",
		".goon/config",
		".goon/description",
		".goon/index",
	}
	// create the directories and files
	_,isgoon := utils.IsGoonRepo(currPath, ".goon")
	fmt.Println(isgoon)
	if !isgoon {
    fmt.Printf("Initializing goon at %s:",rootPath)

		// directories
		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0o755)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		// files
		for _, file := range files {
      f,err := os.Create(filepath.Clean(file))
      if err != nil{
        fmt.Printf("Error creating file %s: %v",file,err)
        return
      }
      defer f.Close()
		}
	}
}

