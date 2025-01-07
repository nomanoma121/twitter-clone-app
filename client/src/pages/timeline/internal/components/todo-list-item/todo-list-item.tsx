import { useActionState, useCallback, useState } from "react";
import { serverFetch } from "../../../../../utils/fetch";
import { TTodo } from "../../../../../types/todo";
import { Input } from "../../../../../components/input";
import { Edit } from "react-feather";

import "./todo-list-item.css";

interface TodoListItemProps {
  todo: TTodo;
  refetch: () => void;
}

type FormStateType = {
  message: string;
};

export const TodoListItem = ({ todo, refetch }: TodoListItemProps) => {
  const DoneAction = useCallback(async () => {
    const res = await serverFetch(`/api/todos/${todo.id}`, {
      method: "PUT",
      body: JSON.stringify({
        title: todo.title,
        completed: !todo.completed,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (res.ok) {
      refetch();
    }

    return;
  }, [todo, refetch]);

  const DeleteAction = useCallback(async () => {
    const res = await serverFetch(`/api/todos/${todo.id}`, {
      method: "DELETE",
    });

    if (res.ok) {
      refetch();
    }

    return;
  }, [todo, refetch]);

  const [editMode, setEditMode] = useState(false);

  const EditAction = useCallback(
    async (
      _prevState: FormStateType,
      formData: FormData
    ): Promise<FormStateType> => {
      const title = formData.get("title");

      const res = await serverFetch(`/api/todos/${todo.id}`, {
        method: "PUT",
        body: JSON.stringify({
          title,
          completed: todo.completed,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
        refetch();
        setEditMode(false);
        return { message: "" };
      }

      return { message: "Todoの更新に失敗しました。" };
    },
    [todo, refetch]
  );

  const [error, submitAction] = useActionState(EditAction, {
    message: "",
  });

  return (
    <div className="TodoList__item">
      {editMode ? (
        <form action={submitAction} className="TodoList__form">
          <Input
            type="text"
            name="title"
            defaultValue={todo.title}
            className="TodoList__input"
          />
          {error.message && (
            <span className="TodoList__error">{error.message}</span>
          )}
          <button type="submit" className="TodoList__update-button">
            更新
          </button>
        </form>
      ) : (
        <>
          <span
            className={`TodoList__title ${
              todo.completed ? "TodoList__title--completed" : ""
            }`}
          >
            {todo.title}
          </span>
          <button
            onClick={() => setEditMode(true)}
            className="TodoList__edit-button"
            aria-label="Edit"
            type="button"
          >
            <Edit />
          </button>
        </>
      )}
      <input
        type="checkbox"
        name="completed"
        className="TodoList__toggle-button"
        defaultChecked={todo.completed}
        onChange={DoneAction}
      />
      <button onClick={DeleteAction} className="TodoList__delete-button">
        削除
      </button>
    </div>
  );
};
