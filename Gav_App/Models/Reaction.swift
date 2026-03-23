//
//  Reaction.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct Reaction: Identifiable {
    let id: UUID
    let messageID: UUID
    let userID: UUID
    let emoji: String
}