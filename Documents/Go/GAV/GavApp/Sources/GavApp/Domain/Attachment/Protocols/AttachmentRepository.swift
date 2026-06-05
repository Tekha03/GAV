import Foundation

public protocol AttachmentRepository {
    func addAttachment(
        messageID: UUID,
        input: AttachmentInput
    ) async throws -> Attachment
}
//    func getAttachments(messageID: UUID) async throws -> [Attachment]