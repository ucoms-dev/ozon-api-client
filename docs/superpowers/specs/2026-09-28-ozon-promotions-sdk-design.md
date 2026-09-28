# Ozon promotions SDK contracts — design

Дата: 2026-09-28. Статус: спецификация принята пользователем для реализации SDK.
Репозиторий: github.com/ucoms-dev/ozon-api-client.
Проверенная база dev и v1.17.0: efe8b033adb4df1060035c89ef77fa52d8740796.

## Цель и граница

Полностью поддержать восемь методов обновления акций из объявления Ozon от сентября 2026 года в существующем Go SDK. Сохранить исходную совместимость текущих потребителей и подготовить отдельный переход UCOMS. Смена бизнес-логики продавца, миграция Mongo, очередей и интерфейса UCOMS не входят в реализацию SDK.

Выбран подход: новые методы и DTO на существующем Promotions, старые методы сохраняют свои запросы и получают Deprecated-комментарии. Это соответствует docs/ozon-seller-api-contract-plan-2026-09-02.md. Замена старых методов на новые внутри прежних сигнатур отвергнута: Money, структура ответа и смысл записи несовместимы. Новый модуль /v2 отвергнут: удалять старую поверхность не требуется.

SDK не вычисляет цену исключения, не узнаёт тип акции дополнительным запросом, не запускает таймер 13 октября, не делает скрытый fallback на старый endpoint, не разбивает запись на партии и не повторяет запись после сетевой ошибки. Один вызов — одна HTTP-операция. Повторное чтение и политика повторов принадлежат потребителю.

## Первичные источники

Карточки и раскрытые поля прочитаны в Brave 2026-09-28:
- https://dev.ozon.ru/start/563-Obnovlenie-metodov-raboty-s-aktsiiami-v-Seller-API-v2/
- https://docs.ozon.ru/api/seller/#operation/Promos
- https://docs.ozon.ru/api/seller/#operation/ActionsCandidates
- https://docs.ozon.ru/api/seller/#operation/ActionsProducts
- https://docs.ozon.ru/api/seller/#operation/ActionsProductsUpdate
- https://docs.ozon.ru/api/seller/#operation/ActionsProductsDeactivate
- https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsListV2
- https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsCandidatesV2
- https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsUpdateV2
- https://docs.ozon.ru/api/seller/#operation/ActionsAutoAddProductsDeleteV2

До 2026-10-13 новые методы имеют прежнюю бизнес-логику. В эту дату отключаются соответствующие старые методы. Это не обещание прежней JSON-схемы. GET /v1/actions остаётся; к его DTO необходимо добавить AutoAddDates []time.Time с json:auto_add_dates,omitempty. SDK сохраняет offset/timezone и наносекунды при date-time round trip, не определяет доступность даты по своим часам.

Заявки покупателей discounts-task и акции продавца seller-actions — отдельные семейства. Их миграция не включена; /v1/actions/discounts-task/list имеет отдельное уведомление об устаревании без установленной здесь даты отключения. approve/decline не перенаправлять.

## Публичные методы

Все методы имеют receiver (c Promotions), принимают context.Context первым и *Params вторым, возвращают (*Response, error). Типы Params и Response в таблице однозначны.

| Метод | Params / Response | POST endpoint |
|---|---|---|
| ProductsAvailableForPromotionV2 | ProductsAvailableForPromotionV2Params / ProductsAvailableForPromotionV2Response | /v2/actions/candidates |
| ProductsInPromotionV2 | ProductsInPromotionV2Params / ProductsInPromotionV2Response | /v2/actions/products |
| UpdateProducts | UpdatePromotionProductsParams / UpdatePromotionProductsResponse | /v1/actions/products/update |
| RemoveProductV2 | RemoveProductFromPromotionV2Params / RemoveProductFromPromotionV2Response | /v2/actions/products/deactivate |
| ListAutoAddProductsV2 | ListAutoAddProductsV2Params / ListAutoAddProductsV2Response | /v2/actions/auto-add/products/list |
| ListAutoAddCandidatesV2 | ListAutoAddCandidatesV2Params / ListAutoAddCandidatesV2Response | /v2/actions/auto-add/products/candidates |
| UpdateAutoAddProductsV2 | UpdateAutoAddProductsV2Params / UpdateAutoAddProductsV2Response | /v2/actions/auto-add/products/update |
| DeleteAutoAddProductsV2 | DeleteAutoAddProductsV2Params / DeleteAutoAddProductsV2Response | /v2/actions/auto-add/products/delete |

UpdateProducts не имеет суффикса V2: фактический endpoint содержит v1. Все ответы включают core.CommonResponse и собственные верхнеуровневые поля, без Result. Старые типы PromotionProduct и AddProductToPromotionParams не изменяются.

## Типы и кодирование

