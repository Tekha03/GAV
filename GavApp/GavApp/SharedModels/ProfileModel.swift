//
//  ProfileModel.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

struct UserProfile {
    let userId: Foundation.UUID

    let name: String
    let surname: String
    let username: String
    let bio: String
    let profilePhotoUrl: String
    var followersCount: Int
    var followingCount: Int
    var isFollowed: Bool
}

