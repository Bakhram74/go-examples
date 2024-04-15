# Single Window Backend

## Конфигурирование

Настройки конфигурации сначала считываются из файла `config/config.yml`, затем перезаписываются из переменных окружения.

Например, если в `config.yml` прописать `pg.poolMax = 2`, а затем в файле `.env` прописать `PG_POOLMAX=1`, то запись из `.env` файла перезапишет запись из `config.yml` и по итогу `poolMax` будет равно 1.

Примечание: `.env` файл не обязателен и должен использоваться исключительно в *local* режиме. Если приложение будет запущено в *prod* или *stage* режиме и в главной директории приложения будет лежать `.env` файл, то приложение упадет с ошибкой.

## Запуск

1. Для начала запускаем контур с помощью команды:

`docker compose up --build` - Linux
`docker-compose up --build` - MacOS

Примечание: при любом изменении в коде необходимо выполнить rebuild сервиса с помощью команды ниже.

`docker-compose up -d --no-deps --force-recreate --build {{container_name}}`

Примечание: при возникновении ошибки `Version in "./docker-compose.yml" is unsupported. You might be seeing this error because you're using the wrong Compose file version. Either specify a supported version (e.g "2.2" or "3.3") and place your service definitions under the 'services' key, or omit the 'version' key and place your service definitions at the root of the file to use version 1.` обновите docker и docker-compose.

2. После запуска, можно посмотреть спецификацию по `http://localhost:8080/`

Примечание: при запуске всего приложения, после успешной миграции вызывается код, который запускает скрипты из папки `sql` в алфавитном порядке. Планируется, что там будут находиться скрипты для заполнения базы тестовыми данными.

## Open API спецификация

Для генерации потребуется Redocly CLI `npm install -g @redocly/cli` так же потребуется `Node.js` для корректной работы `npm`

При внесении каких-то изменений в Open API спецификацию, необходимо перегенерировать модели, с помощью команды:

`sh generate-openapi.sh`

### Подключение авторизации

Для подключения авторизации необходимо добавлять следующую запись в конец **responses** (пример - openapi/routes/dispute/Dispute.yaml):

```
security: 
    - SWCookieAuth: []
```

## Тестирование

