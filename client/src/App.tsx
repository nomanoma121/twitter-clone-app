import { Route, Routes, useLocation } from "react-router";
import { lazy } from "react";
import { AuthGuard } from "./components/auth-guard";
import { TweetDetail } from "./pages/tweet-detail";
import { Modal } from "./components/modal";
import { TopLayout } from "./layouts/top-layouts";
import { DefaultLayout } from "./layouts/default-layouts";

const Top = lazy(() => import("./pages/top"));
const Login = lazy(() => import("./pages/login"));
const SignUp = lazy(() => import("./pages/signup"));
const Home = lazy(() => import("./pages/home"));
const Tweet = lazy(() => import("./pages/tweet"));
const User = lazy(() => import("./pages/user"));
const Follow = lazy(() => import("./pages/follow"));

function App() {
  const location = useLocation();

  const state = location.state as { background: Location };
  const background = state?.background;
  return (
    <>
      <Routes location={background || location}>
        {/* 未ログイン時のルート */}
        <Route element={<TopLayout />}>
          <Route path="/" element={<Top />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
        </Route>

        {/* ログイン後のルート */}
        <Route element={<DefaultLayout />}>
          <Route
            path="/home"
            element={
              <AuthGuard>
                <Home />
              </AuthGuard>
            }
          />
          <Route path="/:displayID/status/:tweetID" element={<TweetDetail />} />
          <Route path="/:displayID" element={<User />} />
          <Route path="/:displayID/following" element={<Follow />} />
          <Route path="/:displayID/followers" element={<Follow />} />
        </Route>
      </Routes>
      {/* モーダル用のルート */}
      {background && location.pathname === "/compose/tweet" && (
        <Modal>
          <Tweet />
        </Modal>
      )}
    </>
  );
}

export default App;
