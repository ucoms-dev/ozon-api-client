# Ozon Seller API: проверка rate limits и SDK

Проверка официального Swagger через существующую вкладку Brave. Время снимка: 2026-10-06T21:47:40.636Z (UTC). Исходник загрузчика: `docs.ozon.ru/api/seller/scripts/loadSwagger.js`, загрузка `./swagger.json?Date.now()`. Полный документ разобран в браузере; сохранён нормализованный перечень всех 482 операций. Rate limits искались по всем строковым полям, включая схемы и общие разделы. После исключения описаний видео, частичных отправлений и поисковой аналитики с ложными совпадениями выделены 23 метода с особыми ограничениями.

Общий бюджет: **50 запросов/сек суммарно на Client ID**, а не по 50 на каждый метод или API-ключ. Он действует вместе с более строгими ограничениями метода/товара. `Ratelimit-Remaining` показывает остаток; `Retry-After` задаёт секунды. Товарные квоты: `Item-Rate-Limit-Remaining` и `Item-Retry-After` (**минуты**); суточный сброс 03:00 Москва. `Circle is open` и повторные одинаковые/ошибочные запросы могут временно ограничить доступ; документация не даёт отдельного численного бюджета для такого ограничения.

## Все опубликованные особые ограничения

| Метод | Область | Ограничение | Примечание |
|---|---|---|---|
| `/v1/product/placement-zone/info` | method | 10 requests/second |  |
| `/v1/barcode/add` | account/method | 20 calls/minute | up to 100 barcodes per product |
| `/v1/barcode/generate` | account/method | 20 calls/minute | up to 100 products/request |
| `/v2/products/stocks` | account/method + offer/warehouse | 80 requests/minute; same pair once/30 seconds | 100 offer/warehouse pairs/request |
| `/v1/product/import/prices` | offer | 10 price updates/hour per product | 1000 prices/request |
| `/v1/warehouse/list` | method | 1 call/minute | description announces sunset 2026-04-07; active v2 is a separate endpoint |
| `/v1/draft/crossdock/create` | method | 2/minute + 50/hour + 500/day |  |
| `/v1/draft/direct/create` | method | 2/minute + 50/hour + 500/day |  |
| `/v1/draft/multi-cluster/create` | method | 2/minute + 50/hour + 500/day |  |
| `/v2/draft/create/info` | method; description needs clarification | 2/minute + 50/hour | read endpoint description says creating drafts; preserve ambiguity; no inferred daily cap |
| `/v2/draft/supply/create` | method | 2/minute + 50/hour + 500/day |  |
| `/v1/report/discounted/create` | account/method | 1 request/minute |  |
| `/v1/report/placement/by-products/create` | report/method | 5/day |  |
| `/v1/report/placement/by-supplies/create` | report/method | 5/day |  |
| `/v1/analytics/turnover/stocks` | Client ID/method | 1 request/minute |  |
| `/v1/analytics/data` | account/subscription | 50/day without Premium Plus | does not infer an unlimited budget with Premium Plus |
| `/v3/product/import` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |
| `/v1/product/import-by-sku` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |
| `/v1/product/attributes/update` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |
| `/v1/product/pictures/import` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow; v1 pictures description announces sunset 2026-10-01 |
| `/v2/product/pictures/import` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |
| `/v1/product/update/offer-id` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |
| `/v1/product/unarchive` | product-operation quota | minute/day quota from provider; Item-Retry-After uses minutes | daily reset 03:00 Moscow |

Разархивация автоматически архивированных товаров: до 100 товаров за сутки, до 100 IDs в запросе, сброс 03:00 Москва. Для товаров, архивированных вручную, это ограничение восстановления не указано. `v3/product/import` — до 100 товаров в запросе; `v2/product/pictures/import` — до 100 элементов. Размер пакета и число HTTP-запросов учитываются отдельно от количества товарных операций. У `/v2/warehouse/list` нет опубликованного собственного лимита 1/мин; его нельзя переносить с устаревшего v1 без доказательства.


## Использование в SDK

Ответы содержат исходные HTTP-заголовки в `CommonResponse.Headers`. Значения доступны через `Headers.Get`, отсутствие заголовка возвращает пустую строку. SDK сохраняет заголовки и на успешных ответах, и на JSON-ошибках, включая HTTP 429. `CopyCommonResponse` копирует карту заголовков. Заголовки исключены из JSON DTO.

SDK остаётся типизированным транспортом: без автоматического ожидания или повторов. Ограничение общего бюджета Client ID, конкуренция нескольких процессов, товарные окна и безопасные повторы принадлежат вызывающему приложению. HTTP 429 по прежнему возвращается как ответ с `StatusCode`, `Code` и `Message`, а не как Go error. Проверять `err` без проверки `StatusCode` недостаточно.

`GetProductRangeLimit` теперь сохраняет `operation_limits`; чтение принимает документированный объект и наблюдавшийся массив. Неизвестные `limit_type` сохраняются. Отсутствующий или null `operation_limits` остаётся nil. `quota_by_category` не описан в текущем Swagger: его контракт не выдумывается.
