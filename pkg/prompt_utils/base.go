package prompt_utils

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// StructToXMLString converts a struct into an XML-like string format with snake_cased tags
func StructToXMLString(v interface{}) (string, error) {
	var result strings.Builder
	val := reflect.ValueOf(v)
	typ := val.Type()

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("input must be a struct, got %v", val.Kind())
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Handle embedded structs
		if fieldType.Anonymous {
			if field.Kind() == reflect.Struct {
				embedded, err := StructToXMLString(field.Interface())
				if err != nil {
					return "", err
				}
				result.WriteString(embedded)
				continue
			}
		}

		// Convert field name to snake_case
		tagName := toSnakeCase(fieldType.Name)

		// Format the field value
		var fieldStr string
		switch {
		case field.Kind() == reflect.Struct:
			embedded, err := StructToXMLString(field.Interface())
			if err != nil {
				return "", err
			}
			fieldStr = embedded
		case field.Kind() == reflect.Interface && !field.IsNil():
			// Try to call GetName() if it exists
			if nameMethod := field.MethodByName("GetName"); nameMethod.IsValid() {
				result := nameMethod.Call(nil)
				if len(result) > 0 {
					fieldStr = fmt.Sprintf("%v", result[0])
				} else {
					fieldStr = field.Elem().Type().String()
				}
			} else {
				// Fallback to type name
				fieldStr = field.Elem().Type().String()
			}
		case field.Kind() == reflect.Slice || field.Kind() == reflect.Array:
			items := make([]string, field.Len())
			for j := 0; j < field.Len(); j++ {
				items[j] = fmt.Sprintf("%v", field.Index(j))
			}
			fieldStr = strings.Join(items, "\n")
		default:
			fieldStr = fmt.Sprintf("%v", field)
		}

		// For embedded struct fields, don't wrap them in tags
		if field.Kind() == reflect.Struct {
			result.WriteString(fieldStr)
		} else {
			result.WriteString(fmt.Sprintf("<%s>\n%s\n</%s>\n\n", tagName, fieldStr, tagName))
		}
	}

	return result.String(), nil
}

// toSnakeCase converts a string from CamelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
