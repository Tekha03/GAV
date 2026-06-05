import Foundation

enum APIError: Error, LocalizedError, CustomNSError {
    case invalidURL
    case invalidResponse(statusCode: Int)
    case serverError(statusCode: Int, code: String?, message: String)
    case decodingError(Error)
    case networkError(Error)
    case userMessage(String)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Некорректный адрес сервера"
        case .invalidResponse(let code):
            return Self.message(forStatusCode: code)
        case .serverError(let statusCode, let code, let message):
            return Self.message(forStatusCode: statusCode, code: code, message: message)
        case .decodingError:
            return "Сервер вернул данные в неожиданном формате"
        case .networkError(let error):
            return Self.networkMessage(for: error)
        case .userMessage(let message):
            return message
        }
    }

    var statusCode: Int? {
        switch self {
        case .invalidResponse(let statusCode):
            return statusCode
        case .serverError(let statusCode, _, _):
            return statusCode
        default:
            return nil
        }
    }

    private static func message(forStatusCode statusCode: Int) -> String {
        switch statusCode {
        case 400:
            return "Проверьте введенные данные"
        case 401:
            return "Неверный email или пароль"
        case 404:
            return "Профиль не найден"
        case 409:
            return "Такие данные уже используются"
        case 500...599:
            return "Сервер временно недоступен. Попробуйте еще раз"
        default:
            return "Не удалось выполнить запрос. Код ответа: \(statusCode)"
        }
    }

    private static func message(forStatusCode statusCode: Int, code: String?, message: String) -> String {
        let lowercasedMessage = message.lowercased()

        if statusCode == 401 || code == "UNAUTHORIZED" || lowercasedMessage.contains("invalid email or password") {
            return "Неверный email или пароль"
        }

        if statusCode == 409 {
            if lowercasedMessage.contains("profile") || lowercasedMessage.contains("username") {
                return "Этот никнейм уже занят"
            }
            return "Пользователь с таким email уже существует"
        }

        if statusCode == 400 || code == "VALIDATION_ERROR" {
            if lowercasedMessage.contains("username") {
                return "Никнейм должен быть от 3 до 30 символов: латинские буквы, цифры, _ или ."
            }
            return "Проверьте введенные данные"
        }

        if statusCode == 404 || code == "NOT_FOUND" {
            return "Профиль не найден"
        }

        if statusCode >= 500 {
            return "Сервер временно недоступен. Попробуйте еще раз"
        }

        return message.isEmpty ? Self.message(forStatusCode: statusCode) : message
    }

    private static func networkMessage(for error: Error) -> String {
        guard let urlError = error as? URLError else {
            return "Не удалось подключиться к серверу. Проверьте соединение"
        }

        switch urlError.code {
        case .notConnectedToInternet, .networkConnectionLost:
            return "Нет соединения с интернетом"
        case .cannotFindHost, .cannotConnectToHost, .dnsLookupFailed:
            return "Не удалось подключиться к серверу"
        case .timedOut:
            return "Сервер не ответил вовремя. Попробуйте еще раз"
        default:
            return "Не удалось подключиться к серверу. Проверьте соединение"
        }
    }
}
