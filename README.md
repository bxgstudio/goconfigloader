# GoConfigLoader
<img src="./assets/CL.png" alt="Texte alternatif" width="200"/>
<br/>
This repository is a golang implementation of ConfigLoader, a project that provides a simple way to load application configuration from file and to override parameters with associated env vars if any.
An evaluate of env vars is done after configuration loading which allows you to customize properly environment variables on tricky scenarii (can happen in k8s contexts to concatenate variables from POD_NAME and add a suffix for example).

Here you can find the python version of this project: [Python repository](https://github.com/bxgstudio/pyconfigloader)

---

Manage following configuration file formats: 
- yaml
- json

---

Manage following types of env vars:
- ENV=string
- ENV=${OTHER_ENV}_{ANOTHER_ONE}_string

---

# Installation

go get github.com/bxgstudio/appconfigloader@1.0.0

# Mecanism

LoadConfig function expects to receive pointer to a struct that implements a Validate function. This feature is designed to force the product that uses this package to do a logical check of loaded config in order to build robust code.

Types that are managed by this package are the following:
- int
- string
- bool
- *int
- *string
- *bool

Pointers on simple types can be used when we want to force a field to be explicitly provided in configuration.
The user will be capable to check his config field is provided or not. 

In Golang, we use the strength of structure tags to use a "env" tag which allows the user to define a environment variable with different name that the structure field of json/yaml tag. It is usefull when you execute your application in a shared environment where environment variables need to be really specific to your app to avoid conflict with other variables.

# Usage

Example below:
```go
package main

import (
    "gitlab.com/gck-prod/configloader/loader"
)

type AppConfig struct {
    Host *string `yaml:"host_field" env:"MY_APP_HOST_ENV_VAR"`
    Port *int `yaml:"port_field" env:"MY_APP_PORT_ENV_VAR"`
    Debug bool `yaml:"debug_field"`
} 

func (appConfig *AppConfig) Validate() error {
    if appConfig.Host == nil {
        return errors.New("'host' should be provided")
    }
    if appConfig.Port == nil {
        return errors.New("'port' should be provided")
    }
}

func main() {
    if err := loader.LoadConfig("/path/to/config.yaml", &AppConfig{}); err != nil {
        log.Fatalf("error while loading config: %s", err.Error())
    }

    if err := appConfig.Validate(); err != nil {
        log.Fatalf("error on config validation: %s", err.Error())
    }

    // ...
}
```

A full example of usage is in the loader_test.go file.

# Authors

- Etienne Galecki - [@galecki](https://github.com/egck)
- Antoine Breton - [@breton](https://github.com/antbreton)