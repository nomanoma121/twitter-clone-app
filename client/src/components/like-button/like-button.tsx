import { TTweet } from '../../types';
import { BsHeart } from 'react-icons/bs';
import './like-button.css';

export const LikeButton = ({ tweet }: TTweet ) => {
  const handleLike = async (e) => {
    e.stopPropagation();
    console.log('いいねしました');
  }

  return (
    <div onClick={(e) => handleLike(e)} style={{ cursor: 'pointer', zIndex: 999 }} className="LikeButton">
      <BsHeart />
      <span>{tweet.interactions.like_count}</span>
    </div>
  );
}
