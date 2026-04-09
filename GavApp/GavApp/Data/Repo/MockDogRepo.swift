//
//  MockDogRepo.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import Foundation

final class MockDogRepository: DogRepository {

    func fetchDogs(ownerId: UUID) async -> [Dog] {

        try? await Task.sleep(nanoseconds: 300_000_000)

        return [
            Dog(
                id: UUID(),
                ownerId: ownerId,
                name: "Джесси",
                breed: "Йоркширский терьер",
                photoUrl: "jessy",
                status: "cautious",
                age: "12",
                gender: "female"
            ),
            Dog(
                id: UUID(),
                ownerId: ownerId,
                name: "Луна",
                breed: "Аусси",
                photoUrl: "luna",
                status: "friendly",
                age: "2",
                gender: "female"
            ),
            Dog(id: UUID(),
                ownerId: ownerId,
                name: "Умка",
                breed: "Бернский зенненхунд",
                photoUrl: "umka",
                status: "friendly",
                age: "5",
                gender: "female")
        ]
    }
}
