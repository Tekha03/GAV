package stats

type StatsRepository interface {
    GetByUserID(userID uint) (*UserStats, error)
    Update(stats *UserStats) error
}