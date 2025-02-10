import { useActionState, useCallback } from "react";
import { serverFetch } from "../../utils/fetch";
import { Input } from "../../components/input";
import { Button } from "../../components/button/button";
import { Link, useNavigate } from "react-router";
import { useAuth } from "../../provider/auth";
import { AuthContainer } from "../../components/auth-container";
import "./login.css";

type LoginFormStateType = {
  message: string;
};

export const Login = () => {
  const { setUser } = useAuth();
  const navigate = useNavigate();

  const LoginFormAction = useCallback(
    async (
      _prevState: LoginFormStateType,
      formData: FormData
    ): Promise<LoginFormStateType> => {
      const email = formData.get("email");
      const password = formData.get("password");

      if (email === null || password === null) {
        return { message: "メールアドレスとパスワードを入力してください" };
      }

      if (typeof email !== "string" || typeof password !== "string") {
        return { message: "不正な入力です" };
      }

      if (password.length < 8) {
        return { message: "パスワードは8文字以上で入力してください" };
      }

      const res = await serverFetch("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
        const data = await res.json();
        localStorage.setItem("token", data.token);
        setUser(data.user);
        navigate("/home");
        return { message: "" };
      }

      return {
        message: "ログインに失敗しました",
      };
    },
    [setUser, navigate]
  );

  const [error, submitAction, isPending] = useActionState(LoginFormAction, {
    message: "",
  });

  return (
    <AuthContainer>
      <form className="Login" onSubmit={submitAction}>
        <h1>Twitterにログイン</h1>
        <table className="Login__FormTable">
          <tbody>
            <tr>
              <td>メールアドレス</td>
              <td>
                <Input type="email" name="email" required />
              </td>
            </tr>
            <tr>
              <td>パスワード</td>
              <td>
                <Input type="password" name="password" required />
              </td>
            </tr>
          </tbody>
        </table>
        <Button type="submit" disabled={isPending} buttonActive={true}>
          ログイン
        </Button>
        {error.message && <p className="Login__Error">{error.message}</p>}
        <p>
          アカウントをお持ちでない場合は
          <Link to="/signup">登録</Link>
        </p>
      </form>
    </AuthContainer>
  );
};
