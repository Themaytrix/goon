/*
Copyright © 2025 NAME HERE <obengraymond81@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
"path/filepath"	

	"github.com/Themaytrix/goon/utils"
  
	"github.com/Themaytrix/goon/goonObjects"
	"github.com/spf13/cobra"
)

var printcontext string


// catfileCmd represents the catfile command
var catfileCmd = &cobra.Command{
	Use:   "catfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
    currdir, _ := os.Getwd()
// get goon directory
    goondir,isgoon := utils.IsGoonRepo(currdir,".goon")
    if isgoon && printcontext != ""{
       objectdir := printcontext[:2]
      filep := printcontext[2:]

      fullpath := filepath.Join(goondir,objectdir,filep) 
      goonobjects.ReadObject(fullpath)
    }else{
      fmt.Println("nothing to see here")
    }
	},
}

func init() {

  catfileCmd.Flags().StringVarP(&printcontext, "print", "p", "", "prints context of hash")
  if err := catfileCmd.MarkFlagRequired("print"); err != nil{panic(err)}
	rootCmd.AddCommand(catfileCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
