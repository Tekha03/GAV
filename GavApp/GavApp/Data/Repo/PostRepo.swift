//
//  PostRepo.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

protocol PostRepository {
    func fetchPosts(userId: UUID) async -> [Post]
    func likePost(postId: UUID) async -> Bool
}
