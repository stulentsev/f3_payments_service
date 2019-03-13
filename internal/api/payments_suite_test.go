package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	eventstore "github.com/looplab/eventhorizon/eventstore/mongodb"
	repo "github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/nvloff/f3_payments_service/gen/restapi"
	"github.com/nvloff/f3_payments_service/internal/api"
	"github.com/nvloff/f3_payments_service/internal/domain"
	"github.com/nvloff/f3_payments_service/internal/domain/mongodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testDomain      *domain.PaymentsDomain
	testRestHandler http.Handler
)

func TestF3PaymentsAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Payments API Suite")
}

/*
Convert JSON data into a map.
*/
func mapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result.(map[string]interface{})
}

func helperLoadBytes(name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	fileBytes, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	return fileBytes
}

func waitForPaymentProjection(id string) {
	Eventually(func() error {
		_, err := testDomain.PaymentRepo.Find(context.Background(), uuid.MustParse(id))
		return err
	}).Should(BeNil())
}

func waitForPaymentUpdate() {
	time.Sleep(100 * time.Millisecond)
}

func postPayment(body io.Reader) *httptest.ResponseRecorder {
	const postPaymentsPath = "/v1/payments"

	req := httptest.NewRequest(http.MethodPost, postPaymentsPath, body)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	testRestHandler.ServeHTTP(recorder, req)

	return recorder
}

func updatePayment(id string, body io.Reader) *httptest.ResponseRecorder {
	const basePath = "/v1/payments/"
	req := httptest.NewRequest(http.MethodPatch, basePath+id, body)

	recorder := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	testRestHandler.ServeHTTP(recorder, req)
	waitForPaymentUpdate()

	return recorder
}

func getPayment(id string) *httptest.ResponseRecorder {
	const basePath = "/v1/payments/"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, basePath+id, nil)
	testRestHandler.ServeHTTP(recorder, req)

	return recorder
}

func getPayments(query string) *httptest.ResponseRecorder {
	const basePath = "/v1/payments/"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, basePath+query, nil)
	testRestHandler.ServeHTTP(recorder, req)

	return recorder
}

func deletePayment(id string) *httptest.ResponseRecorder {
	const basePath = "/v1/payments/"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, basePath+id, nil)
	testRestHandler.ServeHTTP(recorder, req)

	return recorder
}

func createPayment(fixtureName string, id string) {
	reqBody := helperLoadBytes(fixtureName)
	recorder := postPayment(bytes.NewReader(reqBody))

	Expect(recorder.Code).To(Equal(201))

	waitForPaymentProjection(id)
}

func createFixturePayment() {
	fixtureFile := "create_payment.json"
	fixturePaymentId := "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"

	createPayment(fixtureFile, fixturePaymentId)
}

var _ = BeforeSuite(func() {
	url := os.Getenv("TEST_MONGO_HOST")

	testDomain = mongodb.BuildDomain(mongodb.Config{
		URL:        url,
		Collection: "payments_test",
	})

	p := api.New(api.Config{
		CommandBus:  testDomain.CommandBus,
		PaymentRepo: testDomain.PaymentRepo,
	})

	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: p,
	})

	Expect(err).NotTo(HaveOccurred())

	testRestHandler = handler

})

var _ = BeforeEach(func() {
	ctx := context.Background()

	testDomain.EventStore.(*eventstore.EventStore).Clear(ctx)
	testDomain.PaymentRepo.Parent().(*repo.Repo).Clear(ctx)
})
