import { FaTwitter } from "react-icons/fa";
import { GoHomeFill } from "react-icons/go";
import { IoSearchOutline } from "react-icons/io5";
import { PiBell } from "react-icons/pi";
import { LuMail } from "react-icons/lu";
import { GiFeather } from "react-icons/gi";
import { RiLogoutBoxLine } from "react-icons/ri";
import { useNavigate } from "react-router";
import { useLocation } from "react-router";
import { useAuth } from "../../provider/auth";
import { UserIcon } from "../user-icon/user-icon";
import "./sidebar.css";

export const Sidebar = () => {
  const { user } = useAuth();
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <div className="Sidebar">
      <div
        className="Sidebar__icon__container"
        onClick={() => navigate("/home")}
      >
        <FaTwitter className="Sidebar__icon color__blue" />
      </div>
      <div
        className="Sidebar__icon__container"
        onClick={() => navigate("/home")}
      >
        <GoHomeFill className="Sidebar__icon color__blue" />
      </div>
      <div className="Sidebar__icon__container">
        <IoSearchOutline className="Sidebar__icon" />
      </div>
      <div className="Sidebar__icon__container">
        <PiBell className="Sidebar__icon" />
      </div>
      <div className="Sidebar__icon__container">
        <LuMail className="Sidebar__icon" />
      </div>
      <div className="Sidebar__icon__container">
        <RiLogoutBoxLine  className="Sidebar__icon" onClick={() => auth.logout()} />
      </div>
      <div
        className="Tweet__icon__container"
        onClick={() =>
          navigate("/compose/tweet", {
            state: {
              background: location,
              tweet: { type: "tweet" },
            },
          })
        }
      >
        <GiFeather className="Tweet__icon" />
      </div>
      {user && (
        <div
          className="Sidebar__user__container"
          onClick={() => navigate(`/${user.display_id}`)}
        >
          <UserIcon user={user} size={40} />
        </div>
      )}
    </div>
  );
};