Для моков используется [mock](https://github.com/uber-go/mock/tree/main) от Uber

1. После добавления нового функционала(роутера, репозитория, юз кейса), необходимо добавить интерфейсы в `internal/usecase/interfaces.go`
2. Генерируем моки с помощью скрипта `generate-mocks.sh`
3. Запускаем тесты с помощью команды `go test -v -cover -race ./internal/...`

Примечание: для скрипта `generate-mocks.sh` требуется утилита mockgen, установить ее можно с помощью команды `go install go.uber.org/mock/mockgen@latest`

## Миграции

Миграции баз данных — это контролируемый процесс изменения схемы базы данных, позволяющий последовательно вносить и отслеживать изменения структуры данных.

Примечание: механизм миграций лежит в `internal/app/migrate.go` и запускает скрипты из `migrations/`

### Перезапуск миграций

1. Необходимо остановить работу docker стека `backend`.
2. Далее необходимо удалить "Container" `postgres`
3. Следущим шагом удаляем "Volume" `backend_pg_data`
4. После чего выполняем команду `docker-compose up` для запуска приложения.

## Авторизация
Для авторизации используется API Auth v3, подробнее почитать о нем можно [здесь](https://youtrack.wildberries.ru/articles/AUTH-A-28/Bystryj-start).

Примечание: для авторизации можно использовать как stage (`https://auth-stage-single-window.wildberries.ru`), так и prod (`https://auth-single-window.wildberries.ru`) URL. Чтобы протестировать нужный вам URL авторизации необходимо поменять переменную *AUTHCLIENT_WBURL* в файле *.env*.

## Линты
В качестве линтера используется набор линтов [`golangci-lint`](https://golangci-lint.run/usage/linters/)  

Для проверки по линтам необходимо установить последнюю версию `golangci`. Текущая версия линтера `v1.55.2`.

Проверка линтов c возможными исправлениями по ним осуществляется запуском скрипта `run-linter.sh`.  
Проверка идет по файлу `golangci.yml`, включенные линты описаны ниже.

### Включенные

| Линт                      | Описание                                                                            | Ссылки (если есть)                                      |
|---------------------------|-------------------------------------------------------------------------------------|---------------------------------------------------------|
| errcheck                  | Проверка на возврат и обработку ошибок                                              | https://golangci-lint.run/usage/linters/#errcheck       |
| gosimple                  | Проверка того, что связано с наименованиями в го                                    | https://golangci-lint.run/usage/linters/#gosimple       |
| ineffassign               | Обнаруживает, когда назначения существующим переменным не используются              | https://github.com/gordonklaus/ineffassign              |
| staticcheck               | Указано в документации                                                              | https://staticcheck.dev/docs/checks/                    |
| unused                    | Проверка на неиспользуемые переменные                                               |                                                         |
| asciicheck                | Проверка на использование только ascii символов в коде                              |                                                         |
| bidichk                   | Проверка на опасные unicode последовательности                                      |                                                         |
| bodyclose                 | Проверка на закрытие Body                                                           |                                                         |
| containedctx              | Проверка на отсутствие контекста в структуре                                        | https://go.dev/blog/context-and-structs                 |
| decorder                  | Проверка порядка декларации переменных                                              |                                                         |
| dogsled                   | Проверка на большое количество пустых возвращаемых значений                         | https://golangci-lint.run/usage/linters/#dogsleed       |
| dupl                      | Проверка на дублирующийся код                                                       | https://golangci-lint.run/usage/linters/#dupl           |
| dupword                   | Проверка на дублирующиеся слова в коде                                              | https://golangci-lint.run/usage/linters/#dupword        |
| durationcheck             | Проверка на умножение временных величин                                             |                                                         |
| errchkjson                | Проверка на передачу типов в обработку в пакет json                                 | https://golangci-lint.run/usage/linters/#errchkjson     |
| errname                   | Наименование ошибок в стандарте go                                                  |                                                         |
| errorlint                 | Проверка на корректное оборачивание ошибок                                          | https://golangci-lint.run/usage/linters/#errorlint      |
| exhaustive                | Проверка на корректную работу с перечислениями                                      | https://golangci-lint.run/usage/linters/#exhaustive     |
| exportloopref             | Проверяет наличие указателей на включающие переменные цикла.                        |                                                         |
| forbidigo                 | Запрещает кастомные идентификаторы импортов                                         | https://github.com/ashanbrown/forbidigo                 |
| funlen                    | Регулирует длину функций                                                            | https://golangci-lint.run/usage/linters/#funlen         |
| ginkgolinter              | См. документацию                                                                    | https://github.com/nunnatsa/ginkgolinter                |
| gocheckcompilerdirectives | Проверка на валидность директив компилятора в go                                    |                                                         |
| gochecksumtype            |                                                                                     |                                                         |
| goconst                   | Замена дубликатов строк константой                                                  |                                                         |
| gocritic                  | Проверка кода на типичные баги                                                      | https://golangci-lint.run/usage/linters/#gocritic       |
| gocyclo                   | Проверка цикломатической сложности функции                                          | https://golangci-lint.run/usage/linters/#gocyclo        |
| gofmt                     | Заставляет использовать наиболее современные конструкции в go                       | https://golangci-lint.run/usage/linters/#gofmt          |
| gofumpt                   | Проверка (адекватной сложности) порядка импортов                                    | https://golangci-lint.run/usage/linters/#gofumpt        |
| gomnd                     | Проверка на магические числа                                                        | https://golangci-lint.run/usage/linters/#gomnd          |
| gosec                     | Линты связанные с безопасностью                                                     | https://golangci-lint.run/usage/linters/#gosec          |
| gosmopolitan              | Проверка i18n/l10n                                                                  | https://golangci-lint.run/usage/linters/#gosmopolitan   |
| grouper                   | См. доку                                                                            | https://golangci-lint.run/usage/linters/#grouper        |
| inamedparam               | Параметры в интерфейсах должны быть именованными                                    |                                                         |
| interfacebloat            | Защита от разрастания интерфейса                                                    | https://golangci-lint.run/usage/linters/#interfacebloat |
| lll                       | Защита от длинных строк                                                             | https://golangci-lint.run/usage/linters/#lll            |
| loggercheck               | Проверка kv для кастомных полей                                                     | https://golangci-lint.run/usage/linters/#loggercheck    |
| maintidx                  | maintainability index                                                               | https://golangci-lint.run/usage/linters/#maintidx       |
| makezero                  | Позволяет инициализировать только слайсы с нулевой длинной                          | https://golangci-lint.run/usage/linters/#makezero       |
| misspell                  | Опечатки                                                                            | https://golangci-lint.run/usage/linters/#misspell       |
| musttag                   | Все структуры которые маршалятся и анмаршалятся должны иметь тег                    | https://golangci-lint.run/usage/linters/#musttag        |
| nakedret                  | Запрещает [naked return](https://go.dev/tour/basics/7)                              | https://golangci-lint.run/usage/linters/#nakedret       |
| nestif                    | Вложенные if                                                                        | https://golangci-lint.run/usage/linters/#nestif         |
| nilerr                    | Срабатывает когда при ненулевой ошибки функция возвращает нулевую ошибку            |                                                         |
| nolintlint                | Срабатывает при неверном использовании директивы nolint                             | https://golangci-lint.run/usage/linters/#nolintlint     |
| perfsprint                | Замена sprintf более подходящей альтернативой                                       | https://golangci-lint.run/usage/linters/#perfsprint     |
| prealloc                  | Преаллокация слайса                                                                 | https://golangci-lint.run/usage/linters/#prealloc       |
| predeclared               | Нужно для проверки предопределенных переменных для избежания затенения              | https://golangci-lint.run/usage/linters/#predeclared    |
| sloglint                  | Проверка логов согласно стандартам                                                  | https://golangci-lint.run/usage/linters/#sloglint       |
| tenv                      | При тесте чтобы использовалась переменная тестовая для установки env а не системная | https://golangci-lint.run/usage/linters/#tenv           |
| testableexamples          | Следит за тем чтобы примеры были тестируемыми                                       |                                                         |
| unconvert                 | Удаление необязательных приведений типов                                            |                                                         |
| unparam                   | Удаление ненужных параметров в функциях                                             | https://golangci-lint.run/usage/linters/#unparam        |
| whitespace                | Проверка на лишние пустые строки в коде                                             | https://golangci-lint.run/usage/linters/#whitespace     |

### Спорные (но включены)

| Линт              | Описание                                                                                        | Ссылки                                                  |
|-------------------|-------------------------------------------------------------------------------------------------|---------------------------------------------------------|
| gocognit          | Проверка [когнитивной](https://linearb.io/blog/cognitive-complexity-in-software) сложности кода | https://golangci-lint.run/usage/linters/#gocognit       |
| goprintffuncname  | Функции типа printf должны оканчиваться на f                                                    |                                                         |
| importas          | Избежание aliases при импорте                                                                   | https://golangci-lint.run/usage/linters/#importas       |
| noctx             | Срабатывает когда посылается http request без контекста                                         |                                                         |
| nonamedreturns    | Отключат все именованные return                                                                 | https://golangci-lint.run/usage/linters/#nonamedreturns |
| nosprintfhostport | При формирование url не используется Sprintf                                                    |                                                         |
| promlinter        | Проверка наименования метрик Prometheus                                                         | https://golangci-lint.run/usage/linters/#promlinter     |
| protogetter       | Проверка на наличие геттеров у proto                                                            | https://golangci-lint.run/usage/linters/#protogetter    |
| reassign          | Проверка на переопределение стандартных значений другого пакета                                 | https://golangci-lint.run/usage/linters/#reassign       |
| sqlclosecheck     | sql.Rows и sql.Stmt должны быть закрыты после использования                                     |                                                         |
| testifylint       | Следит за использованием библиотеки testify                                                     | https://golangci-lint.run/usage/linters/#testifylint    |
| usestdlibvars     | Проверка на возможность использования констант из стандартной библиотеки                        | https://golangci-lint.run/usage/linters/#usestdlibvars  |
| wrapcheck         | Оборачивание ошибок                                                                             | https://golangci-lint.run/usage/linters/#wrapcheck      |
| zerologlint       | Проверка правильного использования zerolog                                                      |                                                         |

