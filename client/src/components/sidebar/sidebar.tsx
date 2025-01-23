import "./sidebar.css";
import { FaRegUser, FaTwitter } from "react-icons/fa";
import { GoHomeFill } from "react-icons/go";
import { IoSearchOutline } from "react-icons/io5";
import { PiBell } from "react-icons/pi";
import { LuMail } from "react-icons/lu"; 
import { CiCircleMore } from "react-icons/ci";

export const Sidebar = () => {
  return (
    <div className="Sidebar">
      <div className="Sidebar__icon">
        <FaTwitter />
      </div>
      <div className="Sidebar__icon">
      <GoHomeFill />
      </div>
      <div className="Sidebar__icon">
        <IoSearchOutline />
      </div>
      <div className="Sidebar__icon">
        <PiBell />
      </div>
      <div className="Sidebar__icon">
        <LuMail />
      </div>
      <div className="Sidebar__icon">
        <FaRegUser />
      </div>
      <div className="Sidebar__icon">
        <CiCircleMore />
      </div>
    </div>
  )
}
