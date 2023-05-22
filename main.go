package main

import (
	"syscall/js"
	"time"
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
			return nil
		}))

		delButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			todos = append(todos[:index], todos[index+1:]...)
			renderTodos()
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

func main() {
	todos = append(todos, &Todo{Name: "Take out the trash"})
	todos = append(todos, &Todo{Name: "Wash dishes"})
	todos = append(todos, &Todo{Name: "Do laundry"})
	todos = append(todos, &Todo{Name: "Do homework"})

	container.Set("classList", "flex flex-col")

	renderTodos()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			renderTodos()
		}
	}()

	select {}
}
