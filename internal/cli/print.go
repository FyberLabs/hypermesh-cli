package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

func (r *run) printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (r *run) printRawJSON(raw json.RawMessage) error {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	return r.printJSON(v)
}

func unwrapItems(raw json.RawMessage) []map[string]any {
	var arr []map[string]any
	if json.Unmarshal(raw, &arr) == nil {
		return arr
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) != nil {
		return nil
	}
	for _, key := range []string{"items", "catalog", "classes", "leases", "data", "results"} {
		if v, ok := obj[key]; ok {
			b, err := json.Marshal(v)
			if err != nil {
				continue
			}
			var inner []map[string]any
			if json.Unmarshal(b, &inner) == nil {
				return inner
			}
		}
	}
	if _, ok := obj["id"]; ok {
		return []map[string]any{obj}
	}
	return nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch t := v.(type) {
			case string:
				if t != "" {
					return t
				}
			}
		}
	}
	return ""
}

func printIDList(raw json.RawMessage) error {
	items := unwrapItems(raw)
	if items == nil {
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	}
	for _, item := range items {
		id := firstString(item, "id", "catalog_id", "class_id")
		label := firstString(item, "name", "title", "display_name")
		if id == "" {
			continue
		}
		if label != "" && label != id {
			fmt.Printf("%s\t%s\n", id, label)
			continue
		}
		fmt.Println(id)
	}
	return nil
}

func findByID(raw json.RawMessage, id string) (map[string]any, error) {
	for _, item := range unwrapItems(raw) {
		got := firstString(item, "id", "catalog_id", "class_id")
		if got == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", id)
}
