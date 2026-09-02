# Ozon Seller API contract audit

- Generated: 2026-09-02
- Swagger: `swagger (1).json`
- Swagger SHA-256: `48a0c82a411bf08c9769d61be6d444488af3062cc7e97d2d50067e16d4b05bb8`
- Client operations: 224
- Swagger operations: 465
- Exact method/path matches: 199
- Client-only paths: 25
- Method mismatches: 0
- Swagger-only operations: 266
- Deprecated exact matches: 11
- Deprecated matches missing Go doc: 0

## Client paths absent from Swagger

| Method | Path | Go method | Source |
|---|---|---|---|
| POST | `/v1/chat/updates` | `Update` | `chats.go` |
| POST | `/v1/conditional-cancellation/approve` | `Approve` | `cancellations.go` |
| POST | `/v1/conditional-cancellation/get` | `GetInfo` | `cancellations.go` |
| POST | `/v1/conditional-cancellation/list` | `List` | `cancellations.go` |
| POST | `/v1/conditional-cancellation/reject` | `Reject` | `cancellations.go` |
| GET | `/v1/draft/create` | `CreateSupplyDraft` | `fbo.go` |
| GET | `/v1/draft/create/info` | `GetSupplyDraftInfo` | `fbo.go` |
| GET | `/v1/draft/supply/create` | `CreateSupplyFromDraft` | `fbo.go` |
| GET | `/v1/draft/timeslot/info` | `GetDraftTimeslots` | `fbo.go` |
| POST | `/v1/product/import/stocks` | `UpdateStocks` | `products.go` |
| POST | `/v1/product/upload_digital_codes` | `UploadActivationCodes` | `products.go` |
| POST | `/v1/product/upload_digital_codes/info` | `StatusOfUploadingActivationCodes` | `products.go` |
| POST | `/v1/quant/get` | `Get` | `quants.go` |
| POST | `/v1/quant/list` | `List` | `quants.go` |
| POST | `/v1/quant/ship` | `Ship, Status` | `quants.go` |
| POST | `/v2/chat/list` | `List` | `chats.go` |
| POST | `/v2/fbs/posting/sent-by-seller` | `ChangeStatusToSendBySeller` | `fbs.go` |
| POST | `/v2/posting/fbs/product/change` | `AddWeightForBulkProduct` | `fbs.go` |
| POST | `/v2/returns/rfbs/compensate` | `CompensateRFBSReturn` | `returns.go` |
| POST | `/v2/returns/rfbs/receive-return` | `ReceiveRFBSReturn` | `returns.go` |
| POST | `/v2/returns/rfbs/reject` | `RejectRFBSReturn` | `returns.go` |
| POST | `/v2/returns/rfbs/return-money` | `RefundRFBS` | `returns.go` |
| POST | `/v2/returns/rfbs/verify` | `ApproveRFBSReturn` | `returns.go` |
| POST | `/v2/supply-order/get` | `GetSupplyRequestInfo` | `fbo.go` |
| POST | `/v2/supply-order/list` | `ListSupplyRequests` | `fbo.go` |

## Method mismatches

| Path | Client | Swagger | Go method | Source |
|---|---|---|---|---|

## Deprecated exact matches

