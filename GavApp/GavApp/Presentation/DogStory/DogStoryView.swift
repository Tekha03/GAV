//
//  DogStoryView.swift
//  
//
//  Created by Виктория Кашуркина on 01.04.2026.
//

import SwiftUI

struct DogStoryView: View {
    let dog: Dog

    var body: some View {
        VStack(spacing: 6) {

            Image(dog.photoUrl)
                .resizable()
                .scaledToFill()
                .frame(width: 60, height: 60)
                .clipShape(Circle())
                .overlay(
                    Circle().stroke(Color(red: 220/255, green: 255/255, blue: 5/255), lineWidth: 2.5)
                )

            Text(dog.name)
                .font(.caption)
                .foregroundColor(.white)
        }
        .frame(width: 70)
    }
}
