import Domain
import SharedModels

public struct ProfileStatsMapper {
    public static func from(model: ProfileStatsModel) -> Domain.ProfileStats {
        return Domain.ProfileStats(
            userId: model.userId,
            postCount: model.postCount,
            followersCount: model.followersCount,
            followingsCount: model.followingsCount
        )
    }

    public static func to(model: Domain.ProfileStats) -> ProfileStatsModel {
        return ProfileStatsModel(
            userId: model.userId,
            postCount: model.postCount,
            followersCount: model.followersCount,
            followingsCount: model.followingsCount
        )
    }
}