| Method | Path | Go method | Go deprecated | Operation ID | Source |
|---|---|---|---|---|---|
| POST | `/v1/product/certificate/create` | `AddForProducts` | yes | `ProductAPI_ProductCertificateCreate` | `certificates.go` |
| POST | `/v1/review/change-status` | `ChangeStatus` | yes | `ReviewAPI_ReviewChangeStatus` | `reviews.go` |
| POST | `/v1/review/comment/delete` | `DeleteComment` | yes | `ReviewAPI_CommentDelete` | `reviews.go` |
| POST | `/v1/review/count` | `Count` | yes | `ReviewAPI_ReviewCount` | `reviews.go` |
| POST | `/v1/review/info` | `Get` | yes | `ReviewAPI_ReviewInfo` | `reviews.go` |
| POST | `/v1/review/list` | `List` | yes | `ReviewAPI_ReviewList` | `reviews.go` |
| POST | `/v2/posting/fbo/list` | `GetShipmentsList` | yes | `PostingAPI_GetFboPostingList` | `fbo.go` |
| POST | `/v2/posting/fbs/digital/act/check-status` | `GenerateAct` | yes | `PostingAPI_PostingFBSDigitalActCheckStatus` | `fbs.go` |
| POST | `/v2/posting/fbs/digital/act/get-pdf` | `GetDigitalAct` | yes | `PostingAPI_PostingFBSGetDigitalAct` | `fbs.go` |
| POST | `/v3/posting/fbs/list` | `GetFBSShipmentsList` | yes | `PostingAPI_GetFbsPostingListV3` | `fbs.go` |
| POST | `/v3/posting/fbs/unfulfilled/list` | `ListUnprocessedShipments` | yes | `PostingAPI_GetFbsPostingUnfulfilledList` | `fbs.go` |

## Swagger operations not implemented by the client

