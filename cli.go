package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runCLI(args []string) int {
	fs := flag.NewFlagSet("bitwarden-html-converter", flag.ContinueOnError)
	input := fs.String("i", "", "Bitwarden JSON export file (required)")
	output := fs.String("o", "", "HTML file to write (required)")
	fieldList := fs.String("fields", "type,name,username,password,notes",
		"comma-separated fields: type,name,username,password,notes,url,favorite,totp,folder,organization,creation-date,modification-date,password-revision,custom-fields")
	grouping := fs.String("group", "none", "group entries by: none, type or folder")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "error: -i and -o are required")
		fs.Usage()
		return 2
	}

	fields, err := parseFields(*fieldList)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}

	mode, err := parseGrouping(*grouping)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}

	outputPath := *output
	if filepath.Ext(outputPath) != ".html" {
		outputPath += ".html"
	}

	options := ExportOptions{Fields: fields, Grouping: mode}
	if err := ConvertBitwardenToHTML(*input, outputPath, options); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}

	fmt.Println("HTML file created:", outputPath)
	return 0
}

func parseFields(list string) (ExportFields, error) {
	var fields ExportFields
	refs := map[string]*bool{
		"type":              &fields.Type,
		"name":              &fields.Name,
		"username":          &fields.Username,
		"password":          &fields.Password,
		"notes":             &fields.Notes,
		"url":               &fields.URL,
		"favorite":          &fields.Favorite,
		"totp":              &fields.TOTP,
		"folder":            &fields.Folder,
		"organization":      &fields.Organization,
		"creation-date":     &fields.CreationDate,
		"modification-date": &fields.ModificationDate,
		"password-revision": &fields.PasswordRevision,
		"custom-fields":     &fields.CustomFields,
	}

	for _, name := range strings.Split(list, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		ref, ok := refs[name]
		if !ok {
			return fields, fmt.Errorf("unknown field %q", name)
		}
		*ref = true
	}

	return fields, nil
}

func parseGrouping(name string) (GroupingMode, error) {
	switch name {
	case "none":
		return GroupNone, nil
	case "type":
		return GroupByType, nil
	case "folder":
		return GroupByFolder, nil
	}
	return GroupNone, fmt.Errorf("unknown grouping %q (expected none, type or folder)", name)
}
