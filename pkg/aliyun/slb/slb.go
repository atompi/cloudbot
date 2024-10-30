package slb

import (
	"encoding/json"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/atompi/cloudbot/pkg/cloudbot/options"
	"github.com/atompi/cloudbot/pkg/dataio"
	"github.com/atompi/cloudbot/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func createApiClient(opts options.AliyunOptions) (*openapi.Client, error) {
	config := utils.AliyunCreateClientConfig(
		tea.String(opts.AccessKeyId),
		tea.String(opts.AccessKeySecret),
		tea.String(opts.RegionId),
		tea.String(opts.Endpoint),
	)

	return openapi.NewClient(config)
}

func createApiInfo() *openapi.Params {
	return &openapi.Params{
		Action:      tea.String("DescribeLoadBalancers"),
		Version:     tea.String("2014-05-15"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}

func listSLB(c *openapi.Client, queries map[string]interface{}) ([]string, error) {
	params := createApiInfo()
	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, err := c.CallApi(params, request, runtime)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
		return []string{}, err
	}

	m, ok := resp["body"].(map[string]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", m)
		return []string{}, err
	}
	loadBalancers, ok := m["LoadBalancers"].(map[string]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", err)
		return []string{}, err
	}
	loadBalancer, ok := loadBalancers["LoadBalancer"].([]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", err)
		return []string{}, err
	}

	lbs := []string{}
	for _, v := range loadBalancer {
		lb, ok := v.(map[string]any)
		if !ok {
			zap.S().Errorf("convert interface to map failed: %v", err)
			return []string{}, err
		}
		lbs = append(lbs, lb["LoadBalancerId"].(string))
	}
	return lbs, nil
}

func describeSLB(c *openapi.Client, queries map[string]interface{}) (map[string]any, error) {
	params := createApiInfo()
	params.Action = tea.String("DescribeLoadBalancerAttribute")
	runtime := &teautil.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, err := c.CallApi(params, request, runtime)
	if err != nil {
		zap.S().Errorf("call api failed: %v", err)
	}

	m, ok := resp["body"].(map[string]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", err)
		return map[string]any{}, err
	}
	listenerPortsAndProtocol, ok := m["ListenerPortsAndProtocol"].(map[string]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", err)
		return map[string]any{}, err
	}
	listenerPortAndProtocol, ok := listenerPortsAndProtocol["ListenerPortAndProtocol"].([]any)
	if !ok {
		zap.S().Errorf("convert interface to map failed: %v", err)
		return map[string]any{}, err
	}

	s := map[string]any{
		"address":      m["Address"].(string),
		"address_type": m["AddressType"].(string),
		"bandwidth":    m["Bandwidth"].(json.Number),
		"status":       m["LoadBalancerStatus"].(string),
	}

	listeners := []map[string]any{}
	for _, v := range listenerPortAndProtocol {
		la, ok := v.(map[string]any)
		if !ok {
			zap.S().Errorf("convert interface to map failed: %v", err)
			return map[string]any{}, err
		}
		l := map[string]any{
			"ListenerPort":     la["ListenerPort"].(json.Number),
			"ListenerProtocol": la["ListenerProtocol"].(string),
		}
		listeners = append(listeners, l)
	}
	s["listeners"] = listeners

	return s, nil
}

func FetchSLB(t options.TaskOptions, queries map[string]interface{}) error {
	c, err := createApiClient(t.Aliyun)
	if err != nil {
		zap.S().Errorf("create api client failed: %v", err)
		return err
	}

	slbList, err := listSLB(c, queries)
	if err != nil {
		zap.S().Errorf("list slb failed: %v", err)
		return err
	}

	fields := []string{
		"address",
		"address_type",
		"status",
		"bandwidth",
		"listeners",
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, v := range slbList {
		q := map[string]interface{}{
			"LoadBalancerId": v,
		}

		wg.Add(1)
		ch <- 1
		go func() {
			s, err := describeSLB(c, q)
			if err != nil {
				zap.S().Errorf("describe slb failed: %v", err)
				return
			}
			jsonListeners, err := json.Marshal(s["listeners"])
			if err != nil {
				zap.S().Errorf("marshal listeners failed: %v", err)
				return
			}
			values := []interface{}{
				s["address"],
				s["address_type"],
				s["status"],
				s["bandwidth"],
				string(jsonListeners),
			}
			err = dataio.OutputMySQL(fields, values, t.Output)
			if err != nil {
				zap.S().Errorf("insert to db failed: %v", err)
				return
			}

			defer func() { wg.Done(); <-ch }()
		}()
	}

	wg.Wait()
	return nil
}
