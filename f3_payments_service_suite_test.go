package f3_payments_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestF3PaymentsService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "F3PaymentsService Suite")
}
