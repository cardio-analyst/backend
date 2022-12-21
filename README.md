## Подготовка

Сервис состоит из двух компонентов: API-сервер и база данных PostgreSQL, поэтому для его успешной работы необходимо
установить следующее ПО:

- [Go](https://golang.org/doc/install) >=1.17;
- [Docker](https://www.docker.com/get-started) >=20.10.14
- [Git](https://git-scm.com/)

## Запуск

После установки ПО необходимо скачать исходный код сервиса и перейти в директорию с исходным кодом:

```shell
git clone https://github.com/cardio-analyst/backend.git
cd backend
```

Автоматизация запуска системы осуществляется средствами утилиты Makefile и файлов docker-compose.yml, config.yaml. Перед запуском рекомендуется обратить внимание на переменные, указанные в этих файлах, и изменить их, если потребуется.

Переменные среды окружения (см. [docker-compose.yml](https://github.com/cardio-analyst/backend/blob/main/docker-compose.yml)):

```
PORT - адрес, на котором будет запущен сервер
DATABASE_URL - адрес базы данных
ACCESS_TOKEN_SIGNING_KEY - ключ подписи ACCESS TOKEN
REFRESH_TOKEN_SIGNING_KEY - ключ подписи REFRESH TOKEN
ACCESS_TOKEN_TTL_SEC - время жизни ACCESS TOKEN в секундах
REFRESH_TOKEN_TTL_SEC - время жизни REFRESH TOKEN в секундах
SMTP_PASSWORD - пароль от электронной почты для отправки отчётов (используется SMTP-сервер Yandex)
```

### [Docker Compose](https://docs.docker.com/compose/gettingstarted/)

Как было упомянуто выше, система запускается с помощью Docker. Оба компонента системы (API сервер и БД) разворачиваются
в отдельных Docker-контейнерах. Настройки компонентов указываются в
[docker-compose.yml](https://github.com/cardio-analyst/backend/blob/main/docker-compose.yml).

Чтобы запустить систему, необходимо ввести следующую команду:

```shell
# запуск компонентов в отдельных Docker-контейнерах
make compose-up # либо docker-compose up
```

## Эндпойнты

После успешного запуска сервис будет доступен по адресу `http://localhost:8080`.

Ниже описаны возможности RESTful API-сервера.

### Авторизация

* `POST /api/v1/auth/signUp`: регистрация
* `POST /api/v1/auth/signIn`: авторизация
* `POST /api/v1/auth/refreshTokens`: обновление ACCESS TOKEN

### Профиль пользователя

* `GET /api/v1/profile/info`: получить информацию о профиле
* `PUT /api/v1/profile/edit`: изменить информацию о профиле

### Общие показатели здоровья

* `GET /api/v1/diseases/info`: получить информацию об общих показателях здоровья пациента
* `PUT /api/v1/diseases/edit`: изменить информацию об общих показателях здоровья пациента

### Анализы

* `GET /api/v1/analyses`: получить все записи лабораторных и инструментальных исследований пациента
* `POST /api/v1/analyses`: создать новую запись лабораторных и инструментальных исследований пациента
* `PUT /api/v1/analyses/{analysisID}`: изменить запись лабораторных и инструментальных исследований пациента по её идентификатору (analysisID)

### Образ жизни

* `GET /api/v1/lifestyles/info`: получить информацию об образе жизни пользователя
* `PUT /api/v1/lifestyles/edit`: изменить информацию об образе жизни пользователя

### Базовые показатели

* `GET /api/v1/basicIndicators`: получить все записи базовых показателей сердечно-сосудистого здоровья пациента
* `POST /api/v1/basicIndicators`: создать новую запись базовых показателей сердечно-сосудистого здоровья пациента
* `PUT /api/v1/basicIndicators/{basicIndicatorsID}`: изменить запись базовых показателей сердечно-сосудистого здоровья пациента по её идентификатору (basicIndicatorsID)

### SCORE

* `GET /api/v1/score/cveRisk?gender={gender}&smoking={statusSmoking}&sbpLevel={sbpLevel}&totalCholesterolLevel={totalCholesterolLevel}`: получить информацию о риске сердечно-сосудистых событий в течение 10 лет по шкале SCORE
* `GET /api/v1/score/idealAge?gender={gender}&smoking={statusSmoking}&sbpLevel={sbpLevel}&totalCholesterolLevel={totalCholesterolLevel}`: получить идеальный «сердечно-сосудистый возраст»

### Рекомендации

* `GET /api/v1/recommendations`: получить список автоматически сгенерированных пациент-ориентированных рекомендаций
* `POST /api/v1/recommendations/send`: отправить сформированный отчёт на электронную почту
