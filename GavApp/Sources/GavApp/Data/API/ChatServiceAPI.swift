import Foundation

@available(macOS 12.0, *)
final class ChatServiceAPI: ChatUseCase, @unchecked Sendable {
    private let base: BaseAPI
    private let currentUserIdProvider: @Sendable () -> UUID?

    init(
        baseURL: URL,
        session: URLSession = .shared,
        authManager: AuthManager,
        currentUserIdProvider: @escaping @Sendable () -> UUID?
    ) {
        self.base = BaseAPI(baseURL: baseURL, session: session, authManager: authManager)
        self.currentUserIdProvider = currentUserIdProvider
    }

    func createPrivateChat(user1: UUID, user2: UUID) async throws -> Chat {
        let request = CreatePrivateChatRequest(userId1: user1, userId2: user2)
        let data = try await base.request(
            "/api/v1/chats/private",
            method: "POST",
            body: try JSONEncoder().encode(request),
            requiresAuth: false
        )
        return try chatDecoder.decode(ChatEnvelope.self, from: data).chat.domain
    }

    func createGroupChat(title: String, creator: UUID, members: [UUID]) async throws -> Chat {
        let request = CreateGroupChatRequest(title: title, creatorId: creator, memberIds: members)
        let data = try await base.request(
            "/api/v1/chats/group",
            method: "POST",
            body: try JSONEncoder().encode(request),
            requiresAuth: false
        )
        return try chatDecoder.decode(ChatEnvelope.self, from: data).chat.domain
    }

    func getChatByID(id: UUID) async throws -> Chat {
        let data = try await base.request("/api/v1/chats/\(id.uuidString)", requiresAuth: false)
        return try chatDecoder.decode(ChatEnvelope.self, from: data).chat.domain
    }

    func getUserChats(userID: UUID) async throws -> [Chat] {
        let data = try await base.request("/api/v1/users/\(userID.uuidString)/chats", requiresAuth: false)
        return try chatDecoder.decode(ChatListEnvelope.self, from: data).chats.map(\.domain)
    }

    func updateChatTitle(chatID: UUID, title: String) async throws {
        throw ChatAPIError.notImplemented
    }

    func updateChatPhoto(chatID: UUID, photoUrl: String) async throws {
        throw ChatAPIError.notImplemented
    }

    func leaveChat(userID: UUID, chatID: UUID) async throws {
        throw ChatAPIError.notImplemented
    }

    func deleteChat(chatID: UUID) async throws {
        throw ChatAPIError.notImplemented
    }

    func getChatMembers(chatID: UUID) async throws -> [ChatMember] {
        let data = try await base.request("/api/v1/chats/\(chatID.uuidString)/members", requiresAuth: false)
        return try chatDecoder.decode(ChatMemberListEnvelope.self, from: data).members.map(\.domain)
    }

    func addMember(userID: UUID, chatID: UUID) async throws {
        throw ChatAPIError.notImplemented
    }

    func removeMember(userID: UUID, chatID: UUID) async throws {
        throw ChatAPIError.notImplemented
    }

    func getMessages(chatID: UUID, limit: Int, before: UUID?) async throws -> [Message] {
        var path = "/api/v1/chats/\(chatID.uuidString)/messages?limit=\(limit)"
        if let before {
            path += "&before=\(before.uuidString)"
        }
        let data = try await base.request(path, requiresAuth: false)
        return try chatDecoder.decode(MessageListEnvelope.self, from: data).messages.map(\.domain)
    }

    func sendMessage(
        chatID: UUID,
        text: String?,
        attachments: [AttachmentInput]?,
        replyToId: UUID?
    ) async throws -> Message {
        guard let currentUserId = currentUserIdProvider() else {
            throw ChatAPIError.missingCurrentUser
        }

        let request = SendMessageRequest(
            senderId: currentUserId,
            text: text,
            replyToId: replyToId,
            attachments: attachments?.map(AttachmentRequest.init(input:)) ?? []
        )
        let data = try await base.request(
            "/api/v1/chats/\(chatID.uuidString)/messages",
            method: "POST",
            body: try JSONEncoder().encode(request),
            requiresAuth: false
        )
        return try chatDecoder.decode(MessageEnvelope.self, from: data).message.domain
    }

    func markAsRead(chatID: UUID, userID: UUID) async throws {
        let data = try JSONEncoder().encode(["user_id": userID.uuidString])
        _ = try await base.request(
            "/api/v1/chats/\(chatID.uuidString)/read",
            method: "POST",
            body: data,
            requiresAuth: false
        )
    }

    func sendTyping(chatID: UUID) async throws {}

    private var chatDecoder: JSONDecoder {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let raw = try container.decode(String.self)
            if let date = ISO8601DateFormatter.gavDate(from: raw) {
                return date
            }
            throw DecodingError.dataCorruptedError(
                in: container,
                debugDescription: "Invalid ISO8601 date: \(raw)"
            )
        }
        return decoder
    }
}

private enum ChatAPIError: LocalizedError {
    case missingCurrentUser
    case notImplemented

