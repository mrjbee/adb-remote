package kbc

type AdbOperation func() bool

type AndroidKey struct {
	Аction AdbOperation
	Тitle  string
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
