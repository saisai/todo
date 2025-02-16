"use client"
import { useEffect, useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import { z } from "zod";
import { fetchTodos, addTodo, updateTodo, deleteTodo } from "@/store/todoSlice"

const todoSchema = z.object({
  title: z.string().min(1, "Todo cannot be empty").max(100, "Too long!"),
});

export default function TodoList() {
    const dispatch = useDispatch()
    const todos = useSelector((state) => state.todos.todos) || [];
    const [title, setTitle] = useState("")
    const [editId, setEditId] = useState(null);
    const [editTitle, setEditTitle] = useState("");
    const [error, setError] = useState("");

    useEffect(() => {
        dispatch(fetchTodos());
    }, [dispatch])

    const handleSubmit = (e) => {
      e.preventDefault(); // Prevent page refresh
  
      // Validate using Zod only on button click
      const validation = todoSchema.safeParse({ title });
  
      if (!validation.success) {
        setError(validation.error.errors[0].message);
        return;
      }
  
      // Dispatch action and reset state
      dispatch(addTodo(title));
      setTitle("");
      setError(""); // Clear error message after successful submission
    };

    return (
        <div className="max-w-md mx-auto w-full mt-10 p-5 bg-white ">
          <h1 className="text-xl font-bold mb-4">Todo List</h1>
          <form  className="flex flex-col gap-2 p-4">
          <div className="flex gap-2">
            <input 
              type="text"
              
              value={title}
             placeholder="Add a new todo..."
             onChange={(e) => setTitle(e.target.value)} 
              className="border border-gray-300 h-8 p-2 rounded w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
             />
            <button  
              className="bg-blue-500 h-8 text-white px-2 py-1 rounded focus:outline-none focus:ring-2 focus:ring-blue-500" 
              onClick={handleSubmit}>
                Add Todo</button>
          </div>
          {error && <p className="text-red-500 text-sm">{error}</p>}
          </form>
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

