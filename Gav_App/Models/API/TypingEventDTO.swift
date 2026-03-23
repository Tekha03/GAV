//
//  TypingEventDTO.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct TypingEventDTO: Codable {
    let eventType: TypingEventType
    let chatID: UUID
    let userID: UUID
    let timestamp: Date
}

enum TypingEventType: String, Codable {
    case typingStart = "typing_start"
    case typingStop  = "typing_stop"
}