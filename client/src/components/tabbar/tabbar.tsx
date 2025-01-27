import "./tabbar.css";

export const Tabbar = () => {
  // タイムラインとフォロー中を切り替えるタブバー
  return (
    <div className="Tabbar">
      <div className="Tabbar__tab">タイムライン</div>
      <div className="Tabbar__tab">フォロー中</div>
    </div>
 );
}
