//
//  DogProfileHeroView.swift
//  GavApp
//
//  Created by Виктория Кашуркина on 09.04.2026.
//


import SwiftUI

struct DogProfileHeroView: View {
    let dog: Dog

    var body: some View {
        Image(dog.photoUrl)
            .resizable()
            .scaledToFill()
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .clipped()
            .ignoresSafeArea()
    }
}