package countserv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestEstimate(t *testing.T) {
	counter := NewCounter()
	counter.Insert("item1", "user1")
	counter.Insert("item1", "user2")
	counter.Insert("item1", "user3")
	estimate := counter.Estimate("item1")
	if estimate != 3 {
		t.Errorf("Expected estimate to be 3, got %d", estimate)
	}
}

func TestSerializeDeserialize(t *testing.T) {
	counter := NewCounter()
	counter.Insert("item1", "user1")
	counter.Insert("item1", "user2")
	counter.Insert("item1", "user3")
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(counter)
	if err != nil {
		t.Errorf("Error encoding counter: %s", err)
	}
	fmt.Println("json size:", buf.Len())

	var counter2 Counter
	buf2 := bytes.NewBuffer(buf.Bytes())
	json.NewDecoder(buf2).Decode(&counter2)
	estimate2 := counter2.Estimate("item1")
	if estimate2 != 3 {
		t.Errorf("Expected estimate to be 3, got %d", estimate2)
	}

	counter2.Insert("item1", "user4")
	estimate3 := counter2.Estimate("item1")
	if estimate3 != 4 {
		t.Errorf("Expected estimate to be 4, got %d", estimate2)
	}
}

func TestNonInsertedItem(t *testing.T) {
	counter := NewCounter()
	estimate := counter.Estimate("item1")
	if estimate != 0 {
		t.Errorf("Expected estimate to be 0, got %d", estimate)
	}
}

func TestEstimates(t *testing.T) {
	counter := NewCounter()
	counter.Insert("item1", "user1")
	counter.Insert("item1", "user2")
	counter.Insert("item1", "user3")
	counter.Insert("item2", "user1")
	counter.Insert("item2", "user2")
	estimates := counter.Estimates()
	expected := map[string]uint64{
		"item1": 3,
		"item2": 2,
	}
	if !reflect.DeepEqual(estimates, expected) {
		t.Errorf("Expected %v, got %v", expected, estimates)
	}
}

func TestItems(t *testing.T) {
	counter := NewCounter()
	counter.Insert("item1", "user1")
	counter.Insert("item1", "user2")
	counter.Insert("item1", "user3")
	counter.Insert("item2", "user1")
	counter.Insert("item2", "user2")
	items := counter.Items()
	expected := []string{"item1", "item2"}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected %v, got %v", expected, items)
	}
}

func TestClone(t *testing.T) {
	counter := NewCounter()
	counter.Insert("item1", "user1")
	counter.Insert("item1", "user2")
	counter.Insert("item1", "user3")
	counter.Insert("item2", "user1")
	counter.Insert("item2", "user2")
	clone := counter.Clone()
	estimates := clone.Estimates()
	expected := map[string]uint64{
		"item1": 3,
		"item2": 2,
	}
	if !reflect.DeepEqual(estimates, expected) {
		t.Errorf("Expected %v, got %v", expected, estimates)
	}
	clone.Insert("item1", "user4")
	estimates = clone.Estimates()
	expected = map[string]uint64{
		"item1": 4,
		"item2": 2,
	}
	if !reflect.DeepEqual(estimates, expected) {
		t.Errorf("Expected %v, got %v", expected, estimates)
	}
}

func BenchmarkInsertSameUsers(b *testing.B) {
	counter := NewCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Insert("item1", "user1")
	}
}

func BenchmarkInsertDifferentUsers(b *testing.B) {
	counter := NewCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Insert("item1", fmt.Sprintf("user%d", i))
	}
}

func BenchmarkEncode(b *testing.B) {
	counter := NewCounter()
	for i := 0; i < 10000; i++ {
		for j := 0; j < 100; j++ {
			counter.Insert(fmt.Sprintf("item%d", j), fmt.Sprintf("user%d", i))
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(counter)
	}
}

func BenchmarkDecode(b *testing.B) {
	counter := NewCounter()
	for i := 0; i < 10000; i++ {
		for j := 0; j < 100; j++ {
			counter.Insert(fmt.Sprintf("item%d", j), fmt.Sprintf("user%d", i))
		}
	}
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(counter)
	jsonBytes := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var counter2 Counter
		buf2 := bytes.NewBuffer(jsonBytes)
		json.NewDecoder(buf2).Decode(&counter2)
	}
}
