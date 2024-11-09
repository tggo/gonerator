package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	outputDirectory = "internal/"
	moduleName      = "tggo-gonerator"
	domainNameRaw   string
	overwriteAll    bool
	partGenerate    = []string{
		"interfaces",
		"handlers",
		"handler-register",
		"views-response",
		"views",
		"model",
		"model-error",
		"model-filter",
		"repository",
		"services",
	}
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Generate the needed files to quickly create domain in internal folder",
	Long: `Make a new domain in internal folder with the needed files:
- handlers
- handlers
- views
- models
- repository
- services
- mocks
- activities
- workflows
`,
	Run: func(cmd *cobra.Command, _ []string) {
		// clean and add slash at the end
		outputDirectory = strings.Trim(outputDirectory, " \r\n\t/") + "/"

		// we receive the first argument as domain name
		domainName := cases.Lower(language.English).String(domainNameRaw)
		pluralizeClient := pluralize.NewClient()
		domainNamePlural := pluralizeClient.Plural(domainName)

		for i := range partGenerate {
			partGenerate[i] = strings.Trim(strings.ToLower(partGenerate[i]), " \r\n\t")
		}

		fmt.Printf("\n %s %s %s\n", lightGreen("Bootstrap called for domain:"), green(domainNamePlural), cmd.Version)

		// create domain folders
		folders := []string{
			"handlers",
			"views",
			"model",
			"repository",
			"services",
			"mocks",
			"activities",
			"workflows",
		}

		for i := range folders {
			msg := fmt.Sprintf("%s%s/%s", outputDirectory, yellow(domainNamePlural), green(folders[i]))
			fmt.Printf("\tCreating %s\n", msg)

			err := os.MkdirAll(fmt.Sprintf("%s%s/%s", outputDirectory, domainNamePlural, folders[i]), os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}

		data := map[string]string{
			"Module":                     moduleName,
			"OutputDirectory":            outputDirectory,
			"DomainName":                 domainName,
			"DomainNamePlural":           domainNamePlural,
			"DomainNameCapitalize":       cases.Title(language.Und).String(domainName),
			"DomainNamePluralCapitalize": cases.Title(language.Und).String(domainNamePlural),
		}

		createdFiles := []string{}

		fmt.Printf("\n\tGenerating parts %s\n", green(strings.Join(partGenerate, ", ")))

		// write interface file in prefix/domainName
		interfaceFile := fmt.Sprintf("%s%s/interfaces.go", outputDirectory, domainNamePlural)
		upsertFile("interfaces", interfaceFile, func() {
			writeFile(interfaceFile, "cmd/cli/cmd/templates/bootstrap/interfaces.go.tmpl", data)
			createdFiles = append(createdFiles, interfaceFile)
		})

		// create a model file
		modelFile := fmt.Sprintf("%s%s/model/%s.go", outputDirectory, domainNamePlural, domainName)
		upsertFile("model", modelFile, func() {
			writeFile(modelFile, "cmd/cli/cmd/templates/bootstrap/model.go.tmpl", data)
			createdFiles = append(createdFiles, modelFile)
		})

		// create a errors file
		errorsFile := fmt.Sprintf("%s%s/model/errors.go", outputDirectory, domainNamePlural)
		upsertFile("model-error", errorsFile, func() {
			writeFile(errorsFile, "cmd/cli/cmd/templates/bootstrap/model_errors.go.tmpl", data)
			createdFiles = append(createdFiles, errorsFile)
		})

		// create a model filter file
		modelFilterFile := fmt.Sprintf("%s%s/model/filter.go", outputDirectory, domainNamePlural)
		upsertFile("model-filter", modelFilterFile, func() {
			writeFile(modelFilterFile, "cmd/cli/cmd/templates/bootstrap/model_filter.go.tmpl", data)
			createdFiles = append(createdFiles, modelFilterFile)
		})

		// create a response file for htmx
		viewResponseFile := fmt.Sprintf("%s%s/views/responses.go", outputDirectory, domainNamePlural)
		upsertFile("views-response", viewResponseFile, func() {
			writeFile(viewResponseFile, "cmd/cli/cmd/templates/bootstrap/view_response.go.tmpl", data)
			createdFiles = append(createdFiles, viewResponseFile)
		})

		// create a repository file
		repositoryFile := fmt.Sprintf("%s%s/repository/%s.go", outputDirectory, domainNamePlural, domainName)
		upsertFile("repository", repositoryFile, func() {
			writeFile(repositoryFile, "cmd/cli/cmd/templates/bootstrap/repository_mongodb.go.tmpl", data)
			createdFiles = append(createdFiles, repositoryFile)
		})

		// create a handler's register
		handlerRegisterFile := fmt.Sprintf("%s%s/handlers/register.go", outputDirectory, domainNamePlural)
		upsertFile("handler-register", handlerRegisterFile, func() {
			writeFile(handlerRegisterFile, "cmd/cli/cmd/templates/bootstrap/handler_register.go.tmpl", data)
			createdFiles = append(createdFiles, handlerRegisterFile)
		})

		// create list handler
		handlerListFile := fmt.Sprintf("%s%s/handlers/%s_list.go", outputDirectory, domainNamePlural, domainNamePlural)
		upsertFile("handlers", handlerListFile, func() {
			writeFile(handlerListFile, "cmd/cli/cmd/templates/bootstrap/handler_list.go.tmpl", data)
			createdFiles = append(createdFiles, handlerListFile)
		})

		// create domain templ page
		viewsMainFile := fmt.Sprintf("%s%s/views/%s.templ", outputDirectory, domainNamePlural, domainNamePlural)
		upsertFile("views", viewsMainFile, func() {
			writeFile(viewsMainFile, "cmd/cli/cmd/templates/bootstrap/views_main.templ.tmpl", data)
			createdFiles = append(createdFiles, viewsMainFile)
		})

		// create domain service
		serviceFile := fmt.Sprintf("%s%s/services/service.go", outputDirectory, domainNamePlural)
		upsertFile("services", serviceFile, func() {
			writeFile(serviceFile, "cmd/cli/cmd/templates/bootstrap/service.go.tmpl", data)
			createdFiles = append(createdFiles, serviceFile)
		})

		fmt.Printf("\tgenerated %d files\n", len(createdFiles))
	},
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)

	// add prefix flag
	bootstrapCmd.Flags().StringVarP(&outputDirectory, "output", "o", "internal", "Output directory")

	// set module name
	bootstrapCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Module name")

	// set domain name
	bootstrapCmd.Flags().StringVarP(&domainNameRaw, "domain", "d", "", "Domain name, must be singular form, e.g: user, card, account")

	// mark domain name as required
	err := bootstrapCmd.MarkFlagRequired("domain")
	if err != nil {
		fmt.Println(err)
	}

	// add part flag
	bootstrapCmd.Flags().StringSliceVarP(&partGenerate, "part", "p", partGenerate, "Part to generate, e.g: handlers, views, model, repository, services, mocks, activities, workflows")
}

