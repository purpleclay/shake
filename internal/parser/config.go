/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

/*
	prompt "this is question number 1" {
		type = "text"
	}

	prompt "this is question number 2" {
		type = "text"
	}
*/

// PromptsHCL ...
type PromptsHCL struct {
	Prompts []*struct {
		Question string `hcl:",label"`
		Type     string `hcl:"type"`
	} `hcl:"prompt,block"`
}

// TextPrompt ...
type TextPrompt struct {
	Question string
}

// ReadConfig ...
func ReadConfig(f string) ([]TextPrompt, error) {
	parser := hclparse.NewParser()
	src, diag := parser.ParseHCLFile(f)
	if diag.HasErrors() {
		return []TextPrompt{}, fmt.Errorf(
			"error in ReadConfig parsing HCL: %w", diag,
		)
	}

	prmptHCL := &PromptsHCL{}
	if diag := gohcl.DecodeBody(src.Body, &hcl.EvalContext{}, prmptHCL); diag.HasErrors() {
		return []TextPrompt{}, fmt.Errorf(
			"error in ReadConfig decoding HCL configuration: %w", diag,
		)
	}

	prmpts := []TextPrompt{}
	for _, p := range prmptHCL.Prompts {
		switch ptype := p.Type; ptype {
		case "text":
			prmpts = append(prmpts, TextPrompt{Question: p.Question})
		default:
			return []TextPrompt{}, fmt.Errorf("error in ReadConfig: unknown prompt type `%s`", ptype)
		}
	}

	return prmpts, nil
}
