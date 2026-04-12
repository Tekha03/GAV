//
//  DogService.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//


import Foundation

protocol DogService {
    func createDog(ownerId: UUID, input: CreateDogInput) async throws -> Dog
    func updateDog(ownerId: UUID, dogId: UUID, input: UpdateDogInput) async throws
    func deleteDog(ownerId: UUID, dogId: UUID) async throws

    func getPublicDog(dogId: UUID) async throws -> Dog
    func getPrivateDog(ownerId: UUID, dogId: UUID) async throws -> Dog
}