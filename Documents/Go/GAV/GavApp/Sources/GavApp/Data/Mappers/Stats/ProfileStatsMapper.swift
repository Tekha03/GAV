public struct ProfileStatsMapper {
    public static func from(model: ProfileStatsModel) -> ProfileStats {
        return ProfileStats(
            userId: model.userId,
            postCount: model.postCount,
            followersCount: model.followersCount,
            followingsCount: model.followingsCount
        )
    }

    public static func to(model: ProfileStats) -> ProfileStatsModel {
        return ProfileStatsModel(
            userId: model.userId,
            postCount: model.postCount,
            followersCount: model.followersCount,
            followingsCount: model.followingsCount
        )
    }
}