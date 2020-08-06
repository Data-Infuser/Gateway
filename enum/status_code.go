package enum

import grpc_author "gitlab.com/promptech1/infuser-gateway/infuser-protobuf/gen/proto/author"

type ResCode int32

const (
	Valid ResCode = 0
	InternalException ResCode = -1
	ParameterException ResCode = -2
	UnregisteredService ResCode = -3
	UnregisteredToken ResCode = -4
	TerminatedService ResCode = -9
	LimitExceeded ResCode = -10
	Unauthorized ResCode = -401
	Unknown ResCode = -999
)

var (
	ProtoCodeMap = map[grpc_author.ApiAuthRes_Code]ResCode {
		grpc_author.ApiAuthRes_VALID: Valid,
		grpc_author.ApiAuthRes_INTERNAL_EXCEPTION: InternalException,
		grpc_author.ApiAuthRes_PARAMETER_EXCEPTION: ParameterException,
		grpc_author.ApiAuthRes_UNREGISTERED_SERVICE: UnregisteredService,
		grpc_author.ApiAuthRes_UNREGISTERED_TOKEN: UnregisteredToken,
		grpc_author.ApiAuthRes_TERMINATED_SERVICE: TerminatedService,
		grpc_author.ApiAuthRes_LIMIT_EXCEEDED: LimitExceeded,
		grpc_author.ApiAuthRes_UNAUTHORIZED: Unauthorized,
		grpc_author.ApiAuthRes_UNKNOWN: Unknown,
	}
)

var ResCodes = [...]ResCode {
	Valid,
	InternalException,
	ParameterException,
	UnregisteredService,
	UnregisteredToken,
	TerminatedService,
	LimitExceeded,
	Unauthorized,
	Unknown,
}

func (c ResCode) Message() string {
	switch c {
	case Valid:
		return "정상"
	case InternalException:
		return "시스템 내부 오류가 발생하였습니다."
	case ParameterException:
		return "요청하신 파라미터가 적합하지 않습니다."
	case UnregisteredService:
		return "등록되지 않은 서비스 입니다."
	case UnregisteredToken:
		return "등록되지 않은 인증키 입니다."
	case TerminatedService:
		return "종료된 서비스 입니다."
	case LimitExceeded:
		return "트래픽 허용 횟수를 초과하였습니다."
	case Unauthorized:
		return "유효하지 않은 인증키 입니다."
	default:
		return "UNKNOWN"
	}
}

func (c ResCode) HttpCode() int {
	switch c {
	case Valid:
		return 200
	case InternalException:
		return 400
	case ParameterException:
		return 400
	case UnregisteredService:
		return 400
	case UnregisteredToken:
		return 400
	case TerminatedService:
		return 400
	case LimitExceeded:
		return 400
	case Unauthorized:
		return 401
	default:
		return 400
	}
}

func (c ResCode) Valid() bool {
	return c == Valid
}

func FindResCode(grpcCode grpc_author.ApiAuthRes_Code) ResCode {
	if val, ok := ProtoCodeMap[grpcCode]; ok {
		return val
	}

	return Unknown
}