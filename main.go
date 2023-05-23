package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type Todo struct {
	Name     string
	Complete bool
}

var (
	document  = js.Global().Get("document")
	container = document.Call("getElementById", "mycont")

	todos = make([]*Todo, 0)
)

func renderTodos() {
	container.Set("innerHTML", "")
	renderForm()

	heading := document.Call("createElement", "h1")
	heading.Set("innerHTML", "Todos")
	heading.Set("classList", "text-center text-3xl font-bold")
	container.Call("appendChild", heading)
	allCompleted := true

	for i := range todos {
		index := i
		if !todos[i].Complete {
			allCompleted = false
		}

		listItem := document.Call("createElement", "div")
		button := document.Call("createElement", "button")
		if todos[i].Complete {
			button.Set("innerHTML", "Done")
		} else {
			button.Set("innerHTML", "Undone")
		}
		button.Set("classList", "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded")

		delButton := document.Call("createElement", "button")
		delButton.Set("innerHTML", "X")
		delButton.Set("classList", "bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded")

		listItem.Set("innerHTML", todos[index].Name)
		cl := "text-black font-bold py-2 px-4 rounded"
		if todos[index].Complete {
			cl = cl + " line-through"
		}
		listItem.Set("classList", cl)
		container.Call("appendChild", listItem)
		container.Call("appendChild", delButton)
		container.Call("appendChild", button)

		button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			todos[index].Complete = !todos[index].Complete
			renderTodos()
			persistToLocalStorage()
			return nil
		}))

		delButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			todos = append(todos[:index], todos[index+1:]...)
			renderTodos()
			persistToLocalStorage()
			return nil
		}))
	}

	if allCompleted {
		completed := document.Call("createElement", "div")
		completed.Set("innerHTML", "Completed")
		completed.Set("classList", "text-center text-3xl font-bold")
		container.Call("appendChild", completed)
	}
}

func persistToLocalStorage() {
	v, err := json.Marshal(todos)
	if err != nil {
		fmt.Println(err)
	}

	w := js.Global().Get("window")
	ls := w.Get("localStorage")
	ls.Call("setItem", "todos", string(v))
}

func getFromLocalStorage() {
	w := js.Global().Get("window")
	ls := w.Get("localStorage")
	r := ls.Call("getItem", "todos")
	fmt.Println(r.String())

	err := json.Unmarshal([]byte(r.String()), &todos)
	if err != nil {
		fmt.Println(err)
	}
}

func renderForm() {
	parent := document.Call("createElement", "div")
	parent.Set("classList", "flex space-x-5 my-5")
	container.Call("appendChild", parent)

	label := document.Call("createElement", "label")
	label.Set("innerHTML", "New Todo")
	parent.Call("appendChild", label)

	input := document.Call("createElement", "input")
	input.Set("classList", "border border-gray-400 w-full")
	parent.Call("appendChild", input)

	input.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		td := args[0].Get("target").Get("value").String()
		todo := &Todo{
			Name:     td,
			Complete: false,
		}
		todos = append(todos, todo)
		fmt.Println(args[0].Get("target").Get("value").String())
		persistToLocalStorage()
		renderTodos()
		return nil
	}))
}

func main() {
	getFromLocalStorage()

	container.Set("classList", "flex flex-col")

	renderTodos()

	// ticker := time.NewTicker(1 * time.Second)
	// go func() {
	// 	for range ticker.C {
	// 		renderTodos()
	// 	}
	// }()

	select {}
}
