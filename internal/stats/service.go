package stats

type StatsService interface {
    Get(userID uint) (*UserStats, error)

    IncrementPosts(userID uint)
    IncrementFollowers(userID uint)
    DecrementFollowers(userID uint)

    IncrementDogs(userID uint)
    DecrementDogs(userID uint)
}
