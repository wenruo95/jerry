/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : yaml.go
*   coder: zemanzeng
*   date : 2021-09-29 17:10:00
*   desc : read yaml
*
================================================================*/

package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYaml(file string, setting interface{}) error {
	buff, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(buff, setting); err != nil {
		return err
	}
	return nil
}
