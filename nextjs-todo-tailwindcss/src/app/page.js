import Image from "next/image";

import TodoList from "../components/TodoList";

export default function Home() {
  return (
    <main>
      <div className="max-w-md mx-auto mt-10 p-5 bg-white shadow-lg rounded-lg">
        <h1 className="text-3xl font-bold underline yr">My TODO App</h1>
        <TodoList />
      </div>
    </main>
  );
}
