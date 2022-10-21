package render

import (
	"encoding/json"
	"io"
	"net/http"
	"tie.prodigy9.co/config"
)

func Text(resp http.ResponseWriter, r *http.Request, text string) {
	resp.Header().Set("Content-Type", "text/plain")
	resp.WriteHeader(200)
	if _, err := resp.Write([]byte(text)); err != nil {
		Error(resp, r, http.StatusInternalServerError, err)
	}
}

func JSON(resp http.ResponseWriter, r *http.Request, obj interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
	if err := json.NewEncoder(resp).Encode(obj); err != nil {
		Error(resp, r, http.StatusInternalServerError, err)
	}
}

// TODO: status code should be specified by the code originating the error since otherwise
//   we'll have to switch in controllers comparing all possible error cases which is not
//   ideal.
func Error(resp http.ResponseWriter, r *http.Request, status int, err error) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)

	errObj := decorateError(err)
	if err_ := json.NewEncoder(resp).Encode(errObj); err_ != nil {
		config.FromRequest(r).Printf("%s %s %s - %s\n",
			r.RemoteAddr, r.Method, r.RequestURI, err_.Error())
	}
}

func Download(resp http.ResponseWriter, r *http.Request, filename string, reader io.Reader) {
	resp.Header().Set("Content-Description", "File Transfer")
	resp.Header().Set("Content-Transfer-Encoding", "binary")
	resp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.WriteHeader(200)
	if _, err := io.Copy(resp, reader); err != nil {
		Error(resp, r, http.StatusInternalServerError, err)
	}
}

func decorateError(err error) interface{} {
	errObj := &struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "unknown",
		Message: err.Error(),
	}

	if code, ok := err.(interface{ Code() string }); ok {
		errObj.Code = code.Code()
	}

	return errObj
}
