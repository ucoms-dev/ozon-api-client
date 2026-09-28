# Ozon promotions SDK — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox syntax for tracking. Независимые read-only reviewers проверяют итоговый диапазон после task commit. Реализация разрешена пользователем и выполнена; см. verification report.

**Goal:** Добавить все восемь новых методов акций в github.com/ucoms-dev/ozon-api-client, сохранив действующие публичные контракты линейки v1.

**Architecture:** Новые методы на существующем Promotions и отдельные DTO для несовместимых wire-схем. Точное десятичное представление денег, корректное чтение ID, локальная валидация новых success-ответов. SDK не принимает бизнес-решения о цене/участии и не оркестрирует повторные записи.

**Tech Stack:** Go 1.20 minimum (текущий go.mod), net/http, encoding/json, httptest, существующий testify; без новых production dependencies. Проверено на установленном Go 1.26.4 linux/amd64.

**Spec:** ../specs/2026-09-28-ozon-promotions-sdk-design.md

## Исходная точка и полномочия

- Repo /home/xdev/go/ozon-api-client; origin git@github.com:ucoms-dev/ozon-api-client.git.
- Live dev и v1.17.0 на момент подготовки: efe8b033adb4df1060035c89ef77fa52d8740796.
- Worktree плана: /home/xdev/go/.codex-worktrees/ozon-sdk-actions-plan-20260928.
- Ветка: codex/ozon-sdk-actions-plan-20260928.
- Базовый go test ./... выполнен 2026-09-28: все 5 пакетов прошли.
- Реализованы Tasks 1–6 и подготовительная часть Task 7. Фактические команды и ограничения: [verification report](../../ozon-promotions-verification-2026-09-28.md). Публикация и миграция UCOMS не выполнены.
- До реализации заново проверить remote, инструкции и актуальность worktree. SDK сейчас не содержит AGENTS.md/CLAUDE.md и UCOMS agent-context scripts; не переносить receipts между репозиториями. При отдельной миграции UCOMS соблюдать его agent-context/git lifecycle.

## Global Constraints

- Go 1.20 minimum; module github.com/ucoms-dev/ozon-api-client.
- Старые публичные методы и DTO сохраняют сигнатуры, маршруты, JSON и обработку HTTP-ошибок.
- Endpoint обновления: /v1/actions/products/update. GET /v1/actions остаётся действующим.
- Размер страницы новых списков 1..100; размер mutation batch 1..1000.
- Обычные списки: last_id string. Auto-add списки: offset + limit + auto_add_date.
- Money = amount string + currency string; денежный float64 не вводить. Диагностическое numeric value = json.Number.
- Один вызов — один запрос; без повторов, фоновых задач, автоматического chunking, clock-switch и скрытых fallback.
- Незнакомый success JSON не означает пустой результат. Незнакомые дополнительные поля допустимы.
- SDK не вычисляет action_price и не подтверждает доставку изменённой цены/будущего расписания.
- Репозиторий и Go-команды исполняются в Ubuntu WSL через zsh -lc. Browser остаётся Windows.

## Review Focus

1. Числовые ID выше 2^53 и вплоть до uint64 max не теряют точность; тест Task 1.
2. Пустой 200 и чужой result не очищают смысл ответа, при этом существующие empty-200 methods работают; Task 2/3/6.
3. Отсутствующий stock отличается от явного нуля, исходные Params не изменяются; Task 3/5.
4. HTTP 200 с частичным результатом не становится общим success и не вызывает повторную запись; Task 3/5/6.
5. Изменение поведения 13 октября и даты auto_add_date не вызывает локальное изменение маршрута или денег; Task 5/6/7.

## Карта файлов

Создать:
- ozon/promotions_types.go — Money, ID, issue, website prices, ошибки и локальная wire-валидация.
- ozon/promotions_v2.go — два списка и два немедленных mutation methods, их DTO.
- ozon/promotions_auto_add_v2.go — четыре auto-add метода и DTO.
- ozon/promotions_types_test.go, ozon/promotions_v2_contract_test.go, ozon/promotions_auto_add_v2_contract_test.go — проверки сквозного поведения клиента.
- ozon/promotions_v2_examples_test.go — компилируемые примеры без обращений к Ozon.
- ozon/testdata/promotions-v2/ — обезличенные/синтетические JSON fixtures по утверждённой схеме.
- docs/ozon-promotions-contracts-2026-09-28.md — field matrix, источники, расхождения, provenance fixtures.
- docs/ozon-promotions-migration.md — руководство потребителям.

