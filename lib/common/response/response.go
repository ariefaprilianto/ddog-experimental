package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ariefaprilianto/ddog-experimental/lib/common/custerr"
	"github.com/ariefaprilianto/ddog-experimental/lib/common/err"
)

var (
	ErrBadRequest                  = errors.New("Bad request")
	ErrForbiddenResource           = errors.New("Forbidden resource")
	ErrForbiddenNakama             = errors.New("Terms and conditions not accepted")
	ErrNotFound                    = errors.New("Not Found")
	ErrPreConditionFailed          = errors.New("Precondition failed")
	ErrInternalServerError         = errors.New("Internal server error")
	ErrTimeoutError                = errors.New("Timeout error")
	ErrAlreadyRegistered           = errors.New("User already registered")
	ErrNoLinkerExists              = errors.New("No Linker exist")
	ErrNoValidUserFound            = errors.New("No Valid User Found")
	ErrTxnAlreadyDone              = errors.New("Transaction Already Done")
	ErrExcessRefundAmount          = errors.New("Excess Refund Amount")
	ErrRefundAlreadyDone           = errors.New("Refund already done")
	ErrInvalidTxnId                = errors.New("Invalid Transaction Id")
	ErrDuplicateReq                = errors.New("Duplicate Request")
	ErrSignatureMisMatch           = errors.New("Signature MisMatch")
	ErrRefundAgainstInvalidTxnType = errors.New("Refund not allowed against Txn with Non-Payment type")
	ErrSaldoRefundAlreadyProcessed = errors.New("Saldo Redund Already processed")
	ErrNoUnsettledBillFound        = errors.New("No Unsettled Bill Found")
	ErrInvalidOTP                  = errors.New("Invalid OTP")
	ErrPreConditionFailedTest      = errors.New("Precondition failed test")
	ErrTokopediaEmail              = errors.New("Tokopedia email id can not be used")
	ErrNoBillFound                 = errors.New("No Bill found")
	ErrBillDueDatePassed           = errors.New("Last date to pay bill is over, pay your bill in next billing cycle.")
	ErrNoOVOWallet                 = errors.New("OVO wallet not found")
	ErrInvalidToken                = errors.New("Invalid auth token")
	ErrRequiredToken               = errors.New("Auth token is required")
	ErrUnauthorized                = errors.New("No Authorization Found")
)

const (
	STATUSCODE_GENERICSUCCESS                   = "200000"
	STATUSCODE_BADREQUEST                       = "400000"
	STATUSCODE_NOLINKEREXIST                    = "411553"
	STATUSCODE_UNAUTHORIZED                     = "401000"
	STATUS_FORBIDDEN                            = "403000"
	STATUSCODE_NOT_FOUND                        = "404000"
	STATUSCODE_GENERIC_PRECONDITION_FAILED      = "412000" // todo error code change from 412 to 400 dur to nginx issue
	STATUSCODE_INTERNAL_ERROR                   = "500000"
	STATUSCODE_TIMEOUT_ERROR                    = "504000"
	STATUSCODE_ALREADY_REGISTERED               = "400001"
	STATUSCODE_TX_ALREADY_DONE                  = "400003"
	STATUSCODE_GENERIC_PRECONDITION_FAILED_TEST = "412000"
	STATUSCODE_REQUIRED_TOKEN                   = "401000"
	STATUS_DATA_EMPTY                           = "2000001"
)

func GetErrorCode(err error) int {
	code, _ := strconv.ParseInt(GetErrorCodeStr(err), 10, 64)
	return int(code)
}