    var errorDescription: String? {
        switch self {
        case .missingCurrentUser:
            return "Current user is missing"
        case .notImplemented:
            return "This chat operation is not implemented yet"
        }
    }
}

private struct ChatEnvelope: Decodable {
    let chat: ChatDTO
}

private struct ChatListEnvelope: Decodable {
    let chats: [ChatDTO]
}

private struct MessageEnvelope: Decodable {
    let message: MessageDTO
}

private struct MessageListEnvelope: Decodable {
    let messages: [MessageDTO]
}

private struct ChatMemberListEnvelope: Decodable {
    let members: [ChatMemberDTO]
}

private struct ChatDTO: Decodable {
    let id: UUID
    let isGroup: Bool
    let title: String?
    let photoUrl: String?
    let createdAt: Date

    var domain: Chat {
        Chat(
            id: id,
            isGroup: isGroup,
            title: title?.isEmpty == false ? title! : "Чат",
            photoUrl: photoUrl ?? "",
            createdAt: createdAt
        )
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case isGroup = "is_group"
        case title
        case photoUrl = "photo_url"
        case createdAt = "created_at"
    }
}

private struct MessageDTO: Decodable {
    let id: UUID
    let chatId: UUID
    let senderId: UUID
    let text: String?
    let replyToId: UUID?
    let createdAt: Date
    let editedAt: Date?
    let attachments: [AttachmentDTO]?

    var domain: Message {
        Message(
            id: id,
            chatId: chatId,
            senderId: senderId,
            text: text,
            replyToId: replyToId,
            createdAt: createdAt,
            editedAt: editedAt,
            attachments: attachments?.map(\.domain) ?? [],
            reactions: []
        )
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case chatId = "chat_id"
        case senderId = "sender_id"
        case text
        case replyToId = "reply_to_id"
        case createdAt = "created_at"
        case editedAt = "edited_at"
        case attachments
    }
}

private struct AttachmentDTO: Decodable {
    let id: UUID?
    let messageID: UUID?
    let url: String
    let type: String
    let fileName: String
    let fileSize: Int64

    var domain: Attachment {
        Attachment(
            id: id ?? UUID(),
            messageID: messageID ?? UUID(),
            url: url,
            type: AttachmentType(apiValue: type),
            fileName: fileName,
            fileSize: fileSize
        )
    }

    private enum CodingKeys: String, CodingKey {
        case id
        case messageID = "message_id"
        case url
        case type
        case fileName = "file_name"
        case fileSize = "file_size"
    }
}

private struct ChatMemberDTO: Decodable {
    let chatId: UUID
    let userId: UUID
    let role: MemberRole
    let muted: Bool
    let lastReadMessageId: UUID?

    var domain: ChatMember {
        ChatMember(
            userId: userId,
            chatId: chatId,
            role: role,
            muted: muted,
            lastReadMessageId: lastReadMessageId
        )
    }

    private enum CodingKeys: String, CodingKey {
        case chatId = "chat_id"
        case userId = "user_id"
        case role
        case muted
        case lastReadMessageId = "last_read_message_id"
    }
}

private struct CreatePrivateChatRequest: Encodable {
    let userId1: UUID
    let userId2: UUID

    private enum CodingKeys: String, CodingKey {
        case userId1 = "user_id_1"
        case userId2 = "user_id_2"
    }
}

private struct CreateGroupChatRequest: Encodable {
    let title: String
    let creatorId: UUID
    let memberIds: [UUID]

    private enum CodingKeys: String, CodingKey {
        case title
        case creatorId = "creator_id"
        case memberIds = "member_ids"
    }
}

private struct SendMessageRequest: Encodable {
    let senderId: UUID
    let text: String?
    let replyToId: UUID?
    let attachments: [AttachmentRequest]

    private enum CodingKeys: String, CodingKey {
        case senderId = "sender_id"
        case text
        case replyToId = "reply_to_id"
        case attachments
    }
}

private struct AttachmentRequest: Encodable {
    let url: String
    let type: String
    let fileName: String
    let fileSize: Int64

    init(input: AttachmentInput) {
        url = input.url
        type = input.type.apiValue
        fileName = input.fileName
        fileSize = input.fileSize
    }

    private enum CodingKeys: String, CodingKey {
        case url
        case type
        case fileName = "file_name"
        case fileSize = "file_size"
    }
}

private extension AttachmentType {
    init(apiValue: String) {
        switch apiValue {
        case "image":
            self = .image
        case "video":
            self = .video
        case "voice", "audio":
            self = .audio
        default:
            self = .document
        }
    }

    var apiValue: String {
        switch self {
        case .image:
            return "image"
        case .video:
            return "video"
        case .audio:
            return "voice"
        case .document:
            return "file"
        }
    }
}

private extension ISO8601DateFormatter {
    static func gavDate(from value: String) -> Date? {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: value) {
            return date
        }

        formatter.formatOptions = [.withInternetDateTime]
        return formatter.date(from: value)
    }
}
