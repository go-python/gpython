// Dict and StringDict type
//
// The idea is that most dicts just have strings for keys so we use
// the simpler StringDict and promote it into a Dict when necessary

package py

var StringDictType = NewType("dict")

// String to object dictionary
//
// Used for variables etc where the keys can only be strings
type StringDict map[string]Object

// Type of this StringDict object
func (o StringDict) Type() *Type {
	return StringDictType
}

// Make a new dictionary
func NewStringDict() StringDict {
	return make(StringDict)
}
