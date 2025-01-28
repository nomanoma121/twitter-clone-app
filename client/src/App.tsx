import { Route, Routes, useLocation } from "react-router";
import { Sidebar } from "./components/sidebar";
import { lazy } from "react";
import { Container } from "./components/container/container";
import { AuthGuard } from "./components/auth-guard";
import { TweetDetail } from "./pages/tweet-detail";
import { Modal } from "./components/modal";

const Top = lazy(() => import("./pages/top"));
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
      <Container>
        <Sidebar />
        <div style={{ maxWidth: "632px", width: "100%", borderLeft: "1px solid #e1e8ed", borderRight: "1px solid #e1e8ed" }}>
          <Routes location={background || location}>
            <Route path="/" element={<Top />} />
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<SignUp />} />
            <Route
              path="/home"
              element={
                <AuthGuard>
                  <Home />
                </AuthGuard>
              }
            />
            <Route
              path="/:displayID/status/:tweetID"
              element={<TweetDetail />}
            />
            <Route path="/:displayID" element={<User />} />
          </Routes>

          {/* モーダル用のルート */}
          {background && location.pathname === "/compose/tweet" && (
            <Modal>
              <Tweet />
            </Modal>
          )}
        </div>
      </Container>
    </>
  );
}

export default App;
