package cmd

import (
	"fmt"
	"github.com/jlewi/tfctl/pkg/application"
	"github.com/jlewi/tfctl/pkg/tf"
	"github.com/spf13/cobra"
	"os"
)

// NewDeleteCmd returns a command to delete the resource from the terraform file
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <resource type> <resource name> <terraform file>",
		Short: "Delete a resource from a Terraform file",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				if err := tf.DeleteResourceFromFile(args[2], args[0], args[1]); err != nil {
					return err
				}
				fmt.Printf("Successfully deleted resource %s of type %s from file %s", args[1], args[0], args[2])
				return nil
			}()

			if err != nil {
				fmt.Printf("Failed to delete resource from file: %+v", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
