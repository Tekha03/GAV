import SwiftUI

struct VaccinationTabView: View {
    @EnvironmentObject private var appViewModel: AppViewModel

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(alignment: .leading, spacing: 18) {
                    if !appViewModel.dogs.isEmpty {
                        Text("Выберите собаку, чтобы открыть прививки и напоминания.")
                            .font(.subheadline)
                            .foregroundStyle(.white.opacity(0.85))
                            .padding(.top, 4)
                    }

                    dogsSection
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 20)
            }
            .background(
                LinearGradient(
                    colors: [
                        Color(red: 0.42, green: 0.22, blue: 0.72),
                        .black
                    ],
                    startPoint: .top,
                    endPoint: .bottom
                )
                .ignoresSafeArea()
            )
            .navigationTitle("Прививки")
            .preferredColorScheme(.dark)
        }
    }

    private var dogsSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Собаки")
                .font(.headline)
                .foregroundStyle(.white)

            if appViewModel.dogs.isEmpty {
                emptyDogsState
            } else {
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: 14) {
                        ForEach(appViewModel.dogs) { dog in
                            NavigationLink {
                                VaccinationListView(dog: dog)
                            } label: {
                                dogCard(dog)
                            }
                            .buttonStyle(.plain)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
        }
    }

    private var emptyDogsState: some View {
        VStack(spacing: 12) {
            Image(systemName: "pawprint")
                .font(.system(size: 38, weight: .semibold))
                .foregroundStyle(.orange)

            Text("Сначала добавьте собаку")
                .font(.headline)
                .foregroundStyle(.white)

            Text("После этого здесь появится трекер прививок и напоминания.")
                .font(.subheadline)
                .multilineTextAlignment(.center)
                .foregroundStyle(.white.opacity(0.68))
                .frame(maxWidth: 280)
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 28)
        .frame(maxWidth: .infinity)
        .background(.white.opacity(0.07), in: RoundedRectangle(cornerRadius: 18))
        .overlay(
            RoundedRectangle(cornerRadius: 18)
                .stroke(.white.opacity(0.08), lineWidth: 1)
        )
    }

    private func dogCard(_ dog: AppDog) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            AsyncImage(url: dog.photoURL) { phase in
                switch phase {
                case .success(let image):
                    image.resizable().scaledToFill()
                default:
                    ZStack {
                        LinearGradient(
                            colors: [
                                Color.white.opacity(0.18),
                                Color.white.opacity(0.06)
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                        Image(systemName: "dog.fill")
                            .font(.title2)
                            .foregroundStyle(.white.opacity(0.9))
                    }
                }
            }
            .frame(width: 148, height: 170)
            .clipped()
            .overlay(
                LinearGradient(
                    colors: [.clear, .black.opacity(0.78)],
                    startPoint: .center,
                    endPoint: .bottom
                )
            )
            .overlay(alignment: .bottomLeading) {
                VStack(alignment: .leading, spacing: 3) {
                    Text(dog.name)
                        .font(.headline.weight(.semibold))
                        .foregroundStyle(.white)
                    Text(dog.breed)
                        .font(.caption)
                        .foregroundStyle(.white.opacity(0.8))
                }
                .padding(10)
            }
            .clipShape(RoundedRectangle(cornerRadius: 22, style: .continuous))
            .shadow(color: .black.opacity(0.22), radius: 10, x: 0, y: 6)
        }
    }
}

struct VaccinationListView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    let dog: AppDog
    @State private var showAdd = false
    @State private var editingItem: AppVaccination?

    var body: some View {
        VStack(alignment: .leading, spacing: 14) {
            VStack(alignment: .leading, spacing: 8) {
                Text(dog.name)
                    .font(.title.bold())
                    .foregroundStyle(.white)

                Text("\(dog.breed) · \(dog.ageText)")
                    .foregroundStyle(.white.opacity(0.75))
            }

            HStack {
                Spacer()
                Button {
                    showAdd = true
                } label: {
                    Label("Добавить", systemImage: "plus")
                }
                .buttonStyle(.borderedProminent)
            }

            let items = appViewModel.vaccinations(for: dog.id)

            if items.isEmpty {
                ContentUnavailableView(
                    "Прививок пока нет",
                    systemImage: "syringe",
                    description: Text("Добавь первую запись для этой собаки.")
                )
                .foregroundStyle(.white)
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                ScrollView {
                    VStack(spacing: 12) {
                        ForEach(items) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.name)
                                        .font(.headline)
                                        .foregroundStyle(.white)

                                    Spacer()

                                    Text(relativeReminderText(for: item.nextDate))
                                        .font(.caption.weight(.semibold))
                                        .foregroundStyle(.orange)
                                }

                                Text("Дата вакцинации: \(item.vaccinationDate.formatted(date: .abbreviated, time: .omitted))")
                                    .font(.footnote)
                                    .foregroundStyle(.white.opacity(0.75))

                                Text(item.notes)
                                    .font(.footnote)
                                    .foregroundStyle(.white.opacity(0.72))

                                Button {
                                    editingItem = item
                                } label: {
                                    Image(systemName: "pencil")
                                        .font(.footnote.weight(.semibold))
                                        .foregroundStyle(.white)
                                        .padding(.horizontal, 10)
                                        .padding(.vertical, 8)
                                        .background(.white.opacity(0.10), in: Capsule())
                                }
                                .buttonStyle(.plain)
                            }
                            .padding(14)
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .background(.white.opacity(0.08), in: RoundedRectangle(cornerRadius: 16))
                        }
                    }
                    .padding(.top, 2)
                }
            }

            Spacer(minLength: 0)
        }
        .padding(20)
        .background(
            LinearGradient(
                colors: [
                    Color(red: 0.42, green: 0.22, blue: 0.72),
                    .black
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()
        )
        .navigationTitle("Прививки")
        .sheet(isPresented: $showAdd) {
            AddVaccinationView(dog: dog)
        }
        .sheet(item: $editingItem) { item in
            AddVaccinationView(dog: dog, editingItem: item)
        }
        .preferredColorScheme(.dark)
    }

    private func relativeReminderText(for date: Date) -> String {
        let days = max(0, Calendar.current.dateComponents([.day], from: Date(), to: date).day ?? 0)
        if days >= 30 { return "\(max(1, days / 30)) мес." }
        if days >= 7 { return "\(max(1, days / 7)) нед." }
        return "\(days) дн."
    }
}

