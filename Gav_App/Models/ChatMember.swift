//
//  ChatMember.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct ChatMember: Identifiable {
    let userID: UUID
    let chatID: UUID
    let role: String
    let joinedAt: Date
    let leftAt: Date?
    let muted: Bool
    let lastReadMessageID: UUID
}