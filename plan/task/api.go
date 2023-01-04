package task

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/parnurzeal/gorequest"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("API", func() Task { return &apiCallTask{} })
}

type apiCallTask struct {
}

func (t *apiCallTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	R, err := buildRequest(v)
	if err != nil {
		return nil, err
	}
	resp, err := makeRequest(R)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var bodyVal interface{}

	var isString, isBytes bool
	r := v.LookupPath(cue.ParsePath("resp.body"))
	if r.Exists() {
		if r.IncompleteKind() == cue.StringKind {
			isString = true
		}
		if r.IncompleteKind() == cue.BytesKind {
			isBytes = true
		}
	}

	if isString {
		bodyVal = string(body)
	} else if isBytes {
		bodyVal = body
	} else {
		bodyVal = v.Context().CompileBytes(body, cue.Filename("body"))
	}

	Then(ctx, v)
	value := compiler.NewValue()
	respV := map[string]interface{}{
		"status":     resp.Status,
		"statusCode": resp.StatusCode,
		"body":       bodyVal,
		"header":     resp.Header,
		"trailer":    resp.Trailer,
	}
	output := value.FillPath(cue.ParsePath("resp"), respV)
	output = output.FillPath(cue.ParsePath("output"), R.Url)
	return &output, nil
}

func (t *apiCallTask) Name() string {
	return "API"
}

func buildRequest(val *cue.Value) (R *gorequest.SuperAgent, err error) {
	req := val.Eval()
	R = gorequest.New()

	// this should be the default after unifying, and this should always exist
	// R.Method = "GET"
	method := req.LookupPath(cue.ParsePath("method"))
	// so we may not need this is not needed, but it is defensive
	if method.Exists() {
		R.Method, err = method.String()
		if err != nil {
			return R, err
		}
	}

	host := req.LookupPath(cue.ParsePath("host"))
	hostStr, err := host.String()
	if err != nil {
		return
	}

	path := req.LookupPath(cue.ParsePath("path"))
	pathStr, err := path.String()
	if err != nil {
		return
	}
	R.Url = hostStr + pathStr

	headers := req.LookupPath(cue.ParsePath("headers"))
	if headers.Exists() {
		H, err := headers.Struct()
		if err != nil {
			return R, err
		}
		hIter := H.Fields()
		for hIter.Next() {
			label := hIter.Label()
			value, err := hIter.Value().String()
			if err != nil {
				return R, err
			}
			R.Header.Add(label, value)
		}
	}

	query := req.LookupPath(cue.ParsePath("query"))
	if query.Exists() {
		Q, err := query.Struct()
		if err != nil {
			return R, err
		}
		qIter := Q.Fields()
		for qIter.Next() {
			label := qIter.Label()
			value, err := qIter.Value().String()
			if err != nil {
				return R, err
			}
			R.QueryData.Add(label, value)
		}
	}

	form := req.LookupPath(cue.ParsePath("form"))
	if form.Exists() {
		F, err := form.Struct()
		if err != nil {
			return R, err
		}
		fIter := F.Fields()
		for fIter.Next() {
			label := fIter.Label()
			value, err := fIter.Value().String()
			if err != nil {
				return R, err
			}
			R.FormData.Add(label, value)
		}
	}

	data := req.LookupPath(cue.ParsePath("data"))
	if data.Exists() {
		err := data.Decode(&R.Data)
		if err != nil {
			return R, err
		}
	}

	timeout := req.LookupPath(cue.ParsePath("timeout"))
	if timeout.Exists() {
		to, err := timeout.String()
		if err != nil {
			return R, err
		}
		d, err := time.ParseDuration(to)
		if err != nil {
			return R, err
		}
		R.Timeout(d)
	}

	retry := req.LookupPath(cue.ParsePath("retry"))
	if retry.Exists() {
		C := 3
		count := retry.LookupPath(cue.ParsePath("count"))
		if count.Exists() {
			c, err := count.Int64()
			if err != nil {
				return R, err
			}
			C = int(c)
		}

		D := time.Second * 6
		timer := retry.LookupPath(cue.ParsePath("timer"))
		if timer.Exists() {
			t, err := timer.String()
			if err != nil {
				return R, err
			}
			d, err := time.ParseDuration(t)
			if err != nil {
				return R, err
			}
			D = d
		}

		CS := []int{}
		codes := retry.LookupPath(cue.ParsePath("codes"))
		if codes.Exists() {
			L, err := codes.List()
			if err != nil {
				return R, err
			}
			for L.Next() {
				v, err := L.Value().Int64()
				if err != nil {
					return R, err
				}
				CS = append(CS, int(v))
			}
		}

		R.Retry(C, D, CS...)
	}

	// Todo, add 'print curl' option
	// check if set, then fill on return
	//curl, err := R.AsCurlCommand()
	//if err != nil {
	//return nil, err
	//}
	//fmt.Println(curl)

	return
}

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func makeRequest(R *gorequest.SuperAgent) (gorequest.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in HTTP: %v %v\n", R, r)
		}
	}()

	resp, body, errs := R.End()

	if len(errs) != 0 && resp == nil {
		return resp, fmt.Errorf("%v", errs)
	}

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		return resp, fmt.Errorf("internal weirdr error:\b%v\n%s", errs, body)
	}
	if len(errs) != 0 {
		return resp, fmt.Errorf("internal error:\n%v\n%s", errs, body)
	}

	return resp, nil
}