Изменить:
- ozon/promotions.go — AutoAddDates, Deprecated-комментарии для четырёх старых методов.
- ozon/promotions_legacy_contract_test.go — regression AutoAddDates и неизменность старых маршрутов.
- README.md — новые методы, ссылка на guide, условия совместимости.
- .github/workflows/tests.yml — добавить dev к существующим master triggers; проверять минимальный и поддерживаемый Go.
- internal/contractaudit/promotions_test.go — обнаружение восьми новых маршрутов существующим аудитом.
- scripts/vet-promotions.sh — vet затронутого кода и новых тестов с включёнными анализаторами.

client.go/core.go не требуют глобальной смены семантики; не делать unrelated refactor. Если реализация покажет необходимость их менять, сначала обновить spec и compatibility-тесты. go.mod/go.sum не меняются ради версии библиотеки; версия задаётся тегом.

## Task 1: Типы wire-контракта и доказательства схемы

**Files:** promotions_types.go, promotions_types_test.go, testdata/promotions-v2/, docs/ozon-promotions-contracts-2026-09-28.md.
**Consumes:** spec Types and Encoding, подтверждённые карточки восьми методов.
**Produces:** PromotionMoney; PromotionProductID с MarshalJSON/UnmarshalJSON; PromotionIssue; PromotionPriceValidation; PromotionWebsitePrices/PromotionSchemaPrices; PromotionInputError/PromotionProtocolError. Формы и имена полей зафиксированы в spec.

- [x] Зафиксировать field matrix отдельно для каждого endpoint. Пометить fixtures как synthetic-spec или redacted-live, указать URL/дату/hash исходника; примеры из статьи не выдавать за реальные ответы.
- [x] До реализации типов написать TestPromotionProductIDExactWire, TestPromotionMoneyExactWire, TestPromotionOptionalStockWire. Assert: ID 9007199254740993 и 18446744073709551615 одинаково читаются числом/строкой; 1.5, -1, 1e3, overflow и null отвергаются. Money 12345678901234567890.123456789 сохраняет текст/валюту; отсутствие Stock не даёт JSON stock, explicit pointer(0) даёт stock:0.
- [x] Запустить go test ./ozon -run 'TestPromotion(ProductID|Money|OptionalStock)' -count=1, зафиксировать failure отсутствующих новых типов.
- [x] Реализовать типы и narrow validation без float64-конверсий; не переименовывать PostingMoney. Для диагностики value использовать json.Number.
- [x] Повторить команду: PASS. Добавить fuzz seeds для ID/Money JSON в тех же тестах; ошибки не паникуют.

## Task 2: Чтение участников и кандидатов

**Files:** promotions_v2.go, promotions_v2_contract_test.go; fixtures products.json/candidates.json/empty-products.json.
**Consumes:** типы Task 1, существующие Promotions.client и core.Client.Request.
**Produces:** ProductsInPromotionV2 и ProductsAvailableForPromotionV2 с Params/Response из spec; PromotionParticipantV2, PromotionCandidateV2 и общий PromotionProductPrices.

- [x] Написать TestPromotionV2ListsWire через публичный NewClient + локальный httptest.Server. Проверить POST path, заголовки Client-Id/Api-Key, action_id, limit:100, last_id:"cursor-1" и отсутствие offset; ответ с Money/website_prices/карантином проходит без потери полей.
- [x] Добавить TestPromotionV2ListsInvalidShape: empty body, null, {}, result:{products:[]}, products:null, отрицательный total и неверный Money дают ошибку; products:[],total:0,last_id:"" корректен. Дополнительное неизвестное поле допустимо.
- [x] Добавить TestPromotionV2ListBoundsAndOwnership: nil Params, ID=0, limit=0/101 отвергаются до HTTP; limit=1/100 принят; один вызов не следует за LastID автоматически и не мутирует Params. Отсутствующие поля кандидата не синтезируются как Stock=0/AddMode=MANUAL.
- [x] Выполнить go test ./ozon -run '^TestPromotionV2List' -count=1: FAIL до реализации.
- [x] Реализовать методы с явным url := "..." и c.client.Request в каждом методе, чтобы существующий AST route audit видел вызовы. Проверка envelope остаётся локальной новым DTO; зафиксировать наличие products и total отдельно от нулевых значений.
- [x] Повторить команду: PASS.

