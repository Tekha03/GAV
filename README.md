# GAV

GAV — социальная платформа для владельцев собак. Приложение объединяет профили
пользователей и питомцев, ленту публикаций, поиск компании для прогулок, карту,
трекер прививок и мессенджер.

## Возможности

- профили пользователей и карточки собак;
- лента, публикации, комментарии, лайки и подписки;
- поиск владельцев собак поблизости и совместные прогулки;
- карта с видимыми пользователями и питомцами;
- трекер прививок;
- личные и групповые чаты;
- загрузка изображений и push-уведомления через Firebase при отдельной настройке.

## Технологии и структура

| Каталог | Назначение |
| --- | --- |
| `GavApp/` | iOS/iPadOS-клиент на SwiftUI |
| `backend/social_network/` | основной HTTP/gRPC API на Go |
| `backend/messenger/` | HTTP/gRPC/WebSocket-сервис чатов на Go |
| `backend/api/` | protobuf-контракты и сгенерированный Go-код |
| `backend/docs/` | Swagger/OpenAPI-документация |

Локальное окружение поднимает PostgreSQL, Redis, Kafka, основной API и
мессенджер через Docker Compose.

## Требования

Для полного локального запуска нужны:

- macOS;
- Docker Desktop с поддержкой `docker compose`;
- Xcode с iOS Simulator;
- для запуска без Docker и работы с backend — Go 1.25;
- для генерации protobuf — `protoc`, `protoc-gen-go` и
  `protoc-gen-go-grpc`.

Xcode-проект сейчас имеет deployment target iOS 26.2. Для сборки понадобится
совместимая версия Xcode и симулятор либо устройство с подходящей версией iOS.
Swift Package в `GavApp/Package.swift` использует Swift tools 6.2.

## Быстрый запуск

### 1. Поднять backend

Из корня репозитория:

```sh
cd backend
docker compose up -d --build
```

При первом запуске Docker скачает образы и соберёт два Go-сервиса. Compose
самостоятельно:

- создаст базы `gav` и `gav_social`;
- применит SQL-миграции мессенджера;
- выполнит AutoMigrate и добавит тестовые данные для social network;
- создаст Kafka-топики `chat-events`, `message-events` и `reaction-events`.

Проверить состояние контейнеров:

```sh
docker compose ps
docker compose logs -f social-network messenger
```

После запуска доступны:

| Сервис | Адрес |
| --- | --- |
| Social Network HTTP API | `http://localhost:8080` |
| Swagger UI | `http://localhost:8080/swagger/index.html` |
| Social Network gRPC | `localhost:9000` |
| Messenger HTTP/WebSocket API | `http://localhost:8082` |
| Messenger gRPC | `localhost:9090` |
| PostgreSQL | `localhost:5433` |
| Redis | `localhost:6379` |
| Kafka | `localhost:9092` |

Остановить окружение:

```sh
docker compose down
```

Удалить контейнеры вместе с локальными данными PostgreSQL, Kafka и загруженными
файлами можно командой ниже. Это необратимо удалит данные Docker volumes проекта:

```sh
docker compose down -v
```

### 2. Запустить приложение в Xcode

Откройте проект:

```sh
open GavApp/GavApp.xcodeproj
```

В Xcode:

1. Выберите схему `GavApp`.
2. Выберите iPhone Simulator или подключённый iPhone/iPad.
3. Убедитесь, что backend-контейнеры запущены.
4. Нажмите **Run** (`⌘R`).

Debug-сборка по умолчанию обращается к:

```text
http://Viktorias-MacBook-Pro.local:8080
http://Viktorias-MacBook-Pro.local:8082
```

Это Bonjour-имя Mac, заданное в Build Settings проекта. Если имя вашего Mac
другое, откройте target `GavApp` → **Build Settings** и измените значения
`GAV_SOCIAL_BASE_URL` и `GAV_MESSENGER_BASE_URL` для конфигурации Debug.

Для симулятора можно использовать `http://localhost:8080` и
`http://localhost:8082`. Для физического устройства укажите Bonjour-имя или
локальный IP компьютера, например `http://192.168.1.10:8080`. Mac и устройство
должны находиться в одной сети. При первом запуске разрешите приложению доступ к
локальной сети.

Текущий IP Mac для Wi-Fi обычно можно узнать так:

```sh
ipconfig getifaddr en0
```

Если приложение запускается на физическом устройстве, в target `GavApp` →
**Signing & Capabilities** также выберите свою Development Team.

### Сборка iOS из терминала

Проверить Debug-сборку для симулятора можно без запуска UI:

```sh
xcodebuild \
  -project GavApp/GavApp.xcodeproj \
  -scheme GavApp \
  -configuration Debug \
  -destination 'generic/platform=iOS Simulator' \
  CODE_SIGNING_ALLOWED=NO \
  build
```

Чтобы передать IP backend только для конкретной сборки:

```sh
GAV_LOCAL_IP=192.168.1.10

xcodebuild \
  -project GavApp/GavApp.xcodeproj \
  -scheme GavApp \
  -configuration Debug \
  -destination 'generic/platform=iOS' \
  GAV_SOCIAL_BASE_URL="http://${GAV_LOCAL_IP}:8080" \
  GAV_MESSENGER_BASE_URL="http://${GAV_LOCAL_IP}:8082" \
  build
```

