import Foundation

public protocol ReactionRepository {
    func addReaction(
        messageID: UUID,
        userID: UUID,
        emoji: String
    ) async throws

    func removeReaction(
        messageID: UUID,
        userID: UUID
    ) async throws

    func getReactions(messageID: UUID) async throws -> [Reaction]
}