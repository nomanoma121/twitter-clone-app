import { useAuth } from "../../provider/auth";
import { useTodos } from "./internal/hook/use-todos";
import { TodoCreateForm } from "./internal/components/todo-create-form/todo-create-form";
import { TodoList } from "./internal/components/todo-list";

export const Todo = () => {
  const { user } = useAuth();
  const { todos, fetchTodos } = useTodos();

  if (!user) return null;

  return (
    <div className="Todo">
      <h1>Todo</h1>
      <p>ようこそ、{user.name}さん。今日もタスクを頑張りましょう！</p>
      <TodoCreateForm refetch={fetchTodos} />
      <TodoList todos={todos} refetch={fetchTodos} />
    </div>
  );
};
