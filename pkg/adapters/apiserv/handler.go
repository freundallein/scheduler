package apiserv

type Scheduler struct{}

type DummyParams struct {
	Data string `json:"data"`
}

// Dummy implements temporary rpc handler.
// curl -H "Content-Type: application/json" -X POST -d \
//  '{"jsonrpc": "2.0", "method": "Scheduler.Dummy", "params":[{"data":"testParams"}], "id": "1"}' \
//  http://127.0.0.1:8000/rpc/v0
func (handler *Scheduler) Dummy(params *DummyParams, result *map[string]string) error {
	*result = map[string]string{
		"data": params.Data,
	}
	return nil
}
