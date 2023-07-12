package client

import (
	"fmt"
	"reflect"
	"testing"
)

func invoicePartsToString(i *InvoiceParts) string {
	return fmt.Sprintf(`InvoiceParts{
	AmountMsat: %d,
	PaymentHash: %s,
	Description: %s,
	DescriptionHash: %v,
	Signature: %s,
}
`, i.AmountMsat, string(i.PaymentHash), string(i.Description), i.DescriptionHash, string(i.Signature),
	)
}

func TestParseInvoice(t *testing.T) {
	var invoice string
	var want, got *InvoiceParts
	var err error

	invoice = "lnbc1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdpl2pkx2ctnv5sxxmmwwd5kgetjypeh2ursdae8g6twvus8g6rfwvs8qun0dfjkxaq9qrsgq357wnc5r2ueh7ck6q93dj32dlqnls087fxdwk8qakdyafkq3yap9us6v52vjjsrvywa6rt52cm9r9zqt8r2t7mlcwspyetp5h2tztugp9lfyql"
	want = &InvoiceParts{
		AmountMsat:      0,
		PaymentHash:     []byte("qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypq"),
		Description:     []byte("2pkx2ctnv5sxxmmwwd5kgetjypeh2ursdae8g6twvus8g6rfwvs8qun0dfjkxaq"),
		DescriptionHash: false,
		Signature:       []byte("357wnc5r2ueh7ck6q93dj32dlqnls087fxdwk8qakdyafkq3yap9us6v52vjjsrvywa6rt52cm9r9zqt8r2t7mlcwspyetp5h2tztugp"),
	}
	got, err = ParseInvoice([]byte(invoice))
	if !reflect.DeepEqual(want, got) || err != nil {
		t.Fatalf("\nwanted: %s\ngot: %s\n", invoicePartsToString(want), invoicePartsToString(got))
	}

	invoice = "lnbc2500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh"
	want = &InvoiceParts{
		AmountMsat:      250000000,
		PaymentHash:     []byte("qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypq"),
		Description:     []byte("xysxxatsyp3k7enxv4js"),
		DescriptionHash: false,
		Signature:       []byte("uk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgp"),
	}
	got, err = ParseInvoice([]byte(invoice))
	if !reflect.DeepEqual(want, got) || err != nil {
		t.Fatalf("\nwanted: %s\ngot: %s\n", invoicePartsToString(want), invoicePartsToString(got))
	}

	invoice = "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"
	want = &InvoiceParts{
		AmountMsat:      1500000,
		PaymentHash:     []byte("jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3s"),
		Description:     []byte("vfhkcap3xyhx7un8"),
		DescriptionHash: false,
		Signature:       []byte("y4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq"),
	}
	got, err = ParseInvoice([]byte(invoice))
	if !reflect.DeepEqual(want, got) || err != nil {
		t.Fatalf("\nwanted: %s\ngot: %s\n", invoicePartsToString(want), invoicePartsToString(got))
	}
}
