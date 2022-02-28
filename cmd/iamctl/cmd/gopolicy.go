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

// genpolicyCmd represents the genpolicy command
var genpolicyCmd = &cobra.Command{
	Use:   "genpolicy",
	Short: "Used to generate AWS IAM Policy from policy json.",
	Long:  `Used for creating/updating iam policies required by the aws loadbalancer operator.`,
	Run: func(cmd *cobra.Command, args []string) {
		aws.GenerateIAMPolicy(inputFile, outputFile, pkg)
	},
}

func init() {
	rootCmd.AddCommand(genpolicyCmd)

	genpolicyCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Used to specify input JSON file path.")
	genpolicyCmd.MarkPersistentFlagRequired("input-file")

	genpolicyCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Used to specify output Go file path.")
	genpolicyCmd.MarkPersistentFlagRequired("output-file")

	genpolicyCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "main", "Used to specify output Go file path.")
}
