//
//  PostModel.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

struct Post: Identifiable {
    let id: UUID
    let userId: UUID
    let date: Date
    let content: String
    let imageUrl: String?
    var likesCount: Int
    var commentsCount: Int
    var isLiked: Bool
}
