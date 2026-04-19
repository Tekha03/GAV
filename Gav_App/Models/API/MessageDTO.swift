//
//  MessageDTO.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct MessageDTO: Codable, Identifiable {
    let id: UUID
    let chatID: UUID
    let senderID: UUID
    let text: String?
    let replyToID: UUID?
    let createdAt: Date
    let editedAt: Date?
}