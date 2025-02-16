"use client"
import { useEffect, useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import { fetchTodos, addTodo, updateTodo, deleteTodo } from "@/store/todoSlice"

export default function TodoList() {
    const dispatch = useDispatch()
    const todos = useSelector((state) => state.todos.todos)
    const [title, setTitle] = useState("")
    const [editId, setEditId] = useState(null);
    const [editTitle, setEditTitle] = useState("");

    useEffect(() => {
        dispatch(fetchTodos());
    }, [dispatch])

    return (
        <div className="max-w-md mx-auto w-full mt-10 p-5 bg-white ">
          <h1 className="text-xl font-bold mb-4">Todo List</h1>
          <div className="flex gap-2">
            <input className="border w-full rounded" value={title} onChange={(e) => setTitle(e.target.value)} />
            <button  className="bg-blue-500 h-8 text-white px-2 py-1 rounded" onClick={() => dispatch(addTodo(title))}>Add Todo</button>
          </div>
          <ul>
            {todos && todos.map((todo) => (
              <li key={todo.id}>
              {editId === todo.id ? (
                <>
                <div className="flex gap-2 pt-5">
                  <input className="border px-1 py-1 w-full rounded"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                  />
                  <button
                    className="bg-blue-500 h-8 text-white px-2 py-1 rounded"
                    onClick={() => {
                      dispatch(updateTodo({ id: todo.id, title: editTitle, completed: todo.completed }));
                      setEditId(null);
                    }}
                  >
                    Save
                  </button>
                  </div>
                </>
              ) : (
                <>
                  <div className="flex gap-2 pt-5">
                  <input
                    type="checkbox"
                    checked={todo.completed}
                    onChange={() => dispatch(updateTodo({ ...todo, completed: !todo.completed }))}
                  />
                  {todo.title}
                  <button className="bg-blue-500 h-8 text-white px-2 py-1 rounded" onClick={() => { setEditId(todo.id); setEditTitle(todo.title); }}>Edit</button>
                  <button className="bg-blue-500 h-8 text-white px-2 py-1 rounded" onClick={() => dispatch(deleteTodo(todo.id))}>Delete</button>
                  </div>
                </>
              )}
            </li>
            ))}
          </ul>
        </div>
      );

}

