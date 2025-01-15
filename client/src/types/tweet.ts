export type TTweet = {
    id: number
    user: {
        id: number
        name: string
        display_id: string
        icon_url: string
    }
    content: string
    retweet: {
        id: number
        user: {
            id: number
            name: string
            display_id: string
            icon_url: string
        }
        content: string
        interactions: {
            retweet_count: number
            reply_count: number
            like_count: number
        }
        created_at: string
    }
    interactions: {
        retweet_count: number
        reply_count: number
        like_count: number
    }
    created_at: string
}

// {
//   {
//     "id": 1,
//     "user": {
//       "id": 2,
//       "display_id": "hoge",
//       "name": "hoge",
//       "icon_url": "...",
//     },
//     "content": "これはリツイート",
//     "retweet": {
//       "id": 101,
//       "user": {
//         "id": 5,
//         "display_id": "fuga",
//         "name": "fuga",
//         "icon_url": "..."
//       },
//       "content": "これはリツイートされたツイート",
//       "interactions": {
//         "retweet_count": 3,
//         "reply_count": 6,
//         "like_count": 7,
//       }
//       "created_at": "2022-11-05",
//     },
//     "interactions": {
//       "retweet_count": 5,
//       "reply_count": 4,
//       "like_count": 5,
//     }
//     "created_at": "2024-12-11",
//   },
//   {
//     ...
//   }
// }
