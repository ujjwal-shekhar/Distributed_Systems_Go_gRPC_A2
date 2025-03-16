package task1

// INPUT
type Key_in struct {
	Item string
}
type Value_in struct {
	Item string
}
type KV_in struct {
	Key   Key_in
	Value Value_in
}

func (kv KV_in) GetKey() interface{} {
	return kv.Key
}
func (kv KV_in) GetValue() interface{} {
	return kv.Value
}
func (kv *KV_in) SetKey(key interface{}) {
	kv.Key = key.(Key_in)
}
func (kv *KV_in) SetValue(value interface{}) {
	kv.Value = value.(Value_in)
}
func (kv *KV_in) SetKeyVal(key interface{}, value interface{}) {
	kv.Key = key.(Key_in)
	kv.Value = value.(Value_in)
}

// INTERMEDIATE
type Key_intermediate struct {
	Item string
}
type Value_intermediate struct {
	Item int32
}
type KV_intermediate struct {
	Key   Key_intermediate
	Value Value_intermediate
}

func (kv KV_intermediate) GetKey() interface{} {
	return kv.Key.Item
}
func (kv KV_intermediate) GetValue() interface{} {
	return kv.Value.Item
}
func (kv *KV_intermediate) SetKey(key interface{}) {
	kv.Key = key.(Key_intermediate)
}
func (kv *KV_intermediate) SetValue(value interface{}) {
	kv.Value = value.(Value_intermediate)
}
func (kv *KV_intermediate) SetKeyVal(key interface{}, value interface{}) {
	kv.Key = key.(Key_intermediate)
	kv.Value = value.(Value_intermediate)
}

// OUTPUT
type Key_out struct {
	Item string
}
type Value_out struct {
	Item int32
}

type KV_out struct {
	Key   Key_out
	Value Value_out
}

func (kv KV_out) GetKey() interface{} {
	return kv.Key
}
func (kv KV_out) GetValue() interface{} {
	return kv.Value
}
func (kv *KV_out) SetKey(key interface{}) {
	kv.Key = key.(Key_out)
}
func (kv *KV_out) SetValue(value interface{}) {
	kv.Value = value.(Value_out)
}
func (kv *KV_out) SetKeyVal(key interface{}, value interface{}) {
	kv.Key = key.(Key_out)
	kv.Value = value.(Value_out)
}