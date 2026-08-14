// 错误码定义
package errcode

// 通用: 基本错误
// Code must start with 1xxxxx
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001
	// ErrUnknown - 500: Internal server error.
	ErrUnknown
	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind
	// ErrValidation - 400: Validation failed.
	ErrValidation
	// ErrErrTokenInvalid - 401: Token invalid.
	ErrErrTokenInvalid
	// ErrInternalServer - 500: Server Exception.
	ErrInternalServer
	// ErrTooManyRequests - 429: Too many requests
	ErrTooManyRequests
)

// 通用：数据库类错误
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

// 通用：认证授权类错误
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201
	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid
	// ErrExpired - 401: Token expired.
	ErrExpired
	// ErrTokenRevokeFailed - 401: token revoke failed
	ErrTokenRevokeFailed
	// ErrTokenRedisStateSetFailed - 401: redis set state failed
	ErrTokenRedisStateSetFailed
	// ErrTokenCreateFailed - 401: token create failed
	ErrTokenCreateFailed
	// ErrTokenParseFailed - 401: token parse failed
	ErrTokenParseFailed
	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader
	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader
	// ErrPasswordIncorrect - 401: Password was incorrect.
	ErrPasswordIncorrect
	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied
	// ErrCasbinSet - 403: Permission denied.
	ErrCasbinSet
)

// 通用：编解码类错误
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301
	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed
	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON
	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON
	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON
	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml
	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml
	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)

// User(用户)相关错误
const (
	// ErrUserNotFound - 404: User not found
	ErrUserNotFound = iota + 110001
	// ErrUserPasswordFault - 404: Password is fault
	ErrUserPasswordFault
	// ErrUserAlreadyExist - 400: User already exist
	ErrUserAlreadyExist
	// ErrUserForbidden - 400: User is forbidden
	ErrUserForbidden
	// ErrUserRegisteFailed - 400: User registe failed
	ErrUserRegisteFailed
	// ErrUserGetList - 400: Get list failed
	ErrUserGetList
	// ErrUserSwitchRole - 400: Switch role failed
	ErrUserSwitchRole
	// ErrUserSetAuthorities - 400: set authorities failed
	ErrUserSetAuthorities
	// ErrUserDeleteSelf - 400: delete self failed
	ErrUserDeleteSelf
	// ErrUserDeleteFailed - 400: user delete failed
	ErrUserDeleteFailed
	// ErrUserSetInfoFailed - 400: user set info failed
	ErrUserSetInfoFailed
	// ErrUserSetSelfFailed - 400: user set self failed
	ErrUserSetSelfFailed
	// ErrUserResetPassworkFailed - 400: User reset password failed
	ErrUserResetPassworkFailed
)

// 密钥相关错误
const (
	// ErrReachMaxCount - 400: Secret reach the max count
	ErrReachMaxCount = iota + 110101
	// ErrSecretNotFound - 404: Secret not found
	ErrSecretNotFound
)

// 插件相关错误
const (
	// ErrLimitedIP - 500: Limited IP
	ErrLimitedIP = iota + 110201
)

// 权限相关错误
const (
	// ErrAuthorityNotFound - 404: Authority not found
	ErrAuthorityNotFound = iota + 110301
	// ErrGetAuthorities - 500: GET failed.
	ErrGetAuthorities
	// ErrAuthorityAlreadyExist - 400: Authority already exist
	ErrAuthorityAlreadyExist
	// ErrAuthorityCreateFailed - 400: Authority create failed
	ErrAuthorityCreateFailed
	// ErrAuthorityDeleteFailed - 400: Authority delete failed
	ErrAuthorityDeleteFailed
	// ErrAuthorityUpdateFailed - 400: Authority update failed
	ErrAuthorityUpdateFailed
	// ErrAuthoritySearchFailed - 400: Authority Search failed
	ErrAuthoritySearchFailed
	// ErrAuthorityAlreadyAssignedToUser - 400: Authority already assigned to user
	ErrAuthorityAlreadyAssignedToUser
	// ErrAuthorityHasChildren - 400: Authority has children
	ErrAuthorityHasChildren
	// ErrAuthorityNotAssignedToUser - 400: Authority not assigned to user
	ErrAuthorityNotAssignedToUser
	// ErrAuthorityNotPermissionGet - 403: Authority not permission get
	ErrAuthorityNotPermissionGet
)

// 菜单相关错误
const (
	// ErrMenuNotFound - 200: Menu not found
	ErrMenuNotFound = iota + 110401
	// ErrMenuAlreadyExist - 200: Menu already exist
	ErrMenuAlreadyExist
	// ErrAuthorityMenuNotFound - 200: Authority has no corresponding menu
	ErrAuthorityMenuNotFound
	// ErrMenuSetFailed - 200: Menu set failed
	ErrMenuSetFailed
	// ErrMenuHasChildren - 200: Menu has children and cannot be deleted
	ErrMenuHasChildren
	// ErrMenuAsHomePage - 200: Menu is being used as home page by some authority and cannot be deleted
	ErrMenuAsHomePage
	// ErrMenuSetHomePageFailed - 200: Menu set home page failed
	ErrMenuSetHomePageFailed
	// ErrMenuDeleteFailed - 200: Menu delete failed
	ErrMenuDeleteFailed
)

// Api相关错误
const (
	// ErrApiNotFound - 404: Api not found
	ErrApiNotFound = iota + 110501
	// ErrApiAlreadyExist - 400: Api already exist
	ErrApiAlreadyExist
	// ErrApiCreateFailed - 400: Api create failed
	ErrApiCreateFailed
	// ErrApiDeleteFailed - 400: Api delete failed
	ErrApiDeleteFailed
	// ErrApiUpdateFailed - 400: Api update failed
	ErrApiUpdateFailed
	// ErrApiIgnoreFailed - 400: Api ignore failed
	ErrApiIgnoreFailed
	// ErrApiSyncFailed - 400: Api sync failed
	ErrApiSyncFailed
	// ErrApiRegisterListFailed - 400: Api register list failed
	ErrApiRegisterListFailed
	// ErrApiGroupFailed - 400: Api group failed
	ErrApiGroupFailed
	// ErrApiSyncParamsFailed - 400: Api sync params failed
	ErrApiSyncParamsFailed
)