## Task 3: Немедленные изменения и исключение

**Files:** promotions_v2.go, promotions_v2_contract_test.go; fixtures update-mixed.json/deactivate-partial.json.
**Consumes:** Task 1, общий transport contract Task 2.
**Produces:** UpdateProducts(ctx, *UpdatePromotionProductsParams) (*UpdatePromotionProductsResponse, error); RemoveProductV2(ctx, *RemoveProductFromPromotionV2Params) (*RemoveProductFromPromotionV2Response, error).

- [x] Написать TestPromotionV2MutationsWire: update использует /v1/actions/products/update и product_id, amount string/currency; deactivate использует /v2/actions/products/deactivate и массив строк product_ids. В ответе одновременно active=[1],deactivated=[2],rejected=[3],warnings=[2]; все категории сохраняются.
- [x] Добавить TestPromotionV2MutationBounds: 1/1000 accepted, 0/1001/duplicate/zero-ID rejected до сети; nil Stock отсутствует, explicit zero передаётся. Сумма и валюта не подменяются, параметры не изменяются.
- [x] Добавить TestPromotionV2PartialAndUnknownResult: requested [1,2],confirmed [1] не превращается в confirmed [1,2]; пустой или неизвестный envelope — protocol error; recognized empty arrays сохраняются как неопределённый результат для caller. HTTP 200 с rejected не Go transport error.
- [x] Выполнить go test ./ozon -run '^TestPromotionV2(Mutation|Partial)' -count=1: FAIL, затем реализовать и повторить до PASS.
- [x] Godoc объясняет изменение предельной цены, ограничение stock по типу акции и отсутствие гарантий фактического участия; метод не сравнивает price с лимитом и не вызывает второй запрос.

## Task 4: Расписание и чтение auto-add

**Files:** promotions.go, promotions_test.go, promotions_auto_add_v2.go, promotions_auto_add_v2_contract_test.go.
**Consumes:** Task 1, существующий GetAvailablePromotions.
**Produces:** AutoAddDates []time.Time; ListAutoAddProductsV2 и ListAutoAddCandidatesV2; AutoAddProductV2 и AutoAddCandidateV2.

- [x] Написать TestAvailablePromotionsAutoAddDates: поле auto_add_dates читается, отсутствие совместимо со старым fixture, GET /v1/actions остаётся. Дата 2026-10-13T03:00:00.123456789+03:00 сохраняет instant/offset/точность.
- [x] Написать TestPromotionAutoAddV2ListsWire: action_id, auto_add_date, offset:5,limit:100, без last_id; list.product_id и candidate.id не смешиваются; list.add_mode bool и candidate.is_manually_added bool сохраняются; product_id и sku различны в fixture.
- [x] Добавить TestPromotionAutoAddV2ListBounds: limit=0/101, zero time, invalid IDs не отправляются; offset=0/5 принят; unknown optional Money/flag сохраняется nil; правильный empty list отличим от malformed envelope.
- [x] Выполнить go test ./ozon -run 'TestAvailablePromotionsAutoAddDates|TestPromotionAutoAddV2List' -count=1: FAIL, реализовать и повторить до PASS.

## Task 5: Запланированное изменение и удаление

**Files:** promotions_auto_add_v2.go, promotions_auto_add_v2_contract_test.go; fixtures auto-add-update-mixed.json/auto-add-delete.json/auto-add-price-validation.json.
**Consumes:** Tasks 1 and 4.
**Produces:** UpdateAutoAddProductsV2 и DeleteAutoAddProductsV2, их Params/Response по spec; UpdateAutoAddProductV2.

