package main

import (
	"bytes"
	"fmt"
	"reflect"
	"unicode"
)

var typeMapping = map[string]string{
	"bool": "boolean",
	"int":  "integer",
}

type FormBuilder struct {
	Label    string
	Type     string
	Options  []Option
	Selected interface{}
}

type Option struct {
	Label string
	Value interface{}
}

func GenerateFormBuilderMap(input interface{}) map[string]FormBuilder {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	formBuilderMap := make(map[string]FormBuilder)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()
		fieldType := field.Type.Name()

		label := parseLabel(field.Name)

		if newType, ok := typeMapping[fieldType]; ok {
			fieldType = newType
		}

		// if field.Type.Kind() == reflect.Slice {
		// 	fieldType = "slice"
		// } else if field.Type.Kind() == reflect.Map {
		// 	fieldType = "map"
		// } else if field.Type.Kind() == reflect.Struct {
		// 	fieldType = "struct"
		// } else if field.Type.Kind() == reflect.Interface {
		// 	fieldType = "interface"
		// }

		formBuilder := FormBuilder{
			Label:    label,
			Type:     fieldType,
			Selected: fieldValue,
		}

		formBuilderMap[field.Name] = formBuilder
	}

	return formBuilderMap
}

func parseLabel(s string) string {
	buf := &bytes.Buffer{}
	for i, rune := range s {
		if unicode.IsUpper(rune) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(rune)
	}
	return buf.String()
}

type ExampleStruct struct {
	Username string
	Password string
	Age      int
	Active   bool
	Settings map[string]interface{}
	Roles    []string
}


type UsuarioForm struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	Enabled         bool   `json:"enabled"`
	AccountExpired  bool   `json:"account_expired" db:"account_expired"`
	AccountLocked   bool   `json:"account_locked" db:"account_locked"`
	PasswordExpired bool   `json:"password_expired" db:"password_expired"`
	ProjetoPadrao   int    `json:"projeto_padrao_id" db:"projeto_padrao_id"`
}

func main() {
	// example := ExampleStruct{
	// 	Username: "john_doe",
	// 	Password: "password123",
	// 	Age:      30,
	// 	Active:   true,
	// 	Settings: map[string]interface{}{
	// 		"theme": "dark",
	// 	},
	// 	Roles: []string{"admin", "user"},
	// }

	example := UsuarioForm{}

	formBuilderMap := GenerateFormBuilderMap(example)

	for key, fb := range formBuilderMap {
		fmt.Printf("Field: %s, Label: %s, Type: %s, Selected: %v\n", key, fb.Label, fb.Type, fb.Selected)
	}
}
