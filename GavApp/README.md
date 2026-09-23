# Локальный запуск iOS-приложения

Из корня репозитория запустите серверы и зависимости:

```sh
cd backend
docker compose up -d --build
```

Debug-сборка обращается к серверам на Mac по Bonjour-имени
`Viktorias-MacBook-Pro.local`. Это работает как в iOS Simulator, так и на физическом
iPhone в той же Wi-Fi сети. При первом обращении iOS запросит доступ к
локальной сети. Если Bonjour-имя Mac изменится, обновите
`GAV_SOCIAL_BASE_URL` и `GAV_MESSENGER_BASE_URL` в Debug Build Settings.

Если mDNS/Bonjour в сети не работает, можно явно указать IP компьютера при сборке:

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
