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
      <FaTwitter />
      <GoHomeFill />
      <IoSearchOutline />
      <PiBell />
      <LuMail />
      <FaRegUser />
      <CiCircleMore />
    </div>
  )
}