func upsertFile(partName string, filename string, toRun func()) {
	if !ifNeedGeneratePart(partName) {
		return
	}

	fmt.Printf("\tCreating %s\n", yellow(filename))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		toRun()
	}

	if overwriteAll {
		toRun()

		return
	}

	if _, err := os.Stat(filename); err == nil {
		fmt.Printf(lightRed("\tAlready exists, overwrite it? [y/n/Y]: ")) // nolint

		var overwrite string

		_, err = fmt.Scanln(&overwrite)
		if err != nil {
			return
		}

		if overwrite == "y" {
			fmt.Println(red("\tOverwriting file"))
			toRun()
		} else if overwrite == "Y" {
			fmt.Println(red("\tOverwriting ALL files"))

			overwriteAll = true

			toRun()
		}

		return
	}
}

func ifNeedGeneratePart(part string) bool {
	for _, p := range partGenerate {
		if p == part {
			return true
		}
	}

	return false
}

func writeFile(fileName, templateFileName string, data map[string]string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	var buf bytes.Buffer

	// for template use go text/template
	// parse the template
	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		panic(fmt.Errorf("parse: %w", err))
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(fmt.Errorf("execute: %w", err))
	}

	if strings.Contains(templateFileName, ".templ.") {
		_, err = f.Write(buf.Bytes())
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	// run go fmt for our files
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(fmt.Errorf("format: %w", err))
	}

	_, err = f.Write(p)
	if err != nil {
		fmt.Println(err)
	}
}
