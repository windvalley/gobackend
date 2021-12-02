package util

import "reflect"

// HasEntry determine whether an entry exists in a container(slice/array/map)
func HasEntry(entries interface{}, entry interface{}) bool {
	containerValue := reflect.ValueOf(entries)

	switch reflect.TypeOf(entries).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if containerValue.Index(i).Interface() == entry {
				return true
			}
		}
	case reflect.Map:
		if containerValue.MapIndex(reflect.ValueOf(entry)).IsValid() {
			return true
		}
	default:
		return false
	}

	return false
}

// StrSliceSet De-duplicate elements of a slice of string type
func StrSliceSet(slice []string) []string {
	set := make([]string, 0)
	tempMap := make(map[string]bool, len(slice))
	for _, v := range slice {
		if !tempMap[v] {
			set = append(set, v)
			tempMap[v] = true
		}
	}

	return set
}
