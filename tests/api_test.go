package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const ApiURL = "http://localhost:3000"

func TestApi(t *testing.T) {
	testcase := getTestCases()
	ctx := context.Background()
	client := &http.Client{}

	for _, tc := range testcase {
		t.Run(t.Name(), func(t *testing.T) {
			for idx := range tc.Steps {
				step := &tc.Steps[idx]
				request, err := step.Request(t, ctx, &tc)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Accept", "application/json")
				require.NoError(t, err)

				// Send request
				response, err := client.Do(request)

				require.NoError(t, err)
				defer response.Body.Close()

				// Check response
				ReadJsonResult(t, response, step)
				step.Expect(t, ctx, &tc, response, step.Result)
			}
		})
	}
}

func getTestCases() []TestCase {
	id_customer := ""
	id_investment := ""

	return []TestCase{
		{
			Name: "Test to register new customer",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req := map[string]string{
							"name": "user_" + uuid.New().String()[:8],
						}

						body, err := json.Marshal(req)
						require.NoError(t, err)

						return http.NewRequest("POST", ApiURL+"/api/customers", bytes.NewReader(body))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, r *http.Response, m map[string]any) {
						require.Equal(t, http.StatusCreated, r.StatusCode)
						RequireIsUUID(t, m["id"].(string))
						id_customer = m["id"].(string)
					},
				},
			},
		},
		{
			Name: "Test to add product investments",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req := map[string]string{
							"name": "asset_" + uuid.New().String()[:8],
						}

						body, err := json.Marshal(req)
						require.NoError(t, err)

						return http.NewRequest("POST", ApiURL+"/api/investments", bytes.NewReader(body))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, r *http.Response, m map[string]any) {
						require.Equal(t, http.StatusCreated, r.StatusCode)
						id_investment = m["id"].(string)
					},
				},
			},
		},
		{
			Name: "Test to deposit",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req := map[string]interface{}{
							"customer_id":   id_customer,
							"investment_id": id_investment,
							"amount":        100000.00,
						}

						body, err := json.Marshal(req)
						require.NoError(t, err)

						return http.NewRequest("POST", ApiURL+"/api/transactions/deposit", bytes.NewReader(body))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, r *http.Response, m map[string]any) {
						t.Log("POST /api/transactions/deposit")
						t.Log(id_customer)
						t.Log(id_investment)
						require.Equal(t, http.StatusOK, r.StatusCode)
						RequireIsUUID(t, m["transaction_id"].(string))
					},
				},
			},
		},
		{
			Name: "Test to withdraw",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req := map[string]interface{}{
							"customer_id":   id_customer,
							"investment_id": id_investment,
							"amount":        50000.00,
						}

						body, err := json.Marshal(req)
						require.NoError(t, err)

						return http.NewRequest("POST", ApiURL+"/api/transactions/withdraw", bytes.NewReader(body))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, r *http.Response, m map[string]any) {
						t.Log("POST /api/transactions/withdraw")
						t.Log(id_customer)
						t.Log(id_investment)
						require.Equal(t, http.StatusOK, r.StatusCode)
						RequireIsUUID(t, m["transaction_id"].(string))
					},
				},
			},
		},
		{
			Name: "Test to get portfolio",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("GET", ApiURL+"/api/portfolio/"+id_customer+"/"+id_investment, nil)
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, r *http.Response, m map[string]any) {
						t.Log("POST /api/portfolio/" + id_customer + "/" + id_investment)
						t.Log(id_customer)
						t.Log(id_investment)
						t.Log(m)
						require.Equal(t, http.StatusOK, r.StatusCode)
					},
				},
			},
		},
	}
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

type TestCaseStep struct {
	Request RequestFunc
	Expect  ExpectFunc
	Result  map[string]any
}

func ResponseContains(t *testing.T, resp *http.Response, text string) {
	body, err := io.ReadAll(resp.Body)
	bodyStr := string(body)
	require.NoError(t, err)
	require.Contains(t, bodyStr, text)
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	step.Result = result
	require.NoError(t, err)
}

func RequireIsUUID(t *testing.T, value string) {
	_, err := uuid.Parse(value)
	require.NoError(t, err)
}
