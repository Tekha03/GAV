import Foundation

public struct FollowMapper {
    public init() {}

    public func from(model: FollowModel) -> Follow {
        Follow(
            followerId: model.followerId,
            followingId: model.followingId
        )
    }

    public func to(model: Follow) -> FollowModel {
        FollowModel(
            followerId: model.followerId,
            followingId: model.followingId
        )
    }
}