func GetErrorCodeStr(err error) string {
	switch err.(type) {
	case custerr.ErrChain:
		errType := err.(custerr.ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	switch err {
	case ErrBadRequest:
		return STATUSCODE_BADREQUEST
	case ErrForbiddenResource:
		return STATUS_FORBIDDEN
	case ErrNotFound:
		return STATUSCODE_NOT_FOUND
	case ErrPreConditionFailed:
		return STATUSCODE_GENERIC_PRECONDITION_FAILED
	case ErrInternalServerError:
		return STATUSCODE_INTERNAL_ERROR
	case ErrTimeoutError:
		return STATUSCODE_TIMEOUT_ERROR
	case ErrAlreadyRegistered:
		return STATUSCODE_ALREADY_REGISTERED
	case ErrDuplicateReq:
		return STATUSCODE_ALREADY_REGISTERED
	case ErrNoLinkerExists:
		return STATUSCODE_NOLINKEREXIST
	case nil:
		return STATUSCODE_GENERICSUCCESS
	case ErrNoValidUserFound:
		return STATUSCODE_BADREQUEST
	case ErrTxnAlreadyDone:
		return STATUSCODE_TX_ALREADY_DONE
	case ErrInvalidOTP:
		return STATUSCODE_BADREQUEST
	case ErrForbiddenNakama:
		return STATUSCODE_BADREQUEST
	case ErrPreConditionFailedTest:
		return STATUSCODE_GENERIC_PRECONDITION_FAILED_TEST
	case ErrTokopediaEmail:
		return STATUSCODE_BADREQUEST
	case ErrBillDueDatePassed:
		return STATUSCODE_GENERIC_PRECONDITION_FAILED
	case ErrNoBillFound:
		return STATUSCODE_GENERIC_PRECONDITION_FAILED
	case ErrNoOVOWallet:
		return STATUSCODE_BADREQUEST
	case ErrInvalidToken:
		return STATUS_FORBIDDEN
	case ErrRequiredToken:
		return STATUSCODE_REQUIRED_TOKEN
	case ErrUnauthorized:
		return STATUSCODE_UNAUTHORIZED
	default:
		return STATUSCODE_INTERNAL_ERROR
	}
}

func GetHTTPCode(code string) int {
	s := code[0:3]
	i, _ := strconv.Atoi(s)
	return i
}

type JSONResponse struct {
	Code         string                 `json:"code"`
	Message      string                 `json:"message,omitempty"`
	ErrorMessage []*err.ErrorFormat     `json:"error_message,omitempty"`
	ErrorString  string                 `json:"error,omitempty"`
	Data         interface{}            `json:"data,omitempty"`
	Latency      string                 `json:"latency"`
	StatusCode   int                    `json:"-"`
	Error        error                  `json:"-"`
	Log          map[string]interface{} `json:"-"`
	Ctx          context.Context        `json:"-"`
	Success      *int                   `json:"success,omitempty"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{Ctx: context.Background(), Code: STATUSCODE_GENERICSUCCESS, StatusCode: GetHTTPCode(STATUSCODE_GENERICSUCCESS), Log: map[string]interface{}{}}
}

func NewJSONResponseWithCtx(ctx context.Context) *JSONResponse {
	return &JSONResponse{Ctx: ctx, Code: STATUSCODE_GENERICSUCCESS, StatusCode: GetHTTPCode(STATUSCODE_GENERICSUCCESS), Log: map[string]interface{}{}}
}

func (r *JSONResponse) SetSuccess(code int) *JSONResponse {
	r.Success = &code
	return r
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetLatency(latency float64) *JSONResponse {
	r.Latency = fmt.Sprintf("%.2f ms", latency)
	return r
}

func (r *JSONResponse) SetLog(key string, val interface{}) *JSONResponse {
	r.Log[key] = val
	return r
}

func getErrType(err error) error {
	switch err.(type) {
	case custerr.ErrChain:
		errType := err.(custerr.ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	return err
}

func (r *JSONResponse) SetError(err error, a ...string) *JSONResponse {
	err = getErrType(err)
	r.Error = err
	r.ErrorString = err.Error()
	r.Code = GetErrorCodeStr(err)
	r.StatusCode = GetHTTPCode(r.Code)
	if r.StatusCode == http.StatusInternalServerError {
		r.ErrorString = "Internal Server error"
	}
	if len(a) > 0 {
		r.ErrorString = a[0]
	}
	return r
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	b, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
