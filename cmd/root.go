/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/timharris777/go-spawn/internal/utils"
	"github.com/timharris777/go-spawn/internal/version"
)

var inputFilePath string
var inputPipe bool
var templatePath string
var templatePipe bool
var outputPath string
var debug bool
var versionFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-spawn",
	Short: "A cli tool written in go for project templating, scaffolding, and text-replacement",
	Long:  `A cli tool written in go for project templating, scaffolding, and text-replacement`,
	Run: func(cmd *cobra.Command, args []string) {

		// Set debug logging if specified
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if versionFlag {
			version.PrintVersion()
		}

		// Validate proper flag usage
		err := utils.FlagValidation(inputFilePath, inputPipe, templatePath, templatePipe, outputPath)
		if err != nil {
			panic(err)
		}
		// Get template string
		template, err := utils.GetTemplate(templatePath, templatePipe)
		if err != nil {
			panic(err)
		}
		// Get input data
		input, err := utils.GetInput(inputFilePath, inputPipe)
		if err != nil {
			panic(err)
		}
		// Render template with input
		rendered, err := utils.RenderTemplate(template, input)
		// Print rendered output
		fmt.Printf(rendered)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&inputFilePath, "input", "i", "", "provide a yaml file that has inputs for templating")
	rootCmd.PersistentFlags().BoolVarP(&inputPipe, "inputFromPipe", "", false, "provide input from pipe")
	rootCmd.PersistentFlags().StringVarP(&templatePath, "template", "t", "", "path to template file or folder. Folder requires --output option")
	rootCmd.PersistentFlags().BoolVarP(&templatePipe, "templateFromPipe", "", false, "provide template from pipe")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "folder to output rendered templates")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "folder to output rendered templates")
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "print the current version of go-spawn")
}
