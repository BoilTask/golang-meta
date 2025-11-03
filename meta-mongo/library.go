package metamongo

import "go.mongodb.org/mongo-driver/bson"

func StructToMapInclude(input any, includeKeys ...string) bson.M {
	m := bson.M{}
	data, err := bson.Marshal(input)
	if err != nil {
		return m
	}
	if err := bson.Unmarshal(data, &m); err != nil {
		return m
	}
	result := bson.M{}
	for _, key := range includeKeys {
		if value, exists := m[key]; exists {
			result[key] = value
		}
	}
	return result
}

func StructToMapExclude(input any, excludeKeys ...string) bson.M {
	m := bson.M{}
	data, err := bson.Marshal(input)
	if err != nil {
		return m
	}
	if err := bson.Unmarshal(data, &m); err != nil {
		return m
	}
	for _, key := range excludeKeys {
		delete(m, key)
	}
	return m
}
