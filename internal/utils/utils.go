package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	rawyaml, err := getFileContents(file)
	if err != nil {
		return nil, err
	}

	// Map to store the parsed YAML data
	// var dataDocs []map[string]interface{}
	var data map[string]interface{}

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(rawyaml, &data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

func GetYamlContentFromString(content string) (map[string]interface{}, error) {
	// Map to store the parsed YAML data
	var dataDocs []map[string]interface{}
	var data map[string]interface{}

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
		log.Debug(template)
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

func GetInput(file string, pipeFlag bool, printInputFlag bool) (map[string]interface{}, error) {
	var yamlData map[string]interface{}
	// var inputDataOriginal map[string]interface{}
	var inputDataFromIncludes map[string]interface{}
	var inputData map[string]interface{}
	var includeFiles []interface{}
	var content string
	var err error

	if pipeFlag {
		// Get text
		log.Debug("Getting input from pipe")
		content, err = GetPipedContent()
		if err != nil {
			return nil, err
		}
	} else {
		// Get text
		log.Debug("Getting input from file")
		content, err = getFileContentsAsString(file)
		if err != nil {
			return nil, err
		}
	}
	// Read the file
	yamlData, err = GetYamlContentFromString(content)
	if err != nil {
		return nil, err
	}
	includeFiles, ok := yamlData["includes"].([]interface{})
	if !ok {
		includeFiles = make([]interface{}, 0, 0)
	}
	inputDataFromIncludes, err = GetIncludes(includeFiles)
	delete(yamlData, "includes")
	inputData = MergeMaps(inputDataFromIncludes, yamlData)
	inputText, inputData, err := RenderInput(inputData)
	if printInputFlag {
		fmt.Printf(inputText)
		os.Exit(0)
	} else {
		log.Debug("--- InputData Start --- \n", inputText)
		log.Debug("--- InputData End --- ")
	}
	return inputData, nil
}

func GetTemplate(file string, pipeFlag bool) (string, error) {
	var content string
	var err error

	if pipeFlag {
		// Read the file
		log.Debug("Getting template from pipe")
		content, err = GetPipedContent()
		if err != nil {
			return "", err
		}
	} else {
		// Read the file
		log.Debug("Getting template from file")
		content, err = getFileContentsAsString(file)
		if err != nil {
			return "", err
		}
	}

	return content, nil

}

func getFileContents(file string) ([]byte, error) {
	log.Debug("Reading file: ", file)
	contentBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return contentBytes, nil
}

func getFileContentsAsString(file string) (string, error) {
	contentBytes, err := getFileContents(file)
	if err != nil {
		return "", err
	}
	content := string(contentBytes)
	return content, nil
}

func RenderTemplate(template string, input map[string]any) (string, error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		panic(err)
	}

	// Now you can render the template with the given
	// pongo2.Context how often you want to.
	context := pongo2.Context{
		"printMapAsYamlString": func(a map[string]interface{}, defaultIndentSpacing int) string {

			yamlBytes, err := encodeYaml(a, defaultIndentSpacing) // Returns "b: 2\n"
			if err != nil {
				return fmt.Sprintln("error: help!")
			}

			return string(yamlBytes)
		},
		"indentYaml": func(text string, spaces int) string {
			var final string

			var pad string
			for i := 0; i < spaces; i++ {
				pad += " "
			}
			scanner := bufio.NewScanner(strings.NewReader(text))
			for scanner.Scan() {
				text := scanner.Text()
				log.Debug(text)
				if strings.HasPrefix(text, "- ") {
					final += fmt.Sprintln(text)
				} else {
					final += fmt.Sprintln(pad + text)
				}
			}
			return fmt.Sprint(final)

		},
	}
	context.Update(input)
	out, err := tpl.Execute(context)
	if err != nil {
		// panic(err)
		return "", err
	}

	return out, nil

}

func GetIncludes(inputFiles []interface{}) (map[string]interface{}, error) {
	var tmpData map[string]interface{}
	var file string
	for _, include := range inputFiles {
		log.Debug("Processing includes...")
		include, ok := include.(map[string]interface{})
		if !ok {
			log.Debug("not ok")
		}
		filePath := include["path"].(string)
		fileBase := include["base"].(string)
		if fileBase == "currentDirectory" {
			currentDir, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			file = filepath.Join(currentDir, filePath)
		} else if fileBase == "repositoryRoot" {
			repoRoot, err := getRepoRoot()
			if err != nil {
				return nil, err
			}
			file = filepath.Join(strings.TrimSpace(repoRoot), filePath)
		} else {
			log.Error("[base] value must be 'repositoryRoot' or 'currentDirectory'")
			os.Exit(1)
		}
		data_chunk, err := GetYamlContentFromFile(file)
		if err != nil {
			return nil, err
		}
		tmpData = MergeMaps(data_chunk, tmpData)
		// TODO: Possibly support includes within includes
	}
	return tmpData, nil
}

func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			// fmt.Println(k + " is a map")
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

func RenderInput(input map[string]interface{}) (string, map[string]interface{}, error) {
	tmp, err := yaml.Marshal(input)
	if err != nil {
		return "", nil, err
	}
	inputText := string(tmp)
	inputData, err := GetYamlContentFromString(inputText)
	count := 1
	for isTemplated(inputText) {
		log.Debug("Processing input as template... " + fmt.Sprint(count) + " pass")
		inputText, err = RenderTemplate(inputText, inputData)
		inputData, err = GetYamlContentFromString(inputText)
		count++
	}
	return inputText, inputData, nil
}

func getRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func isTemplated(template string) bool {
	r := regexp.MustCompile(`{{(.*?)}}`)
	results := r.FindAllStringSubmatch(template, -1)
	if len(results) > 0 {
		return true
	} else {
		return false
	}
}

func encodeYaml(r map[string]interface{}, spaces int) ([]byte, error) {
	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(spaces)

	if err := e.Encode(&r); err != nil {
		return nil, err
	}
	return b.Bytes(), nil

}
