//
//  DogPresentationMapper.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 01.04.2026.
//


enum DogPresentationMapper {

    static func gender(_ gender: String) -> String {
        switch gender.lowercased() {
        case "male": return "кобель"
        case "female": return "сука"
        default: return gender
        }
    }

    static func breed(_ breed: String) -> String {
        breed
    }

    static func character(status: String, gender: String) -> String {
        switch status.lowercased() {
        case "friendly":
            return gender.lowercased() == "female"
                ? "дружелюбная"
                : "дружелюбный"
        case "cautious":
            return "осторожная"
        default:
            return status
        }
    }
}