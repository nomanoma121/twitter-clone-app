import { useEffect, useState } from "react";
import { serverFetch } from "../../../../utils/fetch";

type Todo = {
  id: number;
  title: string;
  completed: boolean;
};

export const useTodos = () => {
  const [todos, setTodos] = useState<Todo[]>([]);

  const fetchTodos = async () => {
    const res = await serverFetch("/api/todos");
    setTodos(await res.json());
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  return {
    todos,
    fetchTodos,
  };
};
