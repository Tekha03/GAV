import Foundation

struct ChatMessageUI: Identifiable {
    let id: UUID
    let message: Message
    let isMine: Bool
    let isPinned: Bool
}

struct ChatMessageRowModel: Identifiable {
    let id: UUID
    let message: Message
    let isMine: Bool
    let isPinned: Bool
    let isRead: Bool
}

struct ChatDaySection: Identifiable {
    let date: Date
    let rows: [ChatMessageRowModel]

    var id: Date { date }
}
