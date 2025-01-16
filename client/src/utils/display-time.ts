export const displayTime = (created_at: string) => {
  const now = new Date();
  const createdAt = new Date(created_at);
  const diff = now.getTime() - createdAt.getTime();
  const diffSec = Math.floor(diff / 1000);
  if (diffSec < 60) {
    return `${diffSec}秒前`;
  }
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) {
    return `${diffMin}分前`;
  }
  const diffHour = Math.floor(diffMin / 60);
  if (diffHour < 24) {
    return `${diffHour}時間前`;
  }
  const diffDay = Math.floor(diffHour / 24);
  if (diffDay < 7) {
    return `${diffDay}日前`;
  }
  if (now.getFullYear() !== createdAt.getFullYear()) {
    return `${createdAt.getFullYear()}年${
      createdAt.getMonth() + 1
    }月${createdAt.getDate()}日`;
  }
  return `${createdAt.getMonth() + 1}月${createdAt.getDate()}日`;
};