| Method | Path | Operation ID |
|---|---|---|
| POST | `/v1/actions/auto-add/products/candidates` | `ActionsAutoAddProductsCandidates` |
| POST | `/v1/actions/auto-add/products/delete` | `ActionsAutoAddProductsDelete` |
| POST | `/v1/actions/auto-add/products/list` | `ActionsAutoAddProductsList` |
| POST | `/v1/actions/auto-add/products/update` | `ActionsAutoAddProductsUpdate` |
| POST | `/v1/analytics/product-queries/details` | `AnalyticsAPI_AnalyticsProductQueriesDetails` |
| POST | `/v1/analytics/stocks` | `AnalyticsAPI_AnalyticsStocks` |
| POST | `/v1/assembly/carriage/posting/list` | `AssemblyCarriagePostingList` |
| POST | `/v1/assembly/carriage/product/list` | `AssemblyCarriageProductList` |
| POST | `/v1/assembly/fbs/posting/list` | `AssemblyFbsPostingList` |
| POST | `/v1/assembly/fbs/product/list` | `AssemblyFbsProductList` |
| POST | `/v1/cancel-reason/list` | `CancelReasonList` |
| POST | `/v1/cancel-reason/list-by-order` | `CancelReasonListByOrder` |
| POST | `/v1/cancel-reason/list-by-posting` | `CancelReasonAPI_CancelReasonListByPosting` |
| POST | `/v1/cargoes-label/create` | `CargoesAPI_CargoesLabelCreate` |
| GET | `/v1/cargoes-label/file/{file_guid}` | `CargoesAPI_CargoesLabelFile` |
| POST | `/v1/cargoes-label/get` | `CargoesAPI_CargoesLabelGet` |
| POST | `/v1/cargoes/create` | `CargoesAPI_CargoesCreate` |
| POST | `/v1/cargoes/delete` | `CargoesAPI_CargoesDelete` |
| POST | `/v1/cargoes/delete/status` | `CargoesAPI_CargoesDeleteStatus` |
| POST | `/v1/cargoes/get` | `CargoesGet` |
| POST | `/v1/cargoes/label/transport-by-order/create` | `CargoesLabelTransportByOrderCreate` |
| POST | `/v1/cargoes/label/transport-by-order/status` | `CargoesLabelTransportByOrderStatus` |
| POST | `/v1/cargoes/label/transport/create` | `CargoesLabelTransportCreate` |
| POST | `/v1/cargoes/label/transport/status` | `CargoesLabelTransportStatus` |
| POST | `/v1/cargoes/rules/get` | `CargoesAPI_CargoesRulesGet` |
| POST | `/v1/cargoes/supplies/get` | `CargoesSuppliesGet` |
| POST | `/v1/cargoes/transport/activate` | `CargoesTransportActivate` |
| POST | `/v1/cargoes/transport/activate/status` | `CargoesTransportActivateStatus` |
| POST | `/v1/cargoes/transport/bind` | `CargoesTransportBind` |
| POST | `/v1/cargoes/transport/bind/status` | `CargoesTransportBindStatus` |
| POST | `/v1/cargoes/transport/create` | `CargoesTransportCreate` |
| POST | `/v1/cargoes/transport/create/status` | `CargoesTransportCreateStatus` |
| POST | `/v1/carriage/act-discrepancy/pdf` | `CarriageActDiscrepancyPDF` |
| POST | `/v1/carriage/approve` | `CarriageAPI_CarriageApprove` |
| POST | `/v1/carriage/container/approve` | `CarriageContainerApprove` |
| POST | `/v1/carriage/container/cancel` | `CarriageContainerCancel` |
| POST | `/v1/carriage/container/create` | `CarriageContainerCreate` |
| POST | `/v1/carriage/container/document/get` | `CarriageContainerDocumentGet` |
| POST | `/v1/carriage/container/fill` | `CarriageContainerFill` |
| POST | `/v1/carriage/container/get` | `CarriageContainerGet` |
| POST | `/v1/carriage/container/label/get` | `CarriageContainerLabelGet` |
| POST | `/v1/carriage/container/list` | `CarriageContainerList` |
| POST | `/v1/carriage/container/place-into` | `CarriageContainerPlaceInto` |
| POST | `/v1/carriage/container/remove-from` | `CarriageContainerRemoveFrom` |
| POST | `/v1/carriage/container/remove-postings` | `CarriageContainerRemovePostings` |
| POST | `/v1/carriage/container/status/get` | `CarriageContainerStatusGet` |
| POST | `/v1/carriage/container/task/info` | `CarriageContainerTaskInfo` |
| POST | `/v1/carriage/courier-contact/get` | `CarriageCourierContactGet` |
| POST | `/v1/carriage/courier-contact/set` | `CarriageCourierContactSet` |
| POST | `/v1/carriage/create` | `CarriageAPI_CarriageCreate` |
| POST | `/v1/carriage/delivery/list` | `CarriageAPI_CarriageDeliveryList` |
| POST | `/v1/carriage/ettn/status` | `CarriageEttnStatus` |
| POST | `/v1/delivery-method/return/settings/get` | `GetDeliveryMethodReturnSettingsV1` |
| POST | `/v1/delivery/check` | `DeliveryCheck` |
| POST | `/v1/delivery/map` | `DeliveryMap` |
| POST | `/v1/delivery/point/info` | `DeliveryPointInfo` |
| POST | `/v1/delivery/point/list` | `DeliveryAPI_DeliveryPointList` |
| POST | `/v1/description-category/dependent-attributes` | `DescriptionCategoryDependentAttributes` |
| POST | `/v1/description-category/dependent-attributes/values` | `DescriptionCategoryDependentAttributesValues` |
| POST | `/v1/draft/crossdock/create` | `DraftCrossdockCreate` |
| POST | `/v1/draft/direct/create` | `DraftDirectCreate` |
| POST | `/v1/draft/multi-cluster/create` | `DraftMultiClusterCreate` |
| POST | `/v1/fbp/act-from/create` | `FbpAPI_FbpCreateAct` |
| POST | `/v1/fbp/act-from/get` | `FbpAPI_FbpCheckActState` |
| POST | `/v1/fbp/act-to/create` | `FbpAPI_FbpCreateConsignmentNote` |
| POST | `/v1/fbp/act-to/get` | `FbpAPI_FbpCheckConsignmentNoteState` |
| POST | `/v1/fbp/archive/get` | `FbpAPI_FbpArchiveGet` |
| POST | `/v1/fbp/archive/list` | `FbpAPI_FbpArchiveList` |
| POST | `/v1/fbp/draft/direct/create` | `FbpDraftDirectCreate` |
| POST | `/v1/fbp/draft/direct/delete` | `FbpDraftDirectDelete` |
| POST | `/v1/fbp/draft/direct/product/validate` | `FbpDraftDirectProductValidate` |
| POST | `/v1/fbp/draft/direct/registrate` | `FbpDraftDirectRegistrate` |
| POST | `/v1/fbp/draft/direct/seller-dlv/create` | `FbpDraftDirectSellerDlvCreate` |
| POST | `/v1/fbp/draft/direct/seller-dlv/edit` | `FbpDraftDirectSellerDlvEdit` |
| POST | `/v1/fbp/draft/direct/timeslot/edit` | `FbpDraftDirectTimeslotEdit` |
| POST | `/v1/fbp/draft/direct/timeslot/get` | `FbpDraftDirectGetTimeslot` |
| POST | `/v1/fbp/draft/direct/tpl-dlv/create` | `FbpAPI_FbpDraftDirectTplDlvCreate` |
| POST | `/v1/fbp/draft/direct/tpl-dlv/edit` | `FbpAPI_FbpDraftDirectTplDlvEdit` |
| POST | `/v1/fbp/draft/drop-off/create` | `FbpDraftDropOffCreate` |
| POST | `/v1/fbp/draft/drop-off/delete` | `FbpDraftDropOffDelete` |
| POST | `/v1/fbp/draft/drop-off/dlv/edit` | `FbpDraftDropOffDlvEdit` |
| POST | `/v1/fbp/draft/drop-off/point/list` | `FbpDraftDropOffPointList` |
| POST | `/v1/fbp/draft/drop-off/point/timetable` | `FbpDraftDropOffPointTimetable` |
| POST | `/v1/fbp/draft/drop-off/product/validate` | `FbpDraftDropOffProductValidate` |
| POST | `/v1/fbp/draft/drop-off/province/list` | `FbpDraftDropOffProvinceList` |
| POST | `/v1/fbp/draft/drop-off/registrate` | `FbpDraftDropOffRegistrate` |
| POST | `/v1/fbp/draft/get` | `FbpAPI_FbpDraftGet` |
| POST | `/v1/fbp/draft/list` | `FbpAPI_FbpDraftList` |
| POST | `/v1/fbp/draft/pick-up/create` | `FbpAPI_FbpDraftPickupCreate` |
| POST | `/v1/fbp/draft/pick-up/delete` | `FbpAPI_FbpDraftPickUpDelete` |
| POST | `/v1/fbp/draft/pick-up/dlv/edit` | `FbpAPI_FbpDraftPickupDlvEdit` |
| POST | `/v1/fbp/draft/pick-up/product/validate` | `FbpAPI_FbpDraftPickUpProductValidate` |
| POST | `/v1/fbp/draft/pick-up/registrate` | `FbpDraftPickUpRegistrate` |
| POST | `/v1/fbp/label/create` | `FbpAPI_FbpCreateLabel` |
| POST | `/v1/fbp/label/get` | `FbpAPI_FbpGetLabel` |
| POST | `/v1/fbp/order/direct/cancel` | `FbpAPI_FbpOrderDirectCancel` |
| POST | `/v1/fbp/order/direct/seller-dlv/edit` | `FbpAPI_FbpOrderDirectSellerDlvEdit` |
| POST | `/v1/fbp/order/direct/timeslot/edit` | `FbpAPI_FbpEditTimeslot` |
| POST | `/v1/fbp/order/direct/timeslot/list` | `FbpAPI_FbpAvailableTimeslotList` |
| POST | `/v1/fbp/order/drop-off/cancel` | `FbpAPI_FbpOrderDropOffCancel` |
| POST | `/v1/fbp/order/drop-off/dlv/edit` | `FbpAPI_FbpOrderDropOffDlvEdit` |
| POST | `/v1/fbp/order/drop-off/timetable` | `FbpAPI_FbpOrderDropOffTimetable` |
| POST | `/v1/fbp/order/get` | `FbpAPI_FbpOrderGet` |
| POST | `/v1/fbp/order/list` | `FbpAPI_FbpOrderList` |
| POST | `/v1/fbp/order/pick-up/cancel` | `FbpAPI_FbpOrderPickUpCancel` |
| POST | `/v1/fbp/order/pick-up/dlv/edit` | `FbpAPI_FbpOrderPickUpDlvEdit` |
| POST | `/v1/fbp/warehouse/list` | `FbpWarehouseList` |
| POST | `/v1/finance/accrual/by-day` | `GetFinanceAccrualByDay` |
| POST | `/v1/finance/accrual/postings` | `GetFinanceAccrualPostings` |
| POST | `/v1/finance/accrual/types` | `GetFinanceAccrualTypes` |
| POST | `/v1/finance/balance` | `GetFinanceBalanceV1` |
| POST | `/v1/finance/compensation` | `ReportAPI_GetCompensationReport` |
| POST | `/v1/finance/decompensation` | `ReportAPI_GetDecompensationReport` |
| POST | `/v1/finance/document-b2b-sales` | `ReportAPI_CreateDocumentB2BSalesReport` |
| POST | `/v1/finance/document-b2b-sales/json` | `ReportAPI_CreateDocumentB2BSalesJSONReport` |
| POST | `/v1/finance/products/buyout` | `GetFinanceProductsBuyout` |
| POST | `/v1/finance/realization/by-day` | `FinanceAPI_GetRealizationByDayReportV1` |
| POST | `/v1/finance/realization/posting` | `FinanceAPI_GetRealizationReportV1` |
| POST | `/v1/notification/check` | `CheckNotification` |
| POST | `/v1/notification/delete` | `DeleteNotification` |
| POST | `/v1/notification/enable` | `EnableNotification` |
| POST | `/v1/notification/list` | `NotificationList` |
| POST | `/v1/notification/push-type/list` | `GetNotificationPushTypeList` |
| POST | `/v1/notification/set` | `SetNotification` |
| POST | `/v1/notification/update` | `UpdateNotification` |
| POST | `/v1/order/cancel` | `OrderAPI_OrderCancel` |
| POST | `/v1/order/cancel/check` | `OrderAPI_OrderCancelCheck` |
| POST | `/v1/order/cancel/status` | `OrderAPI_OrderCancelStatus` |
| POST | `/v1/polygon/delete` | `PolygonDelete` |
| POST | `/v1/polygon/list` | `PolygonList` |
| POST | `/v1/polygon/time/coordinates/update` | `PolygonTimeCoordinatesUpdate` |
| POST | `/v1/polygon/time/set` | `PolygonTimeSet` |
| POST | `/v1/posting/cancel` | `PostingAPI_PostingCancel` |
| POST | `/v1/posting/cancel/status` | `PostingAPI_PostingCancelStatus` |
| POST | `/v1/posting/digital/codes/upload` | `UploadPostingCodes` |
| POST | `/v1/posting/digital/list` | `ListPostingCodes` |
| POST | `/v1/posting/fbp/get` | `GetFbpPosting` |
| POST | `/v1/posting/fbp/list` | `PostingFbpList` |
| POST | `/v1/posting/fbs/package-label/create` | `PostingAPI_CreateLabelBatch` |
| POST | `/v1/posting/fbs/product/traceable/attribute` | `PostingFbsProductTraceableAttribute` |
| POST | `/v1/posting/fbs/traceable/split` | `PostingFbsTraceableSplit` |
| POST | `/v1/posting/marks` | `PostingAPI_PostingMarks` |
| POST | `/v1/product/action/timer/status` | `ProductAPI_ActionTimerStatus` |
| GET | `/v1/product/certificate/accordance-types` | `ProductAPI_ProductCertificateAccordanceTypes` |
| POST | `/v1/product/certificate/product_status/list` | `ProductStatusList` |
| POST | `/v1/product/certification/list` | `ProductAPI_V1ProductCertificationList` |
| POST | `/v1/product/digital/stocks/import` | `DigitalProductAPI_StocksImport` |
| POST | `/v1/product/info/stocks-by-warehouse/fbo` | `GetProductInfoStocksByWarehouseFbo` |
| POST | `/v1/product/info/warehouse/stocks` | `ProductInfoWarehouseStocks` |
| POST | `/v1/product/info/wrong-volume` | `ProductAPI_ProductInfoWrongVolume` |
| POST | `/v1/product/placement-zone/info` | `ProductAPI_GetProductPlacementZoneInfo` |
| POST | `/v1/product/prices/details` | `ProductPricesDetails` |
| POST | `/v1/product/stairway-discount/by-quantity/get` | `ProductAPI_GetProductStairwayDiscountByQuantity` |
| POST | `/v1/product/stairway-discount/by-quantity/set` | `ProductAPI_SetProductStairwayDiscountByQuantity` |
| POST | `/v1/product/visibility/info` | `ProductVisibilityInfo` |
| POST | `/v1/product/visibility/set` | `ProductVisibilitySet` |
| POST | `/v1/question/answer/create` | `QuestionAnswer_Create` |
| POST | `/v1/question/answer/delete` | `QuestionAnswer_Delete` |
| POST | `/v1/question/answer/list` | `QuestionAnswer_List` |
| POST | `/v1/question/change-status` | `Question_ChangeStatus` |
| POST | `/v1/question/count` | `Question_Count` |
| POST | `/v1/question/info` | `Question_Info` |
| POST | `/v1/question/list` | `Question_List` |
| POST | `/v1/question/top-sku` | `Question_TopSku` |
| POST | `/v1/rating/index/fbs/info` | `RatingAPI_GetFBSRatingIndexInfoV1` |
| POST | `/v1/rating/index/fbs/posting/list` | `RatingAPI_ListFBSRatingIndexPostingsV1` |
| POST | `/v1/receipts/get` | `GetReceipt` |
| POST | `/v1/receipts/seller/list` | `ReceiptsSellerList` |
| POST | `/v1/receipts/upload` | `UploadReceipt` |
| POST | `/v1/removal/from-stock/list` | `GetSupplierReturnsSummaryReport` |
| POST | `/v1/removal/from-supply/list` | `GetSupplyReturnsSummaryReport` |
| POST | `/v1/report/marked-products-sales/create` | `CreateCompanyMarkedProductsSalesReport` |
| POST | `/v1/report/placement/by-products/create` | `CreatePlacementByProductsReport` |
| POST | `/v1/report/placement/by-supplies/create` | `CreatePlacementBySuppliesReport` |
| POST | `/v1/report/realization/posting/create` | `CreateCompanyFinanceRealizationPostingReport` |
| POST | `/v1/returns/rfbs/action/set` | `ReturnsAPI_ReturnsRfbsActionSet` |
| POST | `/v1/returns/settings/utilization/history` | `UtilizationHistory` |
| POST | `/v1/returns/settings/utilization/info` | `UtilizationInfo` |
| POST | `/v1/returns/settings/utilization/update` | `UtilizationUpdate` |
| POST | `/v1/roles` | `AccessAPI_RolesByToken` |
| POST | `/v1/search-queries/text` | `SearchQueriesAPI_SearchQueriesText` |
| POST | `/v1/search-queries/top` | `SearchQueriesAPI_SearchQueriesTop` |
| POST | `/v1/seller-actions/archive` | `SellerActionsArchive` |
| POST | `/v1/seller-actions/change-activity` | `SellerActionsChangeActivity` |
| POST | `/v1/seller-actions/create/discount` | `SellerActionsCreateDiscount` |
| POST | `/v1/seller-actions/create/discount-with-condition` | `SellerActionsCreateDiscountWithCondition` |
| POST | `/v1/seller-actions/create/installment` | `SellerActionsCreateInstallment` |
| POST | `/v1/seller-actions/create/multi-level-discount` | `SellerActionsCreateMultiLevelDiscount` |
| POST | `/v1/seller-actions/create/voucher` | `SellerActionsCreateVoucher` |
| POST | `/v1/seller-actions/list` | `SellerActionsList` |
| POST | `/v1/seller-actions/products/add` | `SellerActionsProductsAdd` |
| POST | `/v1/seller-actions/products/candidates` | `SellerActionsProductsCandidates` |
| POST | `/v1/seller-actions/products/delete` | `SellerActionsProductsDelete` |
| POST | `/v1/seller-actions/products/list` | `SellerActionsProductsList` |
| POST | `/v1/seller-actions/update/discount` | `SellerActionsUpdateDiscount` |
| POST | `/v1/seller-actions/update/discount-with-condition` | `SellerActionsUpdateDiscountWithCondition` |
| POST | `/v1/seller-actions/update/installment` | `SellerActionsUpdateInstallment` |
| POST | `/v1/seller-actions/update/multi-level-discount` | `SellerActionsUpdateMultiLevelDiscount` |
| POST | `/v1/seller-actions/update/voucher` | `SellerActionsUpdateVoucher` |
| POST | `/v1/seller-actions/voucher/get` | `SellerActionsVoucherGet` |
| POST | `/v1/seller/info` | `SellerAPI_SellerInfo` |
| POST | `/v1/seller/ozon-logistics/info` | `SellerAPI_SellerOzonLogisticsInfo` |
| POST | `/v1/supply-order/act/accept` | `SupplyOrderActAccept` |
| POST | `/v1/supply-order/act/accept/status` | `SupplyOrderActAcceptStatus` |
| POST | `/v1/supply-order/act/product/get` | `SupplyOrderActProductGet` |
| POST | `/v1/supply-order/act/summary/get` | `SupplyOrderActSummaryGet` |
| POST | `/v1/supply-order/content/update` | `SupplyOrderAPI_SupplyOrderContentUpdate` |
| POST | `/v1/supply-order/content/update/status` | `SupplyOrderAPI_SupplyOrderContentUpdateStatus` |
| POST | `/v1/supply-order/content/update/validation` | `SupplyOrderContentUpdateValidation` |
| POST | `/v1/supply-order/details` | `SupplyOrderAPI_SupplyOrderDetails` |
| POST | `/v1/warehouse/archive` | `ArchiveWarehouseFBS` |
| POST | `/v1/warehouse/erfbs/aggregator/create` | `WarehouseERFBSAggregatorCreate` |
| POST | `/v1/warehouse/erfbs/aggregator/delivery-method/update` | `WarehouseERFBSAggregatorDeliveryMethodUpdate` |
| POST | `/v1/warehouse/erfbs/non-integrated/create` | `WarehouseERFBSNonIntegratedCreate` |
| POST | `/v1/warehouse/erfbs/non-integrated/delivery-method/update` | `WarehouseERFBSNonIntegratedDeliveryMethodUpdate` |
| POST | `/v1/warehouse/erfbs/update` | `WarehouseERFBSUpdate` |
| POST | `/v1/warehouse/fbo/seller/list` | `WarehouseFboSellerList` |
| POST | `/v1/warehouse/fbs/create` | `WarehouseAPI_CreateWarehouseFBS` |
| POST | `/v1/warehouse/fbs/create/drop-off/list` | `WarehouseAPI_ListDropOffPointsForCreateFBSWarehouse` |
| POST | `/v1/warehouse/fbs/create/drop-off/timeslot/list` | `WarehouseFbsCreateDropOffTimeslotList` |
| POST | `/v1/warehouse/fbs/create/pick-up/timeslot/list` | `WarehouseFbsCreatePickUpTimeslotList` |
| POST | `/v1/warehouse/fbs/create/return-point/list` | `WarehouseFBSCreateReturnPointList` |
| POST | `/v1/warehouse/fbs/first-mile/update` | `UpdateWarehouseFBSFirstMile` |
| POST | `/v1/warehouse/fbs/pickup/courier/cancel` | `WarehouseFbsPickUpCourierCancel` |
| POST | `/v1/warehouse/fbs/pickup/courier/create` | `WarehouseFbsPickUpCourierCreate` |
| POST | `/v1/warehouse/fbs/pickup/history/list` | `WarehouseFbsPickUpHistoryList` |
| POST | `/v1/warehouse/fbs/pickup/planning/list` | `WarehouseFbsPickUpPlanningList` |
| POST | `/v1/warehouse/fbs/return-mile/check` | `WarehouseFbsReturnMileCheck` |
| POST | `/v1/warehouse/fbs/return-mile/info` | `WarehouseFBSReturnMileInfo` |
| POST | `/v1/warehouse/fbs/update` | `UpdateWarehouseFBS` |
| POST | `/v1/warehouse/fbs/update/drop-off/list` | `WarehouseAPI_ListDropOffPointsForUpdateFBSWarehouse` |
| POST | `/v1/warehouse/fbs/update/drop-off/timeslot/list` | `WarehouseFbsUpdateDropOffTimeslotList` |
| POST | `/v1/warehouse/fbs/update/pick-up/timeslot/list` | `WarehouseFbsUpdatePickUpTimeslotList` |
| POST | `/v1/warehouse/fbs/update/return-point/list` | `WarehouseFBSUpdateReturnPointList` |
| POST | `/v1/warehouse/invalid-products/get` | `WarehouseInvalidProductsGet` |
| POST | `/v1/warehouse/operation/status` | `GetWarehouseFBSOperationStatus` |
| POST | `/v1/warehouse/ozon/list` | `WarehouseOZONList` |
| POST | `/v1/warehouse/rfbs/pause` | `WarehouseRfbsPause` |
| POST | `/v1/warehouse/rfbs/unpause` | `WarehouseRfbsUnpause` |
| POST | `/v1/warehouse/unarchive` | `UnarchiveWarehouseFBS` |
| POST | `/v1/warehouse/warehouses-with-invalid-products` | `WarehouseWithInvalidProducts` |
| POST | `/v2/actions/discounts-task/list` | `GetDiscountTaskListV2` |
| POST | `/v2/cargoes/create/info` | `CargoesCreateInfoV2` |
| POST | `/v2/cargoes/delete` | `CargoesDeleteV2` |
| POST | `/v2/cargoes/delete/status` | `CargoesDeleteStatusV2` |
| POST | `/v2/cargoes/get` | `CargoesGetV2` |
| POST | `/v2/carriage/delivery/list` | `CarriageAPI_CarriageDeliveryListV2` |
| POST | `/v2/cluster/list` | `DraftClusterList` |
| POST | `/v2/conditional-cancellation/approve` | `CancellationAPI_ConditionalCancellationApproveV2` |
| POST | `/v2/conditional-cancellation/list` | `CancellationAPI_GetConditionalCancellationListV2` |
| POST | `/v2/conditional-cancellation/reject` | `CancellationAPI_ConditionalCancellationRejectV2` |
| POST | `/v2/delivery-method/list` | `WarehouseAPI_DeliveryMethodListV2` |
| POST | `/v2/delivery/checkout` | `DeliveryCheckout` |
| POST | `/v2/draft/create/info` | `DraftCreateInfo` |
| POST | `/v2/draft/supply/create` | `DraftSupplyCreate` |
| POST | `/v2/draft/supply/create/status` | `DraftSupplyCreateStatus` |
| POST | `/v2/draft/timeslot/info` | `DraftTimeslotInfo` |
| POST | `/v2/order/create` | `OrderAPI_OrderCreate` |
| POST | `/v2/polygon/bind` | `PolygonBind` |
| POST | `/v2/posting/digital/list` | `PostingDigitalList` |
| POST | `/v2/product/info/stocks-by-warehouse/fbs` | `ProductAPI_GetProductInfoStocksByWarehouseFbsV2` |
| POST | `/v2/supply-order/timeslot/list` | `SupplyOrderTimeslotList` |
| POST | `/v2/warehouse/list` | `WarehouseListV2` |
| POST | `/v3/chat/list` | `ChatAPI_ChatListV3` |
| POST | `/v3/supply-order/get` | `SupplyOrderGet` |
| POST | `/v3/supply-order/list` | `SupplyOrderList` |