- PromotionMoney: Amount string, Currency string, теги amount/currency. Это транспортный DTO, а не доменная сущность. Строка сохраняет десятичную точность, никаких ParseFloat или округлений. Проверять синтаксис десятичной суммы и валюту при формировании записи; не навязывать RUB или масштаб 2 всему SDK. Не разрешать отрицательную цену, NaN, Inf, экспоненту и пустую сумму. Ноль — допустимое числовое значение на транспортной границе, решение Ozon о допустимости цены сохраняется в ответе. Отсутствующее Money в ответе — nil, не нулевая цена.
- PromotionProductID uint64: декодирует JSON integer и decimal string без float64; допускает обрамляющие пробельные символы строки, запрещает дробь, знак, экспоненту, переполнение, null и пустую строку. MarshalJSON выдаёт integer. Требование положительности проверяется отдельно для товарных идентификаторов запросов/ответов.
- ActionID, Limit, Offset и количества: uint64; необязательные количества/флаги в ответах — указатели, когда отсутствие имеет смысл. Не переводить max uint64 в int64 внутри SDK.
- Запросы update: product_id/id кодируются числом. Запросы deactivate/delete: публичные ProductIDs []PromotionProductID преобразуются в decimal strings согласно Array of strings<uint64> в спецификации. Ответы принимают оба представления ID, поскольку примеры Ozon отличаются от схемы.
- PromotionIssue: ProductID PromotionProductID, Reason string, теги product_id/reason.
- PromotionPriceValidation: Key PromotionProductID, Value json.Number, теги key/value. Диагностическая цена в auto-add документирована number<double>, это не Money. Сохранять число без float64.
- PromotionWebsitePrices: Price *PromotionMoney, PricesBySchema map[string]PromotionSchemaPrices. PromotionSchemaPrices: BlackPrice *PromotionMoney, GreenPrice *PromotionMoney. Ключи схем доставки сохраняются без жёсткого enum.
- PromotionInputError: Method, Field, Reason; PromotionProtocolError: Method, Reason. errors.As поддерживается; не помещать в ошибки ключи или полный HTTP body. Причина ошибки JSON/контекста должна оставаться доступной через Unwrap там, где она существует.

Эти типы находятся в ozon/promotions_types.go. Существующий PostingMoney относится к контракту отправлений: его не менять и не использовать в публичном API акций для экономии двух строк.

## Запросы и ответы

1. products/candidates: ActionID, Limit (1..100), LastID string, первый курсор пустой. Никакого offset и автоматического полного обхода. Ответ: Products, Total uint64, LastID string. Отдельные PromotionParticipantV2 и PromotionCandidateV2; общие ценовые поля можно встроить из PromotionProductPrices, без искусственных полей Stock у кандидата.
2. Общие поля продукта: ID, Price, ActionPrice, MaxActionPrice, AlertMaxActionPrice, AlertMaxActionPriceFailed, CurrentBoost, PriceMinElastic, PriceMaxElastic, MinBoost, MaxBoost, MinStock, RecommendedStock, MarketplaceSellerPrice, MinSellerPrice, IsQuarantined, WebsitePrices. Денежные поля — *PromotionMoney, boost — float64 (не денежный учёт), количества uint64/указатели, флаги *bool. У участника дополнительно AddMode string (SELLER/AUTO, неизвестное значение сохраняется), Stock *uint64. У кандидата Stock/AddMode в текущей схеме отсутствуют.
3. UpdateProducts: ActionID, Products []UpdatePromotionProduct; элемент ProductID, ActionPrice PromotionMoney, Stock *uint64 с omitempty. nil отличается от явного 0. 1..1000 элементов, положительные уникальные ID. Ответ ActiveProductIDs, DeactivatedProductIDs, Rejected []PromotionIssue, Warnings []PromotionIssue.
4. RemoveProductV2: ActionID, ProductIDs (1..1000, уникальные положительные). Ответ ProductIDs. Не добавлять выдуманный Rejected, отсутствующий в новой схеме. Неполный набор подтверждённых ID возвращается вызывающему коду, не расширяется до исходного запроса.
5. ListAutoAddProductsV2 / ListAutoAddCandidatesV2: ActionID, AutoAddDate time.Time (не zero), Limit 1..100, Offset. Ответ Products и Total. Это offset-пагинация, не last_id. Дата берётся из AutoAddDates, SDK не выдумывает расписание.
6. AutoAddProductV2 и AutoAddCandidateV2 имеют Money-поля Price, BasePrice, MaxDiscountPrice, ActionPriceToAutoAdd, MarketplaceSellerPrice, MinSellerPrice; также Currency, OfferID, SKU, Name, MinActionQuantity, QuantityToAutoAdd, HasExpiredMinSellerPrice, WillBeQuarantined, WebsitePrices. Product использует ProductID и AddMode *bool; candidate — ID и IsManuallyAdded *bool. Нельзя переиспользовать строковый add_mode обычного участника.
7. UpdateAutoAddProductsV2: ActionID, AutoAddDate, Products []UpdateAutoAddProductV2; элемент ID с json:id, ActionPrice PromotionMoney, Stock *uint64. Ответ ProductIDs, DeactivatedIDs (json:deactivated_ids, без product в имени), Rejected, Warnings, BelowMinPrice, ExtremelyLowPrice, FailedPrice. Последние три — []PromotionPriceValidation.
8. DeleteAutoAddProductsV2: ActionID, AutoAddDate, ProductIDs 1..1000. Ответ ProductIDs.

