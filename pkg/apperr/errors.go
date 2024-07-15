package apperr

import (
	"google.golang.org/grpc/codes"
)

var (
	// Common
	// 400 error
	UnknownMediaType          = newAppErr("CM400001", "unknown media type", codes.InvalidArgument)
	UnknownGender             = newAppErr("CM400001", "unknown gender", codes.InvalidArgument)
	UnknownAspType            = newAppErr("CM400003", "unknown asp type", codes.InvalidArgument)
	UnknownCommissionStatus   = newAppErr("CM400004", "unknown commission status", codes.InvalidArgument)
	UnknownMutation           = newAppErr("CM400005", "unknown mutation", codes.InvalidArgument)
	UnknownAffiliatorValidity = newAppErr("CM400006", "unknown affiliator validity", codes.InvalidArgument)
	UnknownCommonMetadata     = newAppErr("CM400007", "unknown common metadata", codes.InvalidArgument)
	// 403 error
	PermissionDenied = newAppErr("CM403001", "permission denied", codes.PermissionDenied)
	RequestCanceled  = newAppErr("CM499000", "request canceled", codes.Canceled)
	// 500 error
	RequestIdNotFound   = newAppErr("CM500001", "request id not found", codes.Internal)
	RequestedByNotFound = newAppErr("CM500002", "requested by not found", codes.Internal)
	UnknownError        = newAppErr("CM500003", "unknown error", codes.Internal)
	UnavailableError    = newAppErr("CM503001", "unavailable", codes.Unavailable)
	// UnKnownError =
	// Affiliate Item
	AffiliateItemValidationError               = newAppErr("AI400000", "validation error", codes.InvalidArgument)
	AffiliateItemRakutenInvalidAffiliatorError = newAppErr("AI400001", "rakuten invalid affiliator", codes.InvalidArgument)
	AffiliateItemInvalidAffiliator             = newAppErr("AI401000", "invalid affiliator", codes.Unauthenticated)
	AffiliateItemInvalidRakutenId              = newAppErr("AI401001", "invalid rakuten affiliate id.", codes.Unauthenticated)
	AffiliateItemNotFound                      = newAppErr("AI404000", "not found", codes.NotFound)
	AffiliateItemRequestsTimeoutError          = newAppErr("AI408000", "requests timeout", codes.Aborted)
	AffiliateItemTooManyRequestError           = newAppErr("AI429000", "too many request", codes.DeadlineExceeded)
	AffiliateItemInternalError                 = newAppErr("AI500000", "internal error", codes.Internal)
	AffiliateItemRepositoryError               = newAppErr("AI500001", "repository error", codes.Internal)
	AffiliateItemFailedToApplyCommissionRate   = newAppErr("AI500002", "failed to apply commission rate", codes.Internal)
	AffiliateItemFailedToImportDatafeed        = newAppErr("AI500002", "failed to import datafeed", codes.Internal)
	AffiliateItemFailedToImportCAWiseItems     = newAppErr("AI500002", "failed to import cawise items", codes.Internal)
	AffiliateNotRegisteredTieupType            = newAppErr("AI500004", "not registered tieup type", codes.Internal)

	// Affiliator
	AffiliatorValidationError                 = newAppErr("AF400000", "validation error", codes.InvalidArgument)
	AffiliatorDateOfBirthParseError           = newAppErr("AF400001", "date of birth parse error", codes.InvalidArgument)
	AffiliatorUpperLimitError                 = newAppErr("AF400002", "upper limit error", codes.InvalidArgument)
	AffiliatorNotMatchRakutenAffiliateIDError = newAppErr("AF400003", "rakuten affiliate id is not match", codes.InvalidArgument)
	AffiliatorInvalidVerificationTokenError   = newAppErr("AF400005", "invalid verification token", codes.InvalidArgument)
	AffiliatorNotFoundError                   = newAppErr("AF404000", "not found", codes.NotFound)
	AffiliatorNotMatchResultCountError        = newAppErr("AF404001", "parameter length and result count do not match", codes.NotFound)
	AffiliatorNotificationBoxNotFoundError    = newAppErr("AF404002", "notification box not found", codes.NotFound)
	AffiliatorNotificationPopupNotFoundError  = newAppErr("AF404003", "notification popup not found", codes.NotFound)
	AffiliatorAlreadyExistsError              = newAppErr("AF409000", "conflict", codes.AlreadyExists)
	AffiliatorAlreadySignedUpError            = newAppErr("AF409001", "already signed up", codes.AlreadyExists)
	AffiliatorInternalError                   = newAppErr("AF500000", "internal error", codes.Internal)
	AffiliatorUnavailableError                = newAppErr("AF503000", "unavailable error", codes.Unavailable)

	// Notification
	NotificationValidationError               = newAppErr("NT400000", "validation error", codes.InvalidArgument)
	NotificationTemplateUnknownVariable       = newAppErr("NT400001", "template unknown variable error", codes.InvalidArgument)
	NotificationScheduleInsufficientVariables = newAppErr("NT400002", "schedule has insufficient variables", codes.InvalidArgument)
	NotificationNotFound                      = newAppErr("NT404000", "not found", codes.NotFound)
	NotificationNicknameNotFound              = newAppErr("NT404001", "nickname not found", codes.NotFound)
	NotificationInternalError                 = newAppErr("NT500000", "internal error", codes.Internal)
	NotificationMediaServerError              = newAppErr("NT500001", "media server error", codes.Internal)

	// Commission
	CommissionUnknownError                  = newAppErr("CO100000", "unknown error", codes.Unknown)
	CommissionEnvVarError                   = newAppErr("CO100001", "invalid environment variable", codes.InvalidArgument)
	CommissionConfigError                   = newAppErr("CO100002", "invalid config", codes.InvalidArgument)
	CommissionValidationError               = newAppErr("CO100010", "validation error", codes.InvalidArgument)
	CommissionRetryableError                = newAppErr("CO100010", "retryable error", codes.Unknown)
	CommissionUnretryableError              = newAppErr("CO100011", "unretryable error", codes.Unknown)
	CommissionInternalError                 = newAppErr("CO100012", "internal error", codes.Internal)
	CommissionHttpServerError               = newAppErr("CO100020", "occurred http server", codes.Internal)
	CommissionHttpClientError               = newAppErr("CO100021", "occurred http client", codes.Internal)
	CommissionIOError                       = newAppErr("CO100022", "io error", codes.Internal)
	CommissionAspUnavailable                = newAppErr("CO100100", "unavailable asp service", codes.Unavailable)
	CommissionAspUnexpectedResultData       = newAppErr("CO100101", "asp was respond that unexpected result data", codes.InvalidArgument)
	CommissionGcpPubSubUnavailable          = newAppErr("CO100110", "unavailable pub/sub of GCP", codes.Unavailable)
	CommissionMeasurementUnavailable        = newAppErr("CO100120", "unavailable measurement context", codes.Unavailable)
	CommissionAffiliatorUnavailable         = newAppErr("CO100121", "unavailable affiliator context ", codes.Unavailable)
	CommissionEstimateError                 = newAppErr("CO100200", "commission estimated but unexpected calculation result", codes.Internal)
	CommissionDotmoneyError                 = newAppErr("CO100300", "failed dotmoney operation", codes.Internal)
	CommissionDatabaseError                 = newAppErr("CO100301", "failed database operation", codes.Internal)
	CommissionTaskTimeoutError              = newAppErr("CO100400", "task timeout", codes.Aborted)
	CommissionTooManyTaskFailedError        = newAppErr("CO100401", "too many failed task set hold status", codes.Internal)
	CommissionAmazonNetRateCalculationError = newAppErr("CO100500", "failed to calculate net rate", codes.Internal)

	// Commission/Dotmoney
	CommissionDotmoneyConstraintViolation    = newAppErr("CO400100", "constraint violation", codes.Unknown)
	CommissionDotmoneyAccountNotFound        = newAppErr("CO400101", "account not found", codes.Unknown)
	CommissionDotmoneyAccountInvalidStatus   = newAppErr("CO400102", "account invalid status", codes.Unknown)
	CommissionDotmoneyBalanceOverflow        = newAppErr("CO400103", "balance overflow", codes.Unknown)
	CommissionDotmoneyRequestIdAlreadyExists = newAppErr("CO400104", "request id already exists", codes.Unknown)
	CommissionDotmoneyUnauthorized           = newAppErr("CO401101", "unauthorized", codes.Unknown)
	CommissionDotmoneyForbidden              = newAppErr("CO403101", "forbidden", codes.Unknown)
	CommissionDotmoneyInternalServerError    = newAppErr("CO500100", "internal server error", codes.Unknown)
	CommissionDotmoneyUnknownError           = newAppErr("CO500101", "unknown error", codes.Unknown)
	CommissionDotmoneyServiceUnavailable     = newAppErr("CO503100", "service unavailable", codes.Unknown)
	CommissionNotFound                       = newAppErr("CO404000", "not found", codes.NotFound)

	// Measurement
	MeasurementValidationError                = newAppErr("MM400000", "validation error", codes.InvalidArgument)
	MeasurementTooMuchAffiliateTagsError      = newAppErr("MM400002", "too much affiliate tags", codes.InvalidArgument)
	MeasurementCanNotUpdateAffiliateTagError  = newAppErr("MM400003", "can not update affiliate tag", codes.InvalidArgument)
	MeasurementRakutenItemInAmemberEntryError = newAppErr("MM400004", "rakuten item in amember entry", codes.InvalidArgument)
	MeasurementNotFound                       = newAppErr("MM404000", "not found", codes.NotFound)
	MeasurementAffiliateTagNotFound           = newAppErr("MM404001", "affiliate tag not found", codes.NotFound)
	MeasurementClickNotFound                  = newAppErr("MM404002", "click not found", codes.NotFound)
	MeasurementAffiliateItemNotFound          = newAppErr("MM404003", "affiliate-item not found", codes.NotFound)
	MeasurementAffiliateItemUnavailable       = newAppErr("MM503001", "unavailable affiliate-item context", codes.Unavailable)
	MeasurementAffiliatorUnavailable          = newAppErr("MM503002", "unavailable affiliator context ", codes.Unavailable)
	MeasurementInternalError                  = newAppErr("MM500000", "internal error", codes.Internal)

	// Personal Information
	PersonalInformationValidationError       = newAppErr("PI400000", "validation error", codes.InvalidArgument)
	PersonalInformationDateOfBirthParseError = newAppErr("PI400001", "date of birth parse error", codes.InvalidArgument)
	PersonalInformationNotFound              = newAppErr("PI404000", "not found", codes.NotFound)
	PersonalInformationInternalError         = newAppErr("PI500000", "internal error", codes.Internal)
	PersonalInformationDatabaseError         = newAppErr("PI500001", "database error", codes.Internal)

	// signup
	SignUpValidationError               = newAppErr("SU400000", "validation error", codes.InvalidArgument)
	SignUpAlreadyApprovedError          = newAppErr("SU400001", "already signed up", codes.InvalidArgument)
	SignUpInvalidVerificationTokenError = newAppErr("SU400002", "invalid verification token", codes.InvalidArgument)
	SignupNotFoundError                 = newAppErr("SU404000", "not found", codes.NotFound)
	SignUpAlreadySignedUpError          = newAppErr("SU409000", "already signed up", codes.AlreadyExists)
	SignUpInternalError                 = newAppErr("SU500000", "internal error", codes.Internal)

	// Summary
	SummaryValidationError = newAppErr("SM400000", "validation error", codes.InvalidArgument)
	SummaryNotFound        = newAppErr("SM404000", "not found", codes.NotFound)
	SummaryInternalError   = newAppErr("SM500000", "internal error", codes.Internal)
	SummaryDatabaseError   = newAppErr("SM500001", "database error", codes.Internal)

	// Media Ameba
	MediaAmebaValidationError        = newAppErr("MA400000", "validation error", codes.InvalidArgument)
	MediaAmebaAdsUserDenyReceiving   = newAppErr("MA400001", "ads user deny receiving", codes.FailedPrecondition)
	MediaAmebaAdsInvalidUserID       = newAppErr("MA400002", "ads invalid user id", codes.InvalidArgument)
	MediaAmebaNotFound               = newAppErr("MA404000", "not found", codes.NotFound)
	MediaAmebaAdsUserTokenNotFound   = newAppErr("MA404001", "ads user token not found", codes.NotFound)
	MediaAmebaInternalError          = newAppErr("MA500000", "internal error", codes.Internal)
	MediaAmebaExternalApiServerError = newAppErr("MA500001", "external api server error", codes.Internal)

	// Image Proxy
	ImageProxyValidationError = newAppErr("IP400000", "validation error", codes.InvalidArgument)
	ImageProxyNotFound        = newAppErr("IP404000", "not found", codes.NotFound)
	ImageProxyInternalError   = newAppErr("IP500000", "internal error", codes.Internal)

	// ASP Rakuten
	AspRakutenValidationError                 = newAppErr("AR400000", "validation error", codes.InvalidArgument)
	AspRakutenNotMatchRakutenAffiliateIDError = newAppErr("AR400001", "rakuten affiliate id is not match", codes.InvalidArgument)
	AspRakutenInvalidAffiliatorError          = newAppErr("AR400002", "rakuten invalid affiliator", codes.InvalidArgument)
	AspRakutenNotFound                        = newAppErr("AR404000", "not found", codes.NotFound)
	AspRakutenTooManyRequestError             = newAppErr("AR429000", "too many request", codes.DeadlineExceeded)
	AspRakutenInternalError                   = newAppErr("AR500000", "internal error", codes.Internal)
	AspRakutenRepositoryError                 = newAppErr("AR500001", "repository error", codes.Internal)

	// Offer Item
	OfferItemValidationError                    = newAppErr("OI400000", "validation error", codes.InvalidArgument)
	OfferItemUpsertUnableError                  = newAppErr("OI400001", "unable to upsert offer item", codes.InvalidArgument)
	OfferItemScheduleUpsertUnableError          = newAppErr("OI400002", "unable to upsert schedule", codes.InvalidArgument)
	OfferItemAssigneeUpsertUnableError          = newAppErr("OI400003", "unable to upsert assignee", codes.InvalidArgument)
	OfferItemMailSettingUpsertUnableError       = newAppErr("OI400004", "unable to upsert mail setting", codes.InvalidArgument)
	OfferItemSampleProductUpsertUnableError     = newAppErr("OI400005", "unable to upsert sample product", codes.InvalidArgument)
	OfferItemFormAlreadyAnsweredError           = newAppErr("OI400005", "form already answered", codes.FailedPrecondition)
	OfferItemScheduleExpiredError               = newAppErr("OI400006", "schedule expired", codes.FailedPrecondition)
	OfferItemNoNeedNotificationTaskCreatedError = newAppErr("OI400007", "no need notification task created", codes.FailedPrecondition)
	OfferItemNotFoundError                      = newAppErr("OI404000", "not found", codes.NotFound)
	OfferItemAffiliateItemNotFoundError         = newAppErr("OI404001", "affiliate-item not found", codes.NotFound)
	OfferItemBloggerPropertyNotFoundError       = newAppErr("OI404002", "blogger property not found", codes.NotFound)
	OfferItemInternalError                      = newAppErr("OI500000", "internal error", codes.Internal)
	OfferItemSendMailPreCheckFailedError        = newAppErr("OI500001", "validation before sending mail failed", codes.Internal)
	OfferItemAffiliateItemUnavailableError      = newAppErr("OI503000", "unavailable affiliate-item context", codes.Unavailable)
	OfferItemMediaAmebaUnavailableError         = newAppErr("OI503001", "unavailable media-ameba context", codes.Unavailable)
	OfferItemSignupUnavailableError             = newAppErr("OI503002", "unavailable signup context", codes.Unavailable)
	OfferItemAffiliatorUnavailableError         = newAppErr("OI503003", "unavailable affiliator context", codes.Unavailable)
	OfferItemMeasurementUnavailableError        = newAppErr("OI503004", "unavailable measurement context", codes.Unavailable)

	// Exporter
	ExporterNotFoundError                 = newAppErr("EX404000", "not found", codes.NotFound)
	ExporterInternalError                 = newAppErr("EX500000", "internal error", codes.Internal)
	ExporterAffiliateItemUnavailableError = newAppErr("EX503000", "unavailable affiliate-item context", codes.Unavailable)

	// Special Select
	SpecialSelectValidationError       = newAppErr("SS400000", "validation error", codes.InvalidArgument)
	SpecialSelectInvalidAffiliator     = newAppErr("SS401000", "invalid affiliator", codes.Unauthenticated)
	SpecialSelectNotFoundError         = newAppErr("SS404000", "not found", codes.NotFound)
	SpecialSelectApplicationLimitError = newAppErr("SS404001", "the application limit has been reached", codes.NotFound)
	SpecialSelectInternalError         = newAppErr("SS500000", "internal error", codes.Internal)
	SpecialSelectDatabaseError         = newAppErr("SS500001", "database error", codes.Internal)

	// Commerce
	CommerceValidationError = newAppErr("CC400000", "validation error", codes.InvalidArgument)
	CommerceOverLimitError  = newAppErr("CC400001", "over limit", codes.InvalidArgument)
	CommerceNotFoundError   = newAppErr("CC404000", "not found", codes.NotFound)
	CommerceAlreadyExists   = newAppErr("CC409000", "already exists", codes.AlreadyExists)
	CommerceInternalError   = newAppErr("CC500000", "internal error", codes.Internal)
	CommerceDatabaseError   = newAppErr("CC500001", "database error", codes.Internal)
)
