//
//  PinnedMessage.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct PinnedMessage: Identifiable {
    let chatID: UUID
    let messageID: UUID
    let pinnedAt: Date
}