//
//  ProfileRepo.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

protocol ProfileRepository {
    func fetchProfile(userId: UUID) async -> UserProfile
    func followUser(userId: UUID) async -> Bool
}
