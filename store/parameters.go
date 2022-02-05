package store

// Types are pattern's types usable in declaring parameter's type and its regex shape
// TODO: Allow custom types in config
var Types = map[string]string{
	"string":  "[^\\s]+",
	"integer": "[0-9]+",
	"rest":    ".*",
}
