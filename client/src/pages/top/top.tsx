import { Link } from "react-router";

export const Top = () => {
  return (
    <div>
      <h1>Welcome to You Do!</h1>
      <p>
        このアプリは、タスク管理アプリです。まずは
        <Link to="/login">ログイン</Link>
        をして、あなたのタスクを管理しましょう。
      </p>
    </div>
  );
};
