package gopd

import (
	"testing"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestWebHookData_FromRequest(t *testing.T) {
	file, _ := ioutil.ReadFile("./fixtures/webhooks/webhook.json")
	r, err := http.NewRequest("POST", "http://some-url/", bytes.NewBuffer(file))
	if err != nil {
		t.Error(err.Error())
	}

	whd := WebHookData{}
	whs, err := whd.FromRequest(r)
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Len(whs, 1)
	wh := whs[0]

	a.Equal("recipient_completed", wh.Event)
	a.Equal("msFYActMfJHqNTKH8YSvF1", wh.Data.Id)
	a.Equal("Sample Document", wh.Data.Name)
	a.Equal("2014-10-06T08:42:13.836022Z", wh.Data.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", wh.Data.DateModified)
	a.Equal("document.draft", wh.Data.Status)

	a.Equal("2016-09-02T22:26:52.227554", wh.Data.ActionDate)
	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", wh.Data.ActionBy.Id)
	a.Equal("john@appleseed.com", wh.Data.ActionBy.Email)
	a.Equal("John", wh.Data.ActionBy.FirstName)
	a.Equal("Appleseed", wh.Data.ActionBy.LastName)

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", wh.Data.CreatedBy.Id)
	a.Equal("john@appleseed.com", wh.Data.CreatedBy.Email)
	a.Equal("John", wh.Data.CreatedBy.FirstName)
	a.Equal("Appleseed", wh.Data.CreatedBy.LastName)
	a.Equal("https://pd-live-media.s3.amazonaws.com/users/FyXaS4SlT2FY7uLPqKD9f2/avatar.jpg", wh.Data.CreatedBy.Avatar)

	a.Len(wh.Data.Recipients, 1)
	recipient := wh.Data.Recipients[0]

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", recipient.Id)
	a.Equal("john@appleseed.com", recipient.Email)
	a.Equal("John", recipient.FirstName)
	a.Equal("Appleseed", recipient.LastName)
	a.Equal("Signer", recipient.RecipientType)
	a.Equal("signer", recipient.Role)
	a.Equal(true, recipient.HasCompleted)

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", wh.Data.SentBy.Id)
	a.Equal("john@appleseed.com", wh.Data.SentBy.Email)
	a.Equal("John", wh.Data.SentBy.FirstName)
	a.Equal("Appleseed", wh.Data.SentBy.LastName)
	a.Equal("https://pd-live-media.s3.amazonaws.com/users/FyXaS4SlT2FY7uLPqKD9f2/avatar.jpg", wh.Data.SentBy.Avatar)

	a.Len(wh.Data.Metadata, 2)
	a.Equal("123456", wh.Data.Metadata["salesforce_opp_id"])
	a.Equal("Panda", wh.Data.Metadata["my_favorite_pet"])

	a.Len(wh.Data.Tokens, 1)
	token := wh.Data.Tokens[0]
	a.Equal("Favorite Animal", token.Name)
	a.Equal("Panda", token.Value)

	a.Len(wh.Data.Fields, 1)
	field := wh.Data.Fields[0]
	a.Equal("YcLBNUKcx45UFxAK3NjLIH", field.UUID)
	a.Equal("Textfield", field.Name)
	a.Equal("Favorite Animal", field.Title)
	a.Equal("Panda", field.Value)
	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", field.AssignedTo.Id)
	a.Equal("john@appleseed.com", field.AssignedTo.Email)
	a.Equal("John", field.AssignedTo.FirstName)
	a.Equal("Appleseed", field.AssignedTo.LastName)
	a.Equal("Signer", field.AssignedTo.Role)
	a.Equal("signer", field.AssignedTo.RecipientType)
	a.Equal(true, field.AssignedTo.HasCompleted)
	a.Equal("recipient", field.AssignedTo.Type)

	tables := wh.Data.Pricing.Tables

	a.Len(tables, 1)
	table := wh.Data.Pricing.Tables[0]
	a.Equal(82307036, table.Id)
	a.Equal("PricingTable1", table.Name)
	a.Equal(true, table.IsIncludedInTotal)
	a.Equal(float32(10), table.Summary.Discount)
	a.Equal(float32(0), table.Summary.Tax)
	a.Equal(float32(60), table.Summary.Total)
	a.Equal(float32(60), table.Summary.Subtotal)
	a.Equal(float32(60), table.Total)

	a.Len(table.Items, 1)
	item := table.Items[0]
	a.Equal("4ElJ4FEsG4PHAVNPR5qoo9", item.Id)
	a.Equal(1, item.Qty)
	a.Equal("Toy Panda", item.Name)
	a.Equal("25", item.Cost)
	a.Equal("53", item.Price)
	a.Equal("Buy a Panda", item.Description)

	a.Len(item.CustomFields, 1)
	a.Equal("Sample Field", item.CustomFields["sampleField"])

	a.Len(item.CustomColumns, 1)
	a.Equal("Sample Column", item.CustomColumns["sampleColumn"])

	a.Equal(float32(10), item.Discount)
	a.Equal(float32(60), item.Subtotal)

	a.Len(wh.Data.Tags, 3)
	a.Contains(wh.Data.Tags, "test tag")
	a.Contains(wh.Data.Tags, "sales")
	a.Contains(wh.Data.Tags, "support")	
}

func TestWebHookData_FromRequest_InvalidData(t *testing.T) {
	file := []byte("some-invalid-json-data")
	r, err := http.NewRequest("POST", "http://some-url/", bytes.NewBuffer(file))
	if err != nil {
		t.Error(err.Error())
	}

	whd := WebHookData{}
	whs, err := whd.FromRequest(r)

	a := assert.New(t)
	a.NotNil(err)
	a.Nil(whs)
}
