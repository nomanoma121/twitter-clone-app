import { useState } from "react";
import "./tabbar.css";

type TabbarProps = {
  titles: {
    first: string;
    second: string;
  };
  switchTab: () => void;
};

export const Tabbar = ({ titles, switchTab }: TabbarProps) => {
  const [isTimeline, setIsTimeline] = useState(true);
  const handleClick = (order) => {
    if (order === "first") {
      setIsTimeline(true);
    } else {
      setIsTimeline(false);
    }
  switchTab();
  };

  return (
    <div className="Tabbar">
      <div className="Tabbar__tab" onClick={() => handleClick("first")}>
        <div className={`Tabbar__tab__title ${isTimeline ? "active" : "" }`}>{titles.first}</div>
        <div className={`Tabbar__tab__border ${isTimeline ? "active" : "" }`} style={{ width: "34%" }} ></div>
      </div>
      <div className="Tabbar__tab" onClick={() => handleClick("second")}>
        <div className={`Tabbar__tab__title ${isTimeline ? "" : "active" }`}>{titles.second}</div>
        <div className={`Tabbar__tab__border ${isTimeline ? "" : "active" }`} style={{ width: "28%" }} ></div>
      </div>
    </div>
  );
};
