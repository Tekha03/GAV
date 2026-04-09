//
//  DogRepo.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

protocol DogRepository {
    func fetchDogs(ownerId: UUID) async -> [Dog]
}
