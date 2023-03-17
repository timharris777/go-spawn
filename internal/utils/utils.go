package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func PasswordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

func UNUSED(x ...interface{}) {}

var dataDocs []map[string]interface{}

func GetYamlContentFromFile(file string) (map[string]interface{}, error) {
	// Read the file
	log.Debug("Reading file: ", file)
	rawyaml, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Map to store the parsed YAML data
	var dataDocs []map[string]interface{}
	var data map[string]interface{}

	// Unmarshal the YAML data into the struct
	// err = yaml.Unmarshal(rawyaml, &data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	if err := UnmarshalAllYamlDocs(rawyaml, &dataDocs); err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, doc := range dataDocs {
		variables, ok := doc["variables"].(map[string]interface{})
		if ok {
			log.Debug("Found variables: ", variables)
			return doc, nil
		}
		data = doc
	}
	return data, nil
}

func GetYamlContentFromString(content string) (map[string]interface{}, error) {
	// Map to store the parsed YAML data
	var dataDocs []map[string]interface{}
	var data map[string]interface{}

	// // Unmarshal the YAML data into the struct
	// err := yaml.Unmarshal([]byte(content), &data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	if err := UnmarshalAllYamlDocs([]byte(content), &dataDocs); err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, doc := range dataDocs {
		variables, ok := doc["variables"].(map[string]interface{})
		if ok {
			log.Debug("Found variables: ", variables)
			return doc, nil
		}
		data = doc
	}
	return data, nil
}

func UnmarshalAllYamlDocs(in []byte, out *[]map[string]interface{}) error {
	r := bytes.NewReader(in)
	decoder := yaml.NewDecoder(r)
	for {
		var bo map[string]interface{}

		if err := decoder.Decode(&bo); err != nil {
			// Break when there are no more documents to decode
			if err != io.EOF {
				return err
			}
			break
		}
		*out = append(*out, bo)
	}
	return nil
}

func GetPipedContent() (string, error) {
	var template string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			template = template + fmt.Sprintln(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	} else {
		// TODO: Clean this up
		fmt.Print("Enter your name: ")

		var content string
		fmt.Scanf("%s", &content)
		fmt.Printf(content)
	}
	return template, nil
}

func FlagValidation(inputFilePath string, inputPipe bool, templatePath string, templatePipe bool, outputPath string) error {
	// is input provided (must be provided or piped)
	// is template provided (must be provided or piped in)
	// if template is dir is output dir provided
	// if output is provided and template is file output must be file
	// if output is provided and template is dir output must be dir
	return nil
}

func GetInput(file string, pipeFlag bool) (map[string]interface{}, error) {
	var yamlData map[string]interface{}
	var inputDataOriginal map[string]interface{}
	var inputDataFromIncludes map[string]interface{}
	var inputDataFinal map[string]interface{}
	var includeFiles []any
	var content string
	var err error

	if pipeFlag {
		// Read the file
		content, err = GetPipedContent()
		if err != nil {
			return nil, err
		}
		// Read the file
		yamlData, err = GetYamlContentFromString(content)
		if err != nil {
			return nil, err
		}
	} else {
		// Read the file
		yamlData, err = GetYamlContentFromFile(file)
		if err != nil {
			return nil, err
		}
	}
	log.Debug("yamlData: ", yamlData)
	// includeFiles, ok := yamlData["includes"].([]any)
	// if !ok {
	includeFiles, ok := yamlData["variableFiles"].([]any)
	if !ok {
		includeFiles = make([]any, 0, 0)
	}
	// }
	// inputDataOriginal, ok = yamlData["data"].(map[string]interface{})
	// if !ok {
	inputDataOriginal, ok = yamlData["variables"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Cannot find ['variables'] key in root of input yaml file.")
	}
	// }
	inputDataFromIncludes, err = GetIncludes(includeFiles)
	inputDataFinal = MergeMaps(inputDataFromIncludes, inputDataOriginal)
	return inputDataFinal, nil
}

func GetTemplate(file string, pipeFlag bool) (string, error) {
	var contentRaw []byte
	var content string
	var err error

	if pipeFlag {
		// Read the file
		content, err = GetPipedContent()
		if err != nil {
			return "", err
		}
	} else {
		// Read the file
		contentRaw, err = os.ReadFile(file)
		content = string(contentRaw)
		if err != nil {
			return "", err
		}
	}

	return content, nil

}

func RenderTemplate(template string, input map[string]any) (string, error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		panic(err)
	}

	// Now you can render the template with the given
	// pongo2.Context how often you want to.
	context := make(pongo2.Context)
	context.Update(input)
	out, err := tpl.Execute(context)
	if err != nil {
		// panic(err)
		return "", err
	}

	return out, nil

}

func GetIncludes(inputFiles []any) (map[string]interface{}, error) {
	var tmpData map[string]interface{}
	for fileIndex := range inputFiles {
		file := inputFiles[fileIndex].(string)
		log.Debug("Getting includes for input: ", file)
		data_chunk, err := GetYamlContentFromFile(file)
		log.Debug("Include file contents: ", data_chunk)
		if err != nil {
			return nil, err
		}
		tmpData = MergeMaps(data_chunk["variables"].(map[string]interface{}), tmpData)
		// TODO: Possibly support includes within includes
	}
	log.Debug("Returning merged map of includes: ", tmpData)
	return tmpData, nil
}

func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			fmt.Println(k + " is a map")
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
