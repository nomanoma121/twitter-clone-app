import { Route, Routes, useLocation } from "react-router";
import { Header } from "./components/header";
import { lazy } from "react";
import { Container } from "./components/container/container";
import { AuthGuard } from "./components/auth-guard";
import { TweetDetail } from "./pages/tweet-detail";
import { Modal } from "./components/modal";

const Top = lazy(() => import("./pages/top"));
const Todo = lazy(() => import("./pages/todo"));
const Login = lazy(() => import("./pages/login"));
const SignUp = lazy(() => import("./pages/signup"));
const Home = lazy(() => import("./pages/home"));
const Tweet = lazy(() => import("./pages/tweet"));
const User = lazy(() => import("./pages/user"));

function App() {
  const location = useLocation();

  const state = location.state as { background: Location };
  const background = state?.background;
  return (
    <>
      <Header />
      <Container>
        <Routes location={background || location}>
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
          <Route path="/:displayId" element={<User />} />
        </Routes>
        
        {/* モーダル用のルート */}
        {background && location.pathname === "/compose/tweet" && (
          <Modal>
            <Tweet />
          </Modal>
        )}

      </Container>
    </>
  );
}

export default App;
