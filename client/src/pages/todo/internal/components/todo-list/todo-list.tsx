import "./todo-list.css";
import { TodoListItem } from "../todo-list-item";
import { TTodo } from "../../../../../types/todo";

interface TodoListProps {
  todos: TTodo[];
  refetch: () => void;
}

export const TodoList = ({ todos, refetch }: TodoListProps) => {
  return (
    <div className="Todo__list">
      {todos.map((todo) => (
        <TodoListItem key={todo.id} todo={todo} refetch={refetch} />
      ))}
    </div>
  );
};
