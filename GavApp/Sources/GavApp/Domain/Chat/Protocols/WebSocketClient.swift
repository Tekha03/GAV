import Foundation

public protocol WebSocketClient: AnyObject {
    func connect(to chatID: UUID) async throws
    func disconnect()

    func sendTyping(chatID: UUID) async throws

    func onMessageReceived(
        handler: @escaping (Message) -> Void
    )
}

extension WebSocketClient {
    func sendTyping(chatID: UUID) async throws {

    }

    func onMessageReceived(handler: @escaping (Message) -> Void) {

    }
}