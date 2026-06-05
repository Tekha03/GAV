import SwiftUI

@available(macOS 12.0, *)
struct EditDogView: View {
    @Environment(\.dismiss) private var dismiss

    @State private var name: String
    @State private var photoURL: String
    @State private var breed: String
    @State private var age: String
    @State private var status: String
    @State private var gender: String

    let onSave: (Dog) -> Void
    let originalDog: Dog

    init(dog: Dog, onSave: @escaping (Dog) -> Void) {
        self.originalDog = dog
        self.onSave = onSave

        _name = State(initialValue: dog.name)
        _photoURL = State(initialValue: dog.photoURL)
        _breed = State(initialValue: dog.breed)
        _age = State(initialValue: String(describing: dog.age))
        _status = State(initialValue: dog.status.rawValue)
        _gender = State(initialValue: dog.gender.rawValue)
    }

    var body: some View {
        NavigationStack {
            Form {
                Section("Основное") {
                    TextField("Имя", text: $name)
                    TextField("Фото URL", text: $photoURL)
                    TextField("Порода", text: $breed)
                    TextField("Возраст", text: $age)
                }

                Section("Параметры") {
                    TextField("Статус", text: $status)
                    TextField("Пол", text: $gender)
                }
            }
            .navigationTitle("Редактировать пса")
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена") { dismiss() }
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("Сохранить") {
                        guard let ageValue = DogAge(rawValue: age) else { return }

                        let updatedDog = Dog(
                            id: originalDog.id,
                            ownerId: originalDog.ownerId,
                            name: name,
                            breed: breed,
                            photoURL: photoURL,
                            status: DogStatus(rawValue: status) ?? originalDog.status,
                            age: ageValue,
                            gender: DogGender(rawValue: gender) ?? originalDog.gender
                        )

                        onSave(updatedDog)
                        dismiss()
                    }
                }
            }
        }
    }
}