## Защита структуры ответа и сохранение транспорта

Для новых methods пустое тело HTTP 200, null, {}, вложенный result вместо верхнеуровневого ответа и неверные типы обязательных полей — PromotionProtocolError. Списки требуют распознанного массива products и total, включая корректные products:[] и total:0. null products не считать пустым массивом. Для мутаций должен присутствовать хотя бы один документированный результат, правильного типа; пустые распознанные массивы допустимы и не означают подтверждение всех запрошенных товаров. Для deactivate/delete требуется product_ids.

Не вводить DisallowUnknownFields на весь ответ: дополнительные поля Ozon не должны ломать чтение. Не суммировать rejected/warnings с подтверждениями как успешные операции. Ошибка декодирования не превращается в HTTP 200 с пустыми данными. Реализовать проверку в новых ответах и методах; глобальное разрешение пустого 200 в client.go необходимо другим существующим операциям и сохраняется.

Сохранять существующий контракт SDK: корректный ответ Ozon 4xx/5xx возвращает Response с CommonResponse.StatusCode/Code/Message/Details; Go error означает транспорт, контекст, некорректный JSON, входной контракт или неверную success-схему. HTTP 200 с rejected — нормальный частичный ответ. Проверку success-схемы запускать только для 200. Чувствительные заголовки никогда не включать в diagnostics.

Параметры вызывающего кода не мутировать: до применения default/нормализации использовать локальную копию; новый код не добавляет default-теги, изменяющие исходный объект. Nil Params и недопустимые значения отсекаются до Request. Сохраняется один запрос на вызов, cancellation/deadline и действующий HttpClient.

## Поведение до и после 13 октября

SDK использует фиксированный endpoint выбранного метода в обе даты. Новые методы не меняются автоматически по локальным часам и не повышают цену для исключения. Godoc и migration guide объясняют: после переключения Ozon UpdateProducts вне промокодов меняет предельную цену; auto-add применяет её в AutoAddDate; stock для бустинга не передаётся; HTTP-ответ на расписание не доказывает применение в будущем.

Уведомление про stock/тип акции документировать, но не делать дополнительный запрос списка акций для скрытой бизнес-валидации: ActionType отсутствует в запросе SDK. Проверять форму запроса, а выбор Stock оставить вызывающему коду.

## Неопределённости первоисточника и решение

- Вводный текст deactivate противоречит подробному полю update о направлении сравнения цены. SDK не вычисляет это сравнение; guide отмечает противоречие. Для включения ценовой автоматики UCOMS требуется подтверждение.
- Скидка на сток и Максимальный бустинг не объявляются SDK синонимами.
- Auto-add diagnostics: таблица схемы задаёт массив {key:uint64,value:number}, пример содержит одиночный объект и пустую строку. 2026-09-28 при реализации получен актуальный OpenAPI 3.0.0, загруженный страницей документации: https://docs.ozon.ru/api/seller/swagger.json?1790601089484. Его components/schemas/actions.v2.ActionsAutoAddProductsUpdateResponse подтверждает массивы BelowMinPrice/ExtremelyLowPrice/FailedPrice с key:uint64 и value:number<double>. Реализуется схема, неверная форма не проглатывается. Нормализованная матрица восьми методов сохранена в ozon/testdata/promotions-v2/swagger-shapes.json; это извлечение полей, не полная копия Swagger и не реальный ответ авторизованного Seller API.
- Никаких invented fallback v1/v2 based on errors или форматов из устаревшего Swagger 2026-09-02.

## Критерий завершения

Восемь методов доступны через Client.Promotions(); DTO покрывают раскрытые поля; старые потребители компилируются; контрактные тесты охватывают запрос, ответ, отсутствие полей, точность, частичные результаты и cancellation. Обновлены guide, route audit и CI для dev. Есть два независимых review на точном диапазоне SDK; P0/P1 устранены. Предлагаемая версия v1.18.0 при свободном теге. Публикация SDK и последующая миграция UCOMS — отдельные проверяемые события; подготовленный план не является разрешением на публикацию или продакшен.
