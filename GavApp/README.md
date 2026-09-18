# Локальный запуск iOS-приложения

Из корня репозитория запустите серверы и зависимости:

```sh
cd backend
docker compose up -d --build
```

Debug-сборка обращается к `http://127.0.0.1:8080` (социальный API) и
`http://127.0.0.1:8082` (мессенджер). Эти адреса подходят для iOS Simulator,
когда оба сервиса запущены на том же Mac.

Для физического iPhone в той же Wi-Fi сети укажите IP компьютера при сборке:

```sh
GAV_LOCAL_IP=192.168.0.10 # замените на IP компьютера
xcodebuild -project GavApp/GavApp.xcodeproj -scheme GavApp -configuration Debug \
  -destination 'generic/platform=iOS' \
  GAV_SOCIAL_BASE_URL="http://${GAV_LOCAL_IP}:8080" \
  GAV_MESSENGER_BASE_URL="http://${GAV_LOCAL_IP}:8082" build
```

Серверы должны слушать порты `8080` и `8082`, а приложение и компьютер должны
быть в одной сети. Для Staging и Production в проекте заданы отдельные HTTPS
адреса; их использование требует настроенного DNS и развёрнутых серверов.
