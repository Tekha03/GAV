import Foundation

public struct UserStatsMapper {
    public init() {}

    public func from(model: UserStatsModel) -> UserStats {
        UserStats(
            id: model.id,
            userId: model.userId,
            postCount: model.postCount,
            followersCount: model.followersCount,
            followingsCount: model.followingsCount,
            dogsCount: model.dogsCount,
            createdAt: model.createdAt,
            updatedAt: model.updatedAt
        )
    }
}