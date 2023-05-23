package storage

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func PersistToLocalStorage(value any) {
	v, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}

	w := js.Global().Get("window")
	ls := w.Get("localStorage")
	ls.Call("setItem", "todos", string(v))
}

func GetFromLocalStorage() []interface{} {
	todos := []interface{}{}
	w := js.Global().Get("window")
	ls := w.Get("localStorage")
	r := ls.Call("getItem", "todos")
	fmt.Println(r.String())

	err := json.Unmarshal([]byte(r.String()), &todos)
	if err != nil {
		fmt.Println(err)
	}

	return todos
}
