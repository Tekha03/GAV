import SwiftUI

@available(macOS 12.0, *)
struct AddDogView: View {
    let viewModel: ProfileViewModel
    @Environment(\\\.dismiss) private var dismiss

    @State private var name: String = ""
    @State private var breed: String = ""
    @State private var age: String = ""
    @State private var description: String = ""

    var body: some View {
        NavigationStack {
            Form {
                Section("Информация о собаке") {
                    TextField("Имя", text: $name)
                    TextField("Порода", text: $breed)
                    TextField("Возраст", text: $age)
                    TextField("Описание", text: $description)
                }

                Button("Создать собаку") {
                    Task {
                        let input = CreateDogInput(
                            name: name,
                            breed: breed,
                            age: Int(age) ?? 0,
                            description: description
                        )
                        // TODO: вызвать dogUseCase.create
                        //  viewModel.dogs.append(newDog)
                        dismiss()
                    }
                }
            }
            .navigationTitle("Добавить собаку")
        }
    }
}