struct AddVaccinationView: View {
    @EnvironmentObject private var appViewModel: AppViewModel
    @Environment(\.dismiss) private var dismiss
    let dog: AppDog
    let editingItem: AppVaccination?

    @State private var name: String
    @State private var vaccinationDate: Date
    @State private var reminderYears: Int
    @State private var reminderWeeks: Int
    @State private var reminderDays: Int
    @State private var notes: String

    init(dog: AppDog, editingItem: AppVaccination? = nil) {
        self.dog = dog
        self.editingItem = editingItem

        _name = State(initialValue: editingItem?.name ?? "")
        _vaccinationDate = State(initialValue: editingItem?.vaccinationDate ?? .now)
        _reminderYears = State(initialValue: 0)
        _reminderWeeks = State(initialValue: 0)
        _reminderDays = State(initialValue: editingItem.map { $0.reminderAfterDays } ?? 0)
        _notes = State(initialValue: editingItem?.notes ?? "")
    }

    var body: some View {
        NavigationStack {
            Form {
                Section("Прививка") {
                    TextField("Название", text: $name)
                    DatePicker("Дата вакцинации", selection: $vaccinationDate, displayedComponents: .date)
                    Text("Интервал следующей вакцинации")
                    timerColumns()
                    TextField("Заметки", text: $notes, axis: .vertical)
                }

                Section {
                    Button("Сохранить") {
                        save()
                    }
                    .disabled(name.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty)
                }
            }
            .navigationTitle(dog.name)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Отмена") { dismiss() }
                }
            }
        }
    }

    private func timerColumns() -> some View {
        HStack(spacing: 12) {
            timerColumn(title: "Годы", selection: $reminderYears, range: 0...5)
            timerColumn(title: "Недели", selection: $reminderWeeks, range: 0...52)
            timerColumn(title: "Дни", selection: $reminderDays, range: 0...30)
        }
    }

    private func timerColumn(title: String, selection: Binding<Int>, range: ClosedRange<Int>) -> some View {
        VStack(spacing: 8) {
            Text(title)
                .font(.caption)
                .foregroundStyle(.secondary)

            Picker(title, selection: selection) {
                ForEach(range, id: \.self) { value in
                    Text("\(value)").tag(value)
                }
            }
            .pickerStyle(.wheel)
            .frame(width: 90, height: 120)
            .clipped()
        }
        .frame(maxWidth: .infinity)
    }

    private func save() {
        var date = vaccinationDate
        date = Calendar.current.date(byAdding: .year, value: reminderYears, to: date) ?? date
        date = Calendar.current.date(byAdding: .weekOfYear, value: reminderWeeks, to: date) ?? date
        date = Calendar.current.date(byAdding: .day, value: reminderDays, to: date) ?? date

        let item = AppVaccination(
            id: editingItem?.id ?? UUID(),
            dogID: dog.id,
            name: name.trimmingCharacters(in: .whitespacesAndNewlines),
            vaccinationDate: vaccinationDate,
            reminderAfterDays: reminderYears * 365 + reminderWeeks * 7 + reminderDays,
            nextDate: date,
            notes: notes.isEmpty ? "Без заметок" : notes
        )

        if let editingItem {
            if let index = appViewModel.vaccinations.firstIndex(where: { $0.id == editingItem.id }) {
                appViewModel.vaccinations[index] = item
            }
        } else {
            appViewModel.vaccinations.append(item)
        }

        dismiss()
    }
}
