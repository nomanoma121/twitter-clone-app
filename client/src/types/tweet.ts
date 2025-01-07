export type TTweet = {
    id: number
    user: {
        id: number
        name: string
        email: string
    }
    content: string
    retweet: TTweet | null
}
