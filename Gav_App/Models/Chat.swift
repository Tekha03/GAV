//
//  Chat.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct Chat: Identifiable {
    let id: UUID
    let isGroup: Bool
    let title: String
    let photoURL: String?
    let createdAt: Date
}