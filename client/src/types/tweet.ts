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
    liked_by_user: boolean
    created_at: string
}
