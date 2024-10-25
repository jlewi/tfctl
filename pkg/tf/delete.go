package tf

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"os"
)

// DeleteResourceFromFile removes the resource with the given type and name from a Terraform file.
func DeleteResourceFromFile(filePath string, resourceType string, resourceName string) error {
	// Parse the Terraform file using HCL
	// Parse the file and get the root AST
	src, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "error reading file: %v", filePath)

	}

	f, diag := hclwrite.ParseConfig(src, filePath, hcl.InitialPos)

	if diag.HasErrors() {
		return errors.Wrapf(diag, "error parsing file: %s", filePath)
	}

	// Find the resource block
	body := f.Body()
	resourceBlocks := body.Blocks()

	// Create a new HCL write file for storing modified content
	newFile := hclwrite.NewEmptyFile()

	foundResource := false
	// Iterate through the body of the HCL file
	for _, block := range resourceBlocks {
		// Check if the block is a resource block
		if block.Type() == "resource" {
			labels := block.Labels()
			if len(labels) >= 2 && labels[0] == resourceType && labels[1] == resourceName {
				foundResource = true
				// Skip the resource block that matches the resource type and name
				continue
			}
		}

		// Append all other blocks (or non-matching resource blocks) to the new file
		newFile.Body().AppendBlock(block)
		// We need to manually insert a new line after every block.
		// hclwrite.Format doesn't appear to handle that automatically.
		newFile.Body().AppendNewline()
	}

	if !foundResource {
		return errors.Errorf("resource %s.%s not found in %s", resourceType, resourceName, filePath)
	}

	formattedBytes := hclwrite.Format(newFile.Bytes())

	// Write the modified content back to the file
	err = os.WriteFile(filePath, formattedBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}

	return nil
}