Кроме Debug в проекте есть конфигурации Staging и Production. Они используют
`https://staging-api.gav.app` / `https://staging-messenger.gav.app` и
`https://api.gav.app` / `https://messenger.gav.app`; для них нужны настроенные
DNS и развёрнутые серверы.

## Тесты

### Backend unit-тесты

Go workspace состоит из нескольких модулей, поэтому `go test ./...` из каталога
`backend/` не является корректной общей командой. Запускайте основные сервисы
отдельно:

```sh
(cd backend/social_network && go test ./...)
(cd backend/messenger && go test ./...)
```

Тесты общих библиотек:

```sh
for module in app_errors events ratelimit retry; do
  (cd "backend/shared/${module}" && go test ./...)
done
```

### Интеграционные тесты PostgreSQL и Kafka

Без переменных `GAV_TEST_POSTGRES_DSN` и `GAV_TEST_KAFKA_BROKERS` тесты,
требующие внешней инфраструктуры, будут пропущены. Для полного прогона сначала
поднимите зависимости:

```sh
cd backend
docker compose up -d postgres redis kafka kafka-init
cd ..
```

Затем запустите тесты с адресами локальных сервисов:

```sh
export GAV_TEST_POSTGRES_DSN='postgres://gav:gav@localhost:5433/gav?sslmode=disable'
export GAV_TEST_KAFKA_BROKERS='localhost:9092'

(cd backend/social_network && go test ./...)
(cd backend/messenger && go test ./...)
```

PostgreSQL-тесты создают временные изолированные схемы и удаляют их после
завершения.

### iOS

Отдельного XCTest-target в Xcode-проекте пока нет. На текущем этапе проверка
iOS-клиента — это успешная сборка через Xcode (`⌘B`) или приведённую выше
команду `xcodebuild`. После добавления test target тесты можно будет запускать в
Xcode сочетанием `⌘U` либо через `xcodebuild test`.

## Запуск backend без Docker

Docker Compose — рекомендуемый способ разработки: он гарантирует наличие баз,
миграций, Redis и Kafka. Для ручного запуска всё равно сначала нужны PostgreSQL
и Redis; Kafka можно отключить.

1. Создайте локальные env-файлы:

```sh
cp .env.example backend/social_network/.env
cp .env.example backend/messenger/.env
```

2. Поднимите инфраструктуру и один раз примените инициализацию баз:

```sh
cd backend
docker compose up -d postgres redis social-db-init messenger-migrate
cd ..
```

3. В `backend/social_network/.env` оставьте DSN базы `gav_social`. В
   `backend/messenger/.env` замените `POSTGRES_DSN` на:

```text
postgres://gav:gav@localhost:5433/gav?sslmode=disable
```

4. Запустите сервисы в двух терминалах:

```sh
cd backend/social_network
go run ./cmd/api
```

```sh
cd backend/messenger
go run ./cmd/grpc-server
```

При ручном запуске миграции мессенджера из `backend/messenger/migrations/`
должны быть применены заранее. Команда из шага 2 выполняет их через Compose;
повторно запускать `messenger-migrate` для уже инициализированной базы не нужно.

Если нужна Kafka, запустите `kafka` и `kafka-init`, затем установите
`KAFKA_ENABLED=true` в env-файлах.

## Firebase push-уведомления

По умолчанию Firebase отключён. Для запуска social network с настоящим
service-account JSON:

```sh
cd backend
FIREBASE_SERVICE_ACCOUNT_FILE=/absolute/path/to/firebase-service-account.json \
  docker compose -f docker-compose.yml -f docker-compose.firebase.yml up -d --build
```

Не добавляйте service-account JSON и другие секреты в Git.

## Полезные команды

Пересобрать только backend-сервисы:

```sh
cd backend
docker compose up -d --build social-network messenger
```

Посмотреть последние логи:

```sh
cd backend
docker compose logs --tail=200 social-network messenger
```

Сгенерировать protobuf-код:

```sh
cd backend/api
make proto-all
```

Проверить список схем и конфигураций Xcode:

```sh
xcodebuild -project GavApp/GavApp.xcodeproj -list
```

## Возможные проблемы

### Приложение не подключается к backend

- проверьте `docker compose ps` и логи сервисов;
- проверьте адреса `GAV_SOCIAL_BASE_URL` и `GAV_MESSENGER_BASE_URL`;
- для iPhone используйте IP/Bonjour-имя Mac, а не `localhost`;
- убедитесь, что iPhone и Mac находятся в одной сети;
- разрешите приложению доступ к локальной сети в настройках iOS;
- проверьте, что firewall не блокирует порты `8080` и `8082`.

### Порт уже занят

Проверьте процесс, занявший нужный порт:

```sh
lsof -nP -iTCP:8080 -sTCP:LISTEN
lsof -nP -iTCP:8082 -sTCP:LISTEN
```

### Нужно начать с чистой базы

```sh
cd backend
docker compose down -v
docker compose up -d --build
```

Команда удаляет все локальные данные проекта в Docker volumes.
