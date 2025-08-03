package loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Validatable interface {
	Validate() error
}

func LoadConfig(configPath string, config Validatable) error {
	// Check that provided interface is a pointer to a struct
	interfaceType := reflect.TypeOf(config)
	if interfaceType.Kind() != reflect.Ptr || interfaceType.Elem().Kind() != reflect.Struct {
		panic("provided interface must be a pointer to struct")
	}

	// Loading yaml config file
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error while opening %s: %s", configPath, err.Error())
	}

	// Try to load file as yaml
	err = yaml.Unmarshal(configFile, config)

	// If error, try to load it as json
	if err != nil {
		err = json.Unmarshal(configFile, config)
		if err != nil {
			return errors.New("config parsing error : " + err.Error())
		}
	}

	// Overwrite config with Env Var if Needed
	if err := parseEnvVars(config); err != nil {
		return errors.New("error while converting environment variables: " + err.Error())
	}

	return nil
}

func parseEnvVars(config Validatable) error {
	// Retrieve struct fields infos
	fieldsInfo := reflect.VisibleFields(reflect.TypeOf(config).Elem())
	configElements := reflect.ValueOf(config).Elem()

	// For each field
	for _, field := range fieldsInfo {

		// Look for associated Env Var
		value := os.Getenv(field.Tag.Get("env"))

		// If env var is defined, overwrite associated field in config
		if value != "" {

			if field.Type.Kind() == reflect.Ptr { // Manage pointers on simple types
				switch field.Type.Elem().Kind() {
				case reflect.Int:
					if intValue, err := strconv.Atoi(value); err == nil {
						configElements.FieldByName(field.Name).Set(reflect.ValueOf(&intValue))
					} else {
						return fmt.Errorf("unable to transform %s as integer pointer", value)
					}
				case reflect.String:
					stringValue := evaluateVars(value)
					configElements.FieldByName(field.Name).Set(reflect.ValueOf(&stringValue))
				case reflect.Bool:
					if boolValue, err := strconv.ParseBool(value); err == nil {
						configElements.FieldByName(field.Name).Set(reflect.ValueOf(&boolValue))
					} else {
						return fmt.Errorf("unable to transform %s as boolean pointer", value)
					}
				}

			} else { // Else, manage simple types
				switch field.Type.Kind() {
				case reflect.Int:
					if intValue, err := strconv.Atoi(value); err == nil {
						configElements.FieldByName(field.Name).Set(reflect.ValueOf(intValue))
					} else {
						return fmt.Errorf("unable to transform %s as integer", value)
					}
				case reflect.String:
					configElements.FieldByName(field.Name).SetString(evaluateVars(value))
				case reflect.Bool:
					if boolValue, err := strconv.ParseBool(value); err == nil {
						configElements.FieldByName(field.Name).SetBool(boolValue)
					} else {
						return fmt.Errorf("unable to transform %s as boolean", value)
					}
				}
			}
		}
	}
	return nil
}

/*
Used to evaluate variables in a input string.
For example, if val is "TEST_${MYVAR}" and export MYVAR="test", the result is "TEST_test"
*/
func evaluateVars(val string) string {

	// This regexp is used to detect a variable formed like ${VAR} in input string.
	// It captures the variable name in 'matches[1]'
	variableRegexp := regexp.MustCompile(`\$\{([a-zA-Z0-9_-]*)\}`)

	// Iterate until no variables left in val
	for {
		// Search for envvar in string
		matches := variableRegexp.FindStringSubmatch(val)

		// If no variable is found, return result
		if len(matches) == 0 {
			return val
		}

		detectedVariableName := matches[1]

		// Look for env var value
		envValue := os.Getenv(detectedVariableName)

		// replace var name by value
		val = strings.Replace(val, fmt.Sprintf("${%s}", detectedVariableName), envValue, 1)
	}
}
