//
//  DogProfileMetaView.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 09.04.2026.
//


import SwiftUI

struct DogProfileMetaView: View {
    let dog: Dog

    var body: some View {
        VStack(spacing: 8) {

            Text(dog.name)
                .font(.largeTitle.bold())
                .foregroundStyle(.white)
                .multilineTextAlignment(.center)

            Text("\(DogPresentationMapper.breed(dog.breed)) • \(DogPresentationMapper.gender(dog.gender)) • \(dog.age) года")
                .font(.subheadline)
                .foregroundStyle(.white.opacity(0.95))
                .multilineTextAlignment(.center)

            Text(DogPresentationMapper.character(status: dog.status, gender: dog.gender))
                .font(.headline)
                .foregroundStyle(.white.opacity(0.95))
                .multilineTextAlignment(.center)
        }
        .frame(maxWidth: .infinity)
        .padding(.horizontal, 24)
    }
}