- [x] Написать TestPromotionAutoAddV2MutationsWire: update содержит products:[{id:...,action_price:{...}}], не to_update и не product_id; delete содержит product_ids и ту же auto_add_date. Ответ update хранит product_ids,deactivated_ids,rejected,warnings и три price-validation массива отдельно.
- [x] Добавить TestPromotionAutoAddV2Diagnostics: схема [{key:123,value:12.34567890123456789}] сохраняется точно; ошибочные object вместо array и empty string вместо number не выдаются за успешную обработку. Если подтверждённый live/OpenAPI контракт иной, обновить spec, DTO и fixtures вместе до реализации этого ответа; не добавлять any/RawMessage в публичный результат ради замалчивания.
- [x] Добавить TestPromotionAutoAddV2OwnershipAndLimits: stock omitted/zero, 1/1000 accepted, 0/1001/duplicates invalid, zero date invalid, исходный запрос не меняется. Полученный HTTP ответ не создаёт таймер или второй HTTP вызов.
- [x] Выполнить go test ./ozon -run '^TestPromotionAutoAddV2(Mutation|Diagnostic|Ownership)' -count=1: FAIL, реализовать, повторить до PASS.
- [x] В guide объяснить различие принятого расписания и применённого изменения в дату AutoAddDate; указать позднюю проверку через list/participants на стороне приложения.

## Task 6: Транспорт, ошибки, совместимость и аудит

**Files:** новые contract tests, promotions.go/promotions_test.go, internal/contractaudit/audit_test.go; client.go не изменять без доказанной необходимости.
**Consumes:** все восемь методов.
**Produces:** подтверждённая совместимость и маршрутная полнота.

- [x] Добавить TestPromotionV2TransportContract для всех восьми методов: 400/403/404/409/429/500 с корректным error body сохраняют CommonResponse; cancelled/deadline errors доступны errors.Is; malformed JSON выдаёт Go error. Server считает число запросов — не больше одного, включая 429/timeout; заранее отменённый context может дать ноль сетевых запросов.
- [x] Добавить TestPromotionLegacyCompatibility: старые методы вызывают прежние URL и сериализуют прежние числовые деньги; GetAvailablePromotions/approve/decline не перенаправляются. Compile-only consumer test использует публичные типы старых вызовов.
- [x] Добавить Deprecated-комментарии к ProductsAvailableForPromotion, ProductsInPromotion, AddToPromotion, RemoveProduct с точной заменой и датой отключения; не менять сами вызовы.
- [x] Добавить TestPromotionV2RouteInventory: existing LoadClientOperations обнаруживает все восемь точных POST путей. Сравнение endpoint inventory не заменяет тестов JSON.
- [x] Запустить go test ./... и go vet ./...; выполнить go test -race ./ozon -run 'TestPromotion|TestAvailablePromotions' -count=1 для проверки неизменности/shared Params и context. Исправить только связанные проблемы; исходные baseline failures фиксировать отдельно.
- [x] При наличии актуального полного OpenAPI выполнить go run ./cmd/contract-audit -swagger <verified-local-file> -client ozon -date 2026-09-28 -output <new-audit-report>; старый snapshot не переписывать и не маркировать текущим. Без полного OpenAPI честно оставить field matrix из карточек и отдельный eight-route test; не генерировать ложный full audit.

## Task 7: Guide, CI, версия и передача потребителю

**Files:** README.md, docs/ozon-promotions-migration.md, promotions_v2_examples_test.go, .github/workflows/tests.yml.
**Consumes:** восемь методов и verified contracts Tasks 1..6.
**Produces:** проверенный релизный кандидат SDK; публикация требует отдельного разрешения.

