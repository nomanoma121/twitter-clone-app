import { IoIosArrowRoundBack } from "react-icons/io";
import { useNavigate } from "react-router";
import "./header.css";

type HeaderProps = {
  title: string;
  subtitle?: string;
};

export const Header = ({ title, subtitle }: HeaderProps) => {
  const navigate = useNavigate();
  return (
    <header className="Header" >
      <div className="Header__backButton" onClick={() => navigate(-1)} >
        <IoIosArrowRoundBack style={{ scale: 1.6 }} />
      </div>
      <div className="Header__content">
        <h1>{title}</h1>
        {subtitle && <h2>{subtitle}</h2>}
      </div>
    </header>
  );
};
