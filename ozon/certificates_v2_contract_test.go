package ozon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProductCertificateV2Contracts(t *testing.T) {
	t.Parallel()

	t.Run("options", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/product/certification/options", "", `{"option":[{"name":"NAME","required":true}]}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Certificates().GetProductCertificateOptionsV2(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if len(response.Options) != 1 || response.Options[0].Name != "NAME" || !response.Options[0].Required {
			t.Fatalf("response = %+v", response)
		}
	})

	params := ProductCertificateV2Params{
		AccordanceType:     "EAEU",
		CertificateCountry: "RU",
		CertificateType:    "CERTIFICATE_OF_CONFORMITY",
		ExpiredDate:        &ProductCertificateV2ExpiredDate{Date: &ProductCertificateV2Date{Day: 1, Month: 9, Year: 2027}},
		Files:              []ProductCertificateV2File{{FileContent: "YmFzZTY0", Name: "certificate.pdf"}},
		IssueDate:          certificateIssueDate(t),
		LinkToRegistry:     "https://registry.example/certificate",
		Name:               "Certificate",
		Number:             "CERT-1",
		ProductType:        "MEDICAL_PRODUCT",
		SKUs:               []string{"123"},
	}
	wantBody := `{"params":{"accordance_type":"EAEU","certificate_country":"RU","certificate_type":"CERTIFICATE_OF_CONFORMITY","expired_date":{"date":{"day":1,"month":9,"year":2027}},"files":[{"file_content":"YmFzZTY0","name":"certificate.pdf"}],"issue_date":"2026-09-01T00:00:00Z","link_to_registry":"https://registry.example/certificate","name":"Certificate","number":"CERT-1","product_type":"MEDICAL_PRODUCT","skus":["123"]}}`

	t.Run("params", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/product/certification/params", wantBody, `{"params":[{"name":"FILES","required":true}]}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Certificates().GetProductCertificateParamsV2(context.Background(), &GetProductCertificateParamsV2Request{Params: params})
		if err != nil {
			t.Fatal(err)
		}
		if len(response.Params) != 1 || response.Params[0].Name != "FILES" || !response.Params[0].Required {
			t.Fatalf("response = %+v", response)
		}
	})

	t.Run("create", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(requestContractHandler(t, http.MethodPost, "/v2/product/certificate/create", wantBody, `{"certificate_id":42,"params":[{"error":"","name":"FILES","state":"VALID"}],"status":"COMPLETED"}`))
		defer server.Close()
		response, err := NewClient(WithURI(server.URL)).Certificates().CreateProductCertificateV2(context.Background(), &CreateProductCertificateV2Request{Params: params})
		if err != nil {
			t.Fatal(err)
		}
		if response.CertificateID == nil || *response.CertificateID != 42 || response.Status != "COMPLETED" || len(response.Params) != 1 || response.Params[0].State != "VALID" {
			t.Fatalf("response = %+v", response)
		}
	})
}

func certificateIssueDate(t *testing.T) *time.Time {
	t.Helper()
	date := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	return &date
}