- [x] Добавить компилируемые примеры чтения следующего LastID, offset auto-add, update с Money/stock:nil, проверки всех категорий ответа. Все примеры работают через local transport fixture и не вызывают реальный Ozon.
- [x] Guide: таблица старое→новое; новые response envelopes; ID/product_id/SKU; Money; отсутствие stock; date semantics; differences before/after 13 октября; limits; HTTP versus per-item errors; retry ambiguity; known spec contradictions.
- [x] Дополнить CI triggers dev, сохранив master; выставить Go matrix 1.20.x и 1.x (check-latest + stable filter) с go test ./... и bash scripts/vet-promotions.sh. Полный go vet имеет те же существующие ошибки, что база, на обеих версиях Go; см. verification report. Не добавлять внешние зависимости и не удалять существующий coverage job без отдельной причины.
- [x] Проверить минимальную версию командой GOTOOLCHAIN=go1.20.14 go test ./..., затем текущий toolchain go test ./... / go vet ./.... Недоступная загрузка toolchain — явно непрошедшая проверка, не PASS.
- [x] Проверить обновление зависимости на неизменённом UCOMS в отдельной свежей consumer-worktree через временный -modfile с replace на SDK candidate. Команда go test -c -modfile <temporary-go.mod> -o <temporary-binary> ./internal/service проверяет compilation без запуска TestMain и бизнес-операций; go.mod/go.sum основной копии не менять. Это доказательство source compatibility, не переход сервиса на новые методы.
- [ ] Сверить live SDK dev и теги; после проверок создать scoped task commit с Conventional Commit и Context/Changes/Verification/Task. Получить два independent read-only review exact base/head: contracts/correctness и tests/operations. P0/P1 устранить с повторной проверкой и review итогового диапазона; если SDK HEAD изменился — портировать diff на свежую базу и повторить проверки.
- [ ] Предложить v1.18.0 только если тег свободен. После разрешения публикации: обычная интеграция в dev, неизменяемый тег, release notes, проверка соответствия tag→reviewed commit/tree и загрузки модуля. Не переносить теги и не считать push доказательством доступности версии потребителю.
- [ ] Проверить отдельно согласованную миграцию UCOMS: заменить unsafe cursor shim на SDK методы, обновить UI/Mongo adapters/read allowlist, сериализацию цен и правила по типам акций, shared price-write fences, ревизии, readback и остановку активных задач. Это следующая задача, SDK release сам её не выполняет.

## Порядок и оценка

Task 1 → Tasks 2/3/4 → Task 5 → Task 6 → Task 7. Для одного исполнителя проходить последовательно: типы общие, файлы Tasks 2/3 и 4/5 пересекаются. Отдельные reviewers подключаются к готовому exact range.

Оценка для планирования, не SLA: 3–5 рабочих дней на восемь методов, fixtures, compatibility и review при доступных контрактах; подтверждение непустой auto-add диагностики может добавить ожидание. Миграция и выпуск UCOMS оцениваются отдельно.

## Definition of Done

- [x] Реализованы восемь endpoints и AutoAddDates, все поля field matrix имеют DTO/test.
- [x] Money/ID точны; отсутствующие значения не превращены в подтверждённые нули.
- [x] Partial success и protocol failures различаются; ошибки и context не потеряны.
- [x] Старые методы не изменили поведение; existing consumer compile проверен.
- [x] Go 1.20 и current tests, scoped vet и race прошли. Полный vet не проходит: диагностика совпадает с неизменённой базой; unrelated legacy-test debt не исправлялся.
- [x] Документированы все противоречия; ответ auto-add diagnostics подтверждён перед заявлением о полной совместимости.
- [ ] Два exact-range review без P0/P1, исходная база свежая; публикация не смешана с deployment UCOMS.

## Проверка самого плана

Проверено при подготовке: восемь строк endpoint matrix имеют отдельные методы и шаги; date/stock/ID/Money/partial-result риски покрыты задачами; имена методов и DTO согласованы со spec; Tasks 2/3 и 4/5 намеренно последовательны; Go минимум взят из свежего go.mod. Исходные тесты SDK прошли. При реализации добавлены и выполнены новые HTTP-контрактные тесты; их фактические имена могут отличаться от плановых. RED для отсутствующих типов/методов зафиксирован перед реализацией, общий набор на обеих версиях Go прошёл. Проверка optional stock находится в mutation tests. Полный OpenAPI просмотрен в браузере; локально сохранён только нормализованный extract, поэтому полный Swagger audit не заявляется. Независимые exact-range reviews выполняются после task commit и отражаются в итоговой передаче.
