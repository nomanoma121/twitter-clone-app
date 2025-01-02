import { Link } from "react-router";
import "./header.css";
import { useAuth } from "../../provider/auth";
import { Button } from "../button";

export const Header = () => {
  const auth = useAuth();

  return (
    <header className="Header">
      <div className="Header__inner">
        <Link to="/">
          <img src="/logo.svg" alt="You Do Logo" height="30" />
        </Link>

        <nav>
          <ul className="Header__nav">
            <li>
              <Link to="/" className="Header__nav-item">
                Home
              </Link>
            </li>
            {auth.user && (
              <li>
                <Link to="/todos" className="Header__nav-item">
                  Todos
                </Link>
              </li>
            )}
            {auth.user ? (
              <li>
                <Button
                  onClick={() => {
                    auth.logout();
                  }}
                  className="Header__nav-item--button-like"
                >
                  Logout
                </Button>
              </li>
            ) : (
              <li>
                <Link to="/login" className="Header__nav-item--button-like">
                  Login
                </Link>
              </li>
            )}
          </ul>
        </nav>
      </div>
    </header>
  );
};
