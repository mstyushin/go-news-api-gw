GoNews API Gateway
==================
Точка входа в бэкенд агрегатора новостей.

Умеет отдавать пагинированный список новостей, конкретную новость по её ID и создавать комментарий под новостью.

---

Внутренние сервисы:

- [go-news-scraper](https://github.com/mstyushin/go-news-scraper)
- [go-news-comments](https://github.com/mstyushin/go-news-comments)
- [go-news-moderation](https://github.com/mstyushin/go-news-moderation)

---

Доступны следующие эндпойнты:

- `GET /news/latest?page={num}&page_size={size}&s={search_string}` - отдаёт новости списком, отсортированным по дате создания (дата берётся в 
том виде, как её вернул тот или иной RSS фид). При этом текст каждой новости сокращён (по-умолчанию до 100 символов, на стороне сервиса 
go-news-scraper). В каждом объекте статьи из списка имеется заполненное поле **link_to_full** - ссылка, по которой можно отправить HTTP GET 
и получить от API Gateway сервиса детальный объект статьи. Поддерживает параметры пагинации **page** - номер страницы и **page_size** - число 
статей в одной странице. Также поддерживает параметр **s** - строка для поиска по заголовкам статей.
Пример ответа:

        {
            "articles":[
                {
                    "id": 1,
                    "title": "some title",
                    "short_content": "something for nothing...",
                    "link_to_full": "http://127.0.0.1:8080/news/1",
                    "pub_time": 1728683086
                },
                {
                    "id": 45,
                    "title": "another title",
                    "short_content": "yet another something...",
                    "link_to_full": "http://127.0.0.1:8080/news/45",
                    "pub_time": 1728684076
                },
            ],
            "paginator": {
                    "num_pages": 8,
                    "cur_page": 1,
                    "page_size": 2
            }
        }

- `GET /news/{ID}` - отдаёт детальный объект статьи с идентификтором ID + все комментарии к ней (если они есть). Объекты комментариев имеют поле **parent_id** - если 
оно равно значению в поле **id**, то это просто комментарий к новости, в общем случае в **parent_id** хранится идентификатор *комментируемого* комментария, т.е. это 
поле позволяет построить дерево комментариев.
Пример ответа:

        {
            "article": {
                "id": 23,
                "title": "some title",
                "content": "something really viral",
                "link": "https://it.slashdot.org/story/24/10/12/1943206/microsofts-take-on-kernel-access-and-safe-deployment-after-crowdstrike-incident",
                "rss_feed_id": 1,
                "pub_time": 1728673086
            },
            "comments": [
                {
                    "id": 1,
                    "article_id": 23,
                    "parent_id": 1,
                    "author": "Alice",
                    "text": "hey there!",
                    "pub_time": 1728693086
                }
            ]
        }

- `POST /news/{ID}?c=true` - принимает в теле запроса объект комментария к новости с идентификатором ID. Ожидает обязательный query parameter `?c=true`, без него 
ответит 400й ошибкой. В случае, если комментарий прошёл модерацию, то в ответе будет идентификатор созданного комментария с unix timestamp и код ответа 200. 
В противном случае сервис ответит лишь 400й ошибкой.
Пример запроса:

        {
            "author": "Bob",
            "text": "bla bla bla"
        }

    Пример ответа:

        {
            "id": 1,
            "pub_time": 1728769897
        }

Все эндпойнты поддерживают query parameter **request_id** - если он присутствует, то запросы к внутренним сервисам будут дополнятся этим же параметром + в заголовки запросов будет 
добавляться **X-Request-Id** с данным значением. В случае когда **request_id** не указан, то он будет сгенерирован и также передан внутренним сервисам.
Ответ от сервиса API Gateway *всегда* будет содержать заголовок **X-Request-Id** содержащий либо переданный **request_id** либо сгенерированное значение.

# Конфигурационные параметры

| Параметр                     | Описание                                       | Значение по-умолчанию   |
|------------------------------|------------------------------------------------|-------------------------|
| `http_port`                  | порт API Gateway сервиса                       | `8080`                  |
| `base_url`                   | базовый URL по которому доступен API           | `http://127.0.0.1:8080` |
| `news_scraper_address`       | IP или имя хоста с сервисом go-news-scraper    | `127.0.0.1`             |
| `news_scraper_port`          | порт, на котором слушает go-news-scraper       | `8081`                  |
| `comments_service_address`   | IP или имя хоста с сервисом go-news-comments   | `127.0.0.1`             |
| `comments_service_port`      | порт, на котором слушает go-news-comments      | `8082`                  |
| `moderation_service_address` | IP или имя хоста с сервисом go-news-moderation | `127.0.0.1`             |
| `moderation_service_port`    | порт, на котором слушает go-news-moderation    | `8083`                  |

# Сборка и запуск

## Требования

-   golang 1.22

-   docker >=23.0

---

Просмотр конфигурации по-умолчанию:

        $ make build
        $ ./bin/go-news-api-gw -print-config

---

Для быстрого запуска с конфигом по-умолчанию:

        $ make run

Сервер будет запущен на `127.0.0.1:8080`. 

---

Запуск сервера с конфигурацией из файла `config.yaml`:

        $ ./bin/go-news-api-gw -config config.yaml

---

Остановить сервер:

        $ make clean

---

Показать версию сборки:

        $ ./bin/go-news-api-gw -version

## Примеры тестовых запросов

        $ curl "http://127.0.0.1:8080/news/latest"
        $ curl -v -XPOST "http://127.0.0.1:8080/news/1?c=true" -d '{"Author": "Bob", "Text": "bla bla bla"}'
        $ curl -v -XPOST "http://127.0.0.1:8080/news/345?c=true" -d '{"Author": "Alice", "Text": "something qwerty"}'
        $ curl "http://127.0.0.1:8080/news/latest?page_size=3&page=29"
        $ curl "http://127.0.0.1:8080/news/latest?s=TCP&page_size=10&page=1"

## Docker
TODO
