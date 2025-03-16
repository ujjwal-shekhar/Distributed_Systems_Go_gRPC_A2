package common

// KV interface for generic key-value pairs
type KV interface {
	GetKey() interface{}
	GetValue() interface{}
	SetKey(interface{})
	SetValue(interface{})
	SetKeyVal(interface{}, interface{})
}