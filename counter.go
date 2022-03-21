package countserv

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"

	"github.com/axiomhq/hyperloglog"
)

type Entry = *hyperloglog.Sketch

type Counter struct {
	entries map[string]Entry
	mutex   sync.Mutex
}

func NewCounter() *Counter {
	return &Counter{
		entries: make(map[string]Entry),
	}
}

func (c *Counter) Insert(item string, user string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.entries[item]; !ok {
		c.entries[item] = hyperloglog.New()
	}
	return c.entries[item].Insert([]byte(user))
}

func (c *Counter) Estimate(item string) uint64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.entries[item]; !ok {
		return 0
	}
	return c.entries[item].Estimate()
}

func (c *Counter) Estimates() map[string]uint64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	estimates := make(map[string]uint64)
	for item, entry := range c.entries {
		estimates[item] = entry.Estimate()
	}
	return estimates
}

func (c *Counter) Items() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	items := make([]string, 0, len(c.entries))
	for item := range c.entries {
		items = append(items, item)
	}
	return items
}

func (c *Counter) Clone() *Counter {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	clone := NewCounter()
	for k, v := range c.entries {
		clone.entries[k] = v.Clone()
	}
	return clone
}

func (c *Counter) MarshalJSON() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	encodedEntries := make(map[string]string)
	for k, v := range c.entries {
		b, err := v.MarshalBinary()
		if err != nil {
			return nil, err
		}
		encodedEntries[k] = base64.StdEncoding.EncodeToString(b)
	}

	return json.Marshal(&struct {
		Version int               `json:"version"`
		Entries map[string]string `json:"entries"`
	}{
		Version: 1,
		Entries: encodedEntries,
	})
}

func (c *Counter) UnmarshalJSON(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var decoded struct {
		Version int               `json:"version"`
		Entries map[string]string `json:"entries"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	if decoded.Version != 1 {
		return errors.New("unsupported version")
	}
	c.entries = make(map[string]Entry)
	for k, v := range decoded.Entries {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		c.entries[k] = hyperloglog.New()
		if err := c.entries[k].UnmarshalBinary(b); err != nil {
			return err
		}
	}

	return nil
}
