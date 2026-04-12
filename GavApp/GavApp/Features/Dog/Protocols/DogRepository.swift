//
//  DogRepository.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 10.04.2026.
//


import Foundation

protocol DogRepository {
    func create(_ dog: Dog) async throws -> Dog
    func update(_ dog: Dog) async throws
    func delete(id: UUID) async throws

    func getById(_ id: UUID) async throws -> Dog
    func getByOwnerId(_ ownerId: UUID) async throws -> [Dog]

    func getPublicDog(id: UUID) async throws -> Dog
}