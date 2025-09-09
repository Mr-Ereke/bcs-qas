# Quotes Alerting Service

## Установка
### 1. Клонируем репозиторий в необходимый нам каталог
```bash
cd path_to_services_directory
git clone ***/quotes-alerting-service.git quotes-alerting-service
```

### 2. Переходим в каталог с сервисом, получаем изменения и переключаемся на необходимую версию
```bash
cd quotes-alerting-service
git fetch --all && git checkout 1.0.0 или другая актуальная версия.
```
[Версии смотреть здесь:](***)

### 3. Копируем конфигурационные файлы с переменными среды и прописываем в нем необходимые конфиги
```bash
cp src/config/config-example.conf src/config/config.conf
cp src/config/env-example.conf src/config/env.conf
```

### 4. Обслуживание
Команда `./service.sh` имеет 5 параметра:

- build - сборка образа
- start - создаёт и запускает контейнер с сервисом если он не запущен.
- stop - останавливает контейнер.
- restart - останавливает, пересоздаёт и запускает контейнер.
- rebuild - останавливает, удаляет контейнер, если он есть и пересобирает образ.

```bash
./service.sh build
./service.sh start
./service.sh stop
./service.sh restart
./service.sh rebuild
```

#### После любых изменений файла с конфигурацией - src/config/config.conf необходимо:
- пересобрать сервис командой ./service.sh rebuild
- запустить заново сервис командой ./service.sh start

### С подробной документацией по сервису можно ознакомиться в Confluence
[Confluence документация сервиса](***)

[Confluence ТЗ для сервиса](***)
