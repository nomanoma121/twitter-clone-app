# twitter-clone-server

このプロジェクトでは、以下のER図に基づいてデータベースを設計しています。

## ER図
![ER Diagram](./docs/twitter-clone-ER.png)

- Draw.ioファイルは[こちら](./docs/twitter-clone-ER.png)

## API設計

### タイムラインのツイートを取得
`GET api/tweets?type=timeline&cursor...=&limit=...`
#### レスポンス例
```json
{
  {
    "id": 1,
    "user": {
      "id": 2,
      "display_id": "hoge",
      "name": "hoge",
      "icon_image": "...",
    },
    "content": "これはリツイート",
    "tweet_image": "..."
    "retweeted_tweet": {
      "id": 101,
      "user": {
        "id": 5,
        "display_id": "fuga",
        "name": "fuga",
        "icon_image": "..."
      },
      "content": "これはリツイートされたツイート",
      "tweet_image": "..."
      "interactions": {
        "retweet": 3,
        "reply": 6,
        "like": 7,
      }
    },
    "interactions": {
      "retweet": 5,
      "reply": 4,
      "like": 5,
    }
    "type": "retweet"
  },
  {
    ...
  }
}
```

### フォローしている人だけのTweetを取得
`GET api/tweets?type=follows&cursor=...&limit=...`

レスポンスはタイムラインのツイート取得と同じ

### あるユーザーのツイートを取得
`GET api/users/:id/tweets?cursor=...&limit=...`

これも上と同じ

### あるツイートの詳細を取得
`GET api/tweets/:id`


### ツイートする
`POST api/tweets`

リクエスト例
```json
{
  "content": "ツイート",
  "tweet_image": "...",
}
```

### リツイートする
`POST api/tweets/:id/retweet`

リクエスト例
```json
{
  "content": "リツイート",
  "tweet_image": "...",
  "retweet_id": 3,
}
```

### リプライする
`POST api/tweets/:id/reply`

```json
{
  "content": "リプライ",
  "tweet_image": "...",
  "reply_id": 4,
}
```
