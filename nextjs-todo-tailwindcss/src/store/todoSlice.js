import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import axios from "axios";

const API_URL = "http://localhost:3001/todos"


export const fetchTodos = createAsyncThunk("todos/fetchTodos", async () => {
    const resposne = await axios.get(API_URL)
    return resposne.data
})

export const addTodo = createAsyncThunk("todos/addTodo", async( title) => {
    const response = await axios.post(API_URL, { title, completed: false})
    return response.data 
})

export const updateTodo = createAsyncThunk("todos/updateTodo", async (todo) => {
    const response = await axios.put(`${API_URL}/${todo.id}`, todo);
    return response.data;
})

export const deleteTodo = createAsyncThunk("todos/deleteTodo", async (id) => {
    console.log("deleteing id", id)
    await axios.delete(`${API_URL}/${id}`);
    return id;
  });

// Redux slice

const todoSlice = createSlice({
    name: "todos",
    initialState: { todos: [], status: "idle"},
    reducers: {},
    extraReducers: (builder) => {
        builder 
        .addCase(fetchTodos.fulfilled, (state, action) => {
            console.log("fetchTodos state", state)
            console.log("fetchTodos state todos", state.todos)
            console.log("fetchTodos testing", action.payload)
            // state.todos.push(action.payload);
            state.todos = action.payload || [];
        })
        .addCase(addTodo.fulfilled, (state, action) => {
            console.log("add state", state)
            console.log("add state todos", state.todos)
            console.log("add testing", action.payload)
            state.todos.push(action.payload);
        })
        .addCase(updateTodo.fulfilled, (state, action) => {
            const index = state.todos.findIndex((t) => t.id === action.payload.id)
            if (index !== -1) state.todos[index] = action.payload
        })
        .addCase(deleteTodo.fulfilled, (state, action) => {
            state.todos = state.todos.filter((t) => t.id !== action.payload);
        })
    }
})

export default todoSlice.reducer;