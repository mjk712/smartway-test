[![CI](https://github.com/mjk712/smartway-test/actions/workflows/ci.yaml/badge.svg)](https://github.com/mjk712/smartway-test/actions/workflows/ci.yaml)

[![Coverage Status](https://coveralls.io/repos/github/mjk712/smartway-test/badge.svg?branch=main)](https://coveralls.io/github/mjk712/smartway-test?branch=main)

Web-сервис для работы с хранилищем данных перелётов

Сборка: 

docker-compose up

- Сервис настраивается через переменные окружения:
SERVER_ADDRESS
POSTGRES_CONN
POSTGRES_USERNAME
POSTGRES_PASSWORD
POSTGRES_HOST
POSTGRES_PORT
POSTGRES_DATABASE
ENV

Доступный функционал:
⁃ Чтение списка билетов
⁃ Чтение списка пассажиров по билету
⁃ Чтение списка документов по пассажиру
⁃ Чтение полной информации по билету, включая информацию о билете,
пассажирах и документе сразу
⁃ Редактирование информации о билете
⁃ Редактирование информации о пассажире
⁃ Редактирование информации о документе
⁃ Удаление информации о билете
⁃ Удаление информации о пассажире
⁃ Удаление информации о документе
⁃ Получение отчёта по пассажиру за определённый период времени.

Подробное описание можно посмотреть в swagger