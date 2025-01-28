import { useState } from "react";
import "./tabbar.css";

export const Tabbar = () => {
  const [isTimeline, setIsTimeline] = useState(true);

  return (
    <div className="Tabbar">
      <div className="Tabbar__tab" onClick={() => setIsTimeline(true)}>
        <div className={`Tabbar__tab__title ${isTimeline ? "active" : "" }`}>タイムライン</div>
        <div className={`Tabbar__tab__border ${isTimeline ? "active" : "" }`} style={{ width: "34%" }} ></div>
      </div>
      <div className="Tabbar__tab" onClick={() => setIsTimeline(false)}>
        <div className={`Tabbar__tab__title ${isTimeline ? "" : "active" }`}>フォロー中</div>
        <div className={`Tabbar__tab__border ${isTimeline ? "" : "active" }`} style={{ width: "28%" }} ></div>
      </div>
    </div>
  );
};
