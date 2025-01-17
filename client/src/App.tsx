import { Route, Routes } from "react-router";
import { Header } from "./components/header";
import { lazy } from "react";
import { Container } from "./components/container/container";
import { AuthGuard } from "./components/auth-guard";
import { TweetDetail } from "./pages/tweet-detail";

const Top = lazy(() => import("./pages/top"));
const Todo = lazy(() => import("./pages/todo"));
const Login = lazy(() => import("./pages/login"));
const SignUp = lazy(() => import("./pages/signup"));
const Home = lazy(() => import("./pages/home"));
const Tweet = lazy(() => import("./pages/tweet"));
const User = lazy(() => import("./pages/user"));

function App() {
  return (
    <>
      <Header />
      <Container>
        <Routes>
          <Route path="/" element={<Top />} />
          <Route
            path="/todos"
            element={
              <AuthGuard>
                <Todo />
              </AuthGuard>
            }
          />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/home" element={<Home />} />
          <Route path="/:displayId/status/:tweetId" element={<TweetDetail />} />
          <Route path="/compose/tweet" element={<Tweet />} />
          <Route path="/:displayId" element={<User />} />
        </Routes>
      </Container>
    </>
  );
}

export default App;
