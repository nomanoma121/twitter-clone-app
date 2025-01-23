import "./sidebar.css";
import { FaRegUser, FaTwitter } from "react-icons/fa";
import { GoHomeFill } from "react-icons/go";
import { IoSearchOutline } from "react-icons/io5";
import { PiBell } from "react-icons/pi";
import { LuMail } from "react-icons/lu";
import { CiCircleMore } from "react-icons/ci";
import { GiFeather } from "react-icons/gi";
export const Sidebar = () => {
  return (
    <div className="Sidebar">
      <div className="Sidebar__icon__container">
        <FaTwitter className="Sidebar__icon" />
      </div>
      <div className="Sidebar__icon__container">
        <GoHomeFill className="Sidebar__icon" />
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
        <FaRegUser className="Sidebar__icon" />
      </div>
      <div className="Sidebar__icon__container">
        <CiCircleMore className="Sidebar__icon" />
      </div>
      <div className="Tweet__icon__container">
        <GiFeather className="Tweet__icon" />
      </div>
    </div>
  );
};
