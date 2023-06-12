package filters

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/flosch/pongo2/v6"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Regulars
	pongo2.RegisterFilter("mapToYaml", filterMapToYaml)
	pongo2.RegisterFilter("indentYaml", filterIndentYaml)

}

func filterMapToYaml(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	defaultIndentSpacing := param.Integer()
	if defaultIndentSpacing <= 0 {
		defaultIndentSpacing = 2
	}

	text := in.Interface().(map[string]interface{})

	yamlBytes, err := encodeYaml(text, defaultIndentSpacing) // Returns "b: 2\n"
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:mapToYaml",
			OrigError: errors.New("Failed converting input to yaml. Is input a map[string]interface{}?"),
		}
	}
	return pongo2.AsSafeValue(string(yamlBytes)), nil
}

func filterIndentYaml(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var final string

	indentSpacing := param.Integer()
	if indentSpacing <= 0 {
		indentSpacing = 0
	}

	log.Debug(param)

	var pad string
	for i := 0; i < indentSpacing; i++ {
		pad += " "
	}
	scanner := bufio.NewScanner(strings.NewReader(in.String()))
	for scanner.Scan() {
		text := scanner.Text()
		log.Debug(text)
		if strings.HasPrefix(text, "- ") {
			final += fmt.Sprintln(text)
		} else {
			final += fmt.Sprintln(pad + text)
		}
	}
	return pongo2.AsSafeValue(fmt.Sprint(final)), nil
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
