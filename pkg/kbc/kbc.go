package kbc

type AdbOperation func() bool

type AndroidKey struct {
	Action AdbOperation
	Title  string
}

type KeyBindConfig struct {
	Ops map[string]AndroidKey
}

func NewKeyBindConfig(ops map[string]AndroidKey) *KeyBindConfig {
	return &KeyBindConfig{Ops: ops}
}

func (c *KeyBindConfig) ForEach(apply func(key string, ak AndroidKey)) {
	for k, v := range c.Ops {
		apply(k, v)
	}
}

func (c *KeyBindConfig) GetByKeyOk(key string) (AndroidKey, bool) {
	if c == nil || c.Ops == nil {
		return AndroidKey{}, false
	}
	ak, ok := c.Ops[key]
	return ak, ok
}
