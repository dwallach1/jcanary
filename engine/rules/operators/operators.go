package operators

import (
	"encoding/json"
	"fmt"
	"io"
	"jcanary/interoperator"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/mitchellh/mapstructure"
)

type OperatorType string

const (
	WebRequestOperator OperatorType = "webrequest"
	EqualsOperatorType OperatorType = "equals"
)

type Result struct {
	Container *gabs.Container
	Err       error
}

func (r Result) HasError() bool {
	return r.Err != nil
}

type OperatorT interface {
	HttpRequestOperator | EqualsOperator | ArrayGetOperator
}

type Operator interface {
	Operate(interoperator.VariableBag, *[]*Result) *Result
}

func New(t OperatorType, operatorConfig map[string]interface{}) (Operator, error) {
	switch t {
	case WebRequestOperator:
		var op HttpRequestOperator
		err := mapstructure.Decode(operatorConfig, &op)
		if err != nil {
			fmt.Printf("unable to initalize webrequest operator: %v", err)
			return &NilOperator{}, err
		}
		return &op, nil
	case EqualsOperatorType:
		var op EqualsOperator
		err := mapstructure.Decode(operatorConfig, &op)
		if err != nil {
			fmt.Printf("unable to initalize equals operator: %v", err)
			return &NilOperator{}, err
		}
		return &op, nil
	default:
	}
	return &NilOperator{}, fmt.Errorf("no matching operator found")
}

type NilOperator struct{}

func (o *NilOperator) Operate(varBag interoperator.VariableBag, pipeline *[]*Result) *Result {
	return nil
}

// HttpRequestOperator
type HttpRequestOperator struct {
	Type       OperatorType `json:"type"`
	Connection struct {
		Url     string `json:"url"`
		Method  string `json:"method"`
		Body    string `json:"body"`
		Headers []struct {
			Key string `json:"key"`
			Val string `json:"val"`
		} `json:"headers"`
		QueryParams []struct {
			Key string `json:"key"`
			Val string `json:"val"`
		} `json:"queryParams"`
	} `json:"connection"`
}

func (o *HttpRequestOperator) Operate(varBag interoperator.VariableBag, pipeline *[]*Result) *Result {
	var result Result
	httpclient := &http.Client{
		Timeout: time.Second * 30,
	}
	var body io.Reader
	if o.Connection.Body != "" {
		body = strings.NewReader(interoperator.BuildString(o.Connection.Body, varBag))
	}
	url := interoperator.BuildString(o.Connection.Url, varBag)
	method := strings.ToUpper(o.Connection.Method)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		result.Err = fmt.Errorf("unable to create request: %w", err)
		return &result
	}
	for _, header := range o.Connection.Headers {
		req.Header.Add(header.Key, interoperator.BuildString(header.Val, varBag))
	}

	q := req.URL.Query()
	for _, qparam := range o.Connection.QueryParams {
		q.Add(qparam.Key, interoperator.BuildString(qparam.Val, varBag))
	}
	req.URL.RawQuery = q.Encode()

	fmt.Printf("invoking url: %v\n", url)

	resp, err := httpclient.Do(req)
	if err != nil {
		result.Err = fmt.Errorf("unable to do http request: %w", err)
		return &result
	}
	m := make(map[string]interface{})
	var responseBody interface{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		result.Err = fmt.Errorf("unable to unmarshal response body: %w", err)
		return &result
	}
	m["statusCode"] = resp.StatusCode
	m["responseBody"] = responseBody
	fmt.Printf(" -> HTTP Operation: %v::%v < statuscode: %v > \n", method, url, resp.StatusCode)
	if resp.StatusCode > 399 {
		fmt.Printf("body: %v\n", responseBody)
	}
	container, err := gabs.Consume(m)
	if err != nil {
		result.Err = fmt.Errorf("failed to consume container: %w", err)
		return &result
	}
	result.Container = container
	return &result
}

type Operand struct {
	Type string      `json:"type"`
	Val  interface{} `json:"val"`
}

type EqualsOperator struct {
	Type         OperatorType `json:"type"`
	StepRef      int          `json:"stepRef"`
	LeftOperand  Operand      `json:"leftOperand"`
	RightOperand Operand      `json:"rightOperand"`
}

func (o *EqualsOperator) Operate(varBag interoperator.VariableBag, pipeline *[]*Result) *Result {
	var result Result

	right := o.RightOperand.Val
	left := o.LeftOperand.Val
	if o.RightOperand.Type == "reference" {
		path := right.(string)
		right = (*pipeline)[o.StepRef].Container.Path(path).Data()
	}
	if o.LeftOperand.Type == "reference" {
		path := left.(string)
		left = (*pipeline)[o.StepRef].Container.Path(path).Data()
	}
	res := reflect.DeepEqual(right, left)
	m := make(map[string]interface{})
	m["result"] = res
	container, err := gabs.Consume(m)
	if err != nil {
		result.Err = fmt.Errorf("failed to consume container: %w", err)
		return &result
	}
	fmt.Printf(" -> Equals Operation: < %v > == < %v > resulted in %v\n", left, right, res)
	result.Container = container
	return &result
}

type ArrayGetOperator struct{}
