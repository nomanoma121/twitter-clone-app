import { useNavigate } from "react-router";
import { FaTwitter } from "react-icons/fa";
import { IoSearchOutline } from "react-icons/io5";
import { FaRegComment } from "react-icons/fa";
import { FiUsers } from "react-icons/fi";
import { Button } from "../../components/button";
import "./top.css";

export const Top = () => {
  const navigate = useNavigate();
  return (
    <div className="Top__container">
      <div className="Top__left__container">
        <div>
          <div className="Top__left__message">
            <IoSearchOutline size={30} style={{ scale: "1.2" }} />
            <div className="Top__left__message__content">あなたの「好き」をフォローしよう。</div>
          </div>
          <div className="Top__left__message">
            <FiUsers size={30} />
            <div className="Top__left__message__content">話題のトピックを追いかけましょう。</div>
          </div>
          <div className="Top__left__message">
            <FaRegComment size={30} style={{ scale: "0.9" }} />
            <div className="Top__left__message__content">Twitterに参加しましょう。</div>
          </div>
        </div>
      </div>
      <div className="Top__right__container">
        <div className="Top__right__header__container">
          <FaTwitter size={40} color="#1da1f2" />
          <div className="Top__right__header__message">
            「いま」起きていることを見つけよう
          </div>
          <div className="Top__right__nav__message">Twitterをはじめよう</div>
          <Button
            className="Top__right__button"
            onClick={() => navigate("/signup")}
            buttonActive={true}
          >
            アカウント作成
          </Button>
          <Button
            className="Top__right__button Top__right__button__login"
            onClick={() => navigate("/login")}
          >
            ログイン
          </Button>
        </div>
      </div>
    </div>
  );
};
