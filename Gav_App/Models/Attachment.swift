//
//  Attachment.swift
//  
//
//  Created by Виктория Кашуркина on 23.03.2026.
//


struct Attachment: Identifiable {
    let id: UUID
    let messageID: UUID
    let url: String
    let type: String
    let fileName: String
    let fileSize: Int64
}