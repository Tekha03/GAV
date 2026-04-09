//
//  DogModel.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

struct Dog: Identifiable {
    let id: UUID
    let ownerId: UUID

    let name: String
    let breed: String
    let photoUrl: String

    let status: String
    let age: String
    let gender: String
}
