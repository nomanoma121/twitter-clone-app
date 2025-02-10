import { useNavigate } from "react-router";
import { FaTwitter } from "react-icons/fa";
import { RxCross2 } from "react-icons/rx";
import "./auth-container.css";

export const AuthContainer = ({ children }: React.PropsWithChildren<{}>) => {
  const navigate = useNavigate();
  return (
    <div className="AuthContainer">
      <div className="AuthContainer__content">
        <div className="AuthContainer__header">
          <div className="AuthContainer__header__side">
            <div className="AuthContainer__header__cross">
              <RxCross2
                style={{ scale: "1.5" }}
                onClick={() => navigate("/")}
              />
            </div>
          </div>
          <div className="AuthContainer__header__icon">
            <FaTwitter style={{ scale: "2" }} />
          </div>
          <div className="AuthContainer__header__side"></div>
        </div>
        {children}
      </div>
    </div>
  );
};
