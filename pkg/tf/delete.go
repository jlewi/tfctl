package tf

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pkg/errors"
	"os"
)

// DeleteResourceFromFile removes the resource with the given type and name from a Terraform file.
func DeleteResourceFromFile(filePath string, resourceType string, resourceName string) error {
	// Parse the Terraform file using HCL
	parser := hclparse.NewParser()
	hclFile, diags := parser.ParseHCLFile(filePath)
	if diags.HasErrors() {
		return errors.Wrapf(diags, "error parsing file")
	}

	// Create a new HCL write file for storing modified content
	newFile := hclwrite.NewEmptyFile()

	resourceBlocks := body.Blocks()

	content, _, diags := hclFile.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "resource"},
			{Type: "data"},
			// Add other top-level block types you're interested in
		},
	})
	if diags.HasErrors() {
		return errors.Wrapf(diags, "error getting partial content")
	}

	foundResource := false
	// Iterate through the body of the HCL file
	for _, block := range content.Blocks {
		// Check if the block is a resource block
		if block.Type == "resource" {
			labels := block.Labels
			if len(labels) >= 2 && labels[0] == resourceType && labels[1] == resourceName {
				foundResource = true
				// Skip the resource block that matches the resource type and name
				continue
			}
		}
		// Convert the hcl.Block to hclwrite.Block
		newBlock := hclwrite.NewBlock(block.Type, block.Labels)

		// Write the block attributes and nested blocks
		writeBlockBody(hclFile, block.Body, newBlock.Body())

		// Append all other blocks (or non-matching resource blocks) to the new file
		newFile.Body().AppendBlock(newBlock)
	}

	if !foundResource {
		return errors.Errorf("resource %s.%s not found in %s", resourceType, resourceName, filePath)
	}
	// Write the modified content back to the file
	err := os.WriteFile(filePath, newFile.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}

	fmt.Printf("Resource %s.%s deleted successfully from %s!\n", resourceType, resourceName, filePath)
	return nil
}

func writeBlockBody(srcFile *hcl.File, src hcl.Body, dest *hclwrite.Body) {
	content, _ := src.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{{Name: "*"}},
		Blocks:     []hcl.BlockHeaderSchema{{Type: "*"}},
	})

	for _, attr := range content.Attributes {
		// Try to get the value of the expression
		val, diags := attr.Expr.Value(nil)
		if !diags.HasErrors() {
			// If we can get a value, set it directly
			dest.SetAttributeValue(attr.Name, val)
		} else {
			// TODO(jeremy)
			//fmt.Printf("WHAT TO DO HERE")
			// If we can't get a value, use the expression's source code
			srcRange := attr.Expr.Range()
			exprSrc := srcRange.SliceBytes(srcFile.Bytes)
			tokens := hclwrite.Tokens{
				{
					Type:  hclsyntax.TokenIdent,
					Bytes: exprSrc,
				},
			}
			dest.SetAttributeRaw(attr.Name, tokens)
		}
	}

	for _, block := range content.Blocks {
		newBlock := dest.AppendNewBlock(block.Type, block.Labels)
		writeBlockBody(srcFile, block.Body, newBlock.Body())
	}
}
