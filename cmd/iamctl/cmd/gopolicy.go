/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/openshift/aws-load-balancer-operator/cmd/iamctl/aws"
	"github.com/spf13/cobra"
)

var (
	// input file specifies the location for the input json.
	inputFile string

	// outputfile specifies the location of the generated code.
	outputFile string

	// pkg specifies the package with which the code is generated.
	pkg string
)

// gopolicyCmd represents the gopolicy command
var gopolicyCmd = &cobra.Command{
	Use:   "gopolicy",
	Short: "Used to generate AWS IAM Policy from policy json.",
	Long:  `Used for creating/updating iam policies required by the aws loadbalancer operator.`,
	Run: func(cmd *cobra.Command, args []string) {
		aws.GenerateIAMPolicy(inputFile, outputFile, pkg)
	},
}

func init() {
	rootCmd.AddCommand(gopolicyCmd)

	gopolicyCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Used to specify input JSON file path.")
	gopolicyCmd.MarkPersistentFlagRequired("input-file")

	gopolicyCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Used to specify output Go file path.")
	gopolicyCmd.MarkPersistentFlagRequired("output-file")

	gopolicyCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "main", "Used to specify output Go file path.")
}
