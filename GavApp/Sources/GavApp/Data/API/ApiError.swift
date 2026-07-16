import Foundation

enum APIError: Error, LocalizedError, CustomNSError {
    case invalidURL
    case invalidResponse(statusCode: Int)
    case decodingError(Error)
    case networkError(Error)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Некорректный адрес сервера"

        case .invalidResponse(let statusCode):
            return responseManager(for: statusCode)

        case .decodingError(_):
            return "Не удалось обработать ответ сервера"

        case .networkError(let error):
            return error.localizedDescription
        }
    }

    private func responseManager(for statusCode: Int) -> String {
        switch statusCode {
        case 400:
            return "Некорректный запрос"

        case 401:
            return "Необходима авторизация"

        case 403:
            return "Недостаточно прав для выполнения действия"

        case 404:
            return "Запрошенные данные не найдены"

        case 409:
            return "Такие данные уже существуют"

        case 422:
            return "Проверьте введённые данные"

        case 500...599:
            return "Сервис временно недоступен"

        default:
            return "Ошибка сервера: \(statusCode)"
        }
    }
}
