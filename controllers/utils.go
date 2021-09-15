package controllers

import "fmt"

func bodyChecker(data *map[string]string, fields []string) ([]string, int) {
	var error []string

	for _, key := range fields {
		if (*data)[key] == "" {
			message := fmt.Sprintf("%v can't be empty", key)
			error = append(error, message)
		}
	}

	totalError := len(error)

	return error, totalError
}


