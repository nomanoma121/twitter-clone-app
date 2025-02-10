import "./signup.css";
import { useActionState, useCallback } from "react";
import { AuthPageContainer } from "../../components/auth-page-container";
import { serverFetch } from "../../utils/fetch";
import { Input } from "../../components/input";
import { Button } from "../../components/button/button";
import { Link, useNavigate } from "react-router";
import { useAuth } from "../../provider/auth";

type SignUpFormStateType = {
  message: string;
};

export const SignUp = () => {
  const { setUser } = useAuth();
  const navigate = useNavigate();

  const SignUpFormAction = useCallback(
    async (
      _prevState: SignUpFormStateType,
      formData: FormData
    ): Promise<SignUpFormStateType> => {
      const name = formData.get("name");
      const email = formData.get("email");
      const password = formData.get("password");
      const display_id = formData.get("display_id");

      if (
        name === null ||
        display_id === null ||
        email === null ||
        password === null
      ) {
        return { message: "メールアドレスとパスワードを入力してください" };
      }

      if (
        typeof name !== "string" ||
        typeof display_id !== "string" ||
        typeof email !== "string" ||
        typeof password !== "string"
      ) {
        return { message: "不正な入力です" };
      }

      if (password.length < 8) {
        return { message: "パスワードは8文字以上で入力してください" };
      }

      const res = await serverFetch("/auth/signup", {
        method: "POST",
        body: JSON.stringify({ name, email, password, display_id }),
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
        message: "登録に失敗しました",
      };
    },
    [setUser, navigate]
  );

  const [error, submitAction, isPending] = useActionState(SignUpFormAction, {
    message: "",
  });

  return (
    <AuthPageContainer>
      <form action={submitAction} className="SignUp">
        <h1>アカウント作成</h1>
        <table className="SignUp__FormTable">
          <tbody>
            <tr>
              <td>名前</td>
              <td>
                <Input type="text" name="name" required />
              </td>
            </tr>
            <tr>
              <td>ユーザーID</td>
              <td>
                <Input type="text" name="display_id" required />
              </td>
            </tr>
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
        <Button type="submit" disabled={isPending}>
          登録
        </Button>
        {error.message && <p className="SignUp__Error">{error.message}</p>}
        <p>
          アカウントをお持ちの方は
          <Link to="/login">こちら</Link>
        </p>
      </form>
    </AuthPageContainer>
  );
};
