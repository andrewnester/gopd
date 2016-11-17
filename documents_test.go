package gopd

import (
	"testing"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"bytes"
)

func TestFromTemplateDocument_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	create, _ := ioutil.ReadFile("./fixtures/documents/create.json")
	created, _ := ioutil.ReadFile("./fixtures/documents/created.json")
	var td FromTemplateDocument

	httpmock.RegisterResponder("POST", DOCUMENT_API_ENDPOINT,
		func(req *http.Request) (*http.Response, error) {
			var doc FromTemplateDocument
			if err := json.NewDecoder(req.Body).Decode(&doc); err != nil {
				t.Error("Wrong request body")
			}

			a := assert.New(t)
			a.Equal(td, doc)

			// Going to compare JSON string request with expected JSON string
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, create); err != nil {
				t.Error(err)
			}
			jsonStr, _ := json.Marshal(doc)
			a.Equal(string(buffer.Bytes()[:]), string(jsonStr[:]))

			resp := httpmock.NewBytesResponse(201, created)
			return resp, nil
		},
	)

	td = FromTemplateDocument{
		Name: "From API",
		TemplateUuid: "nkpMcUE75tBo6xdBbPXUPJ",
		Recipients: []Recipient{
			{Email:"andrew.nester.dev@gmail.com", FirstName:"Andrew", LastName:"Nester", Role:"Signer"},
			{Email:"andrew.nester.dev2@gmail.com", FirstName:"Andrew2", LastName:"Nester2"},
		},
		Tokens: []Token{
			{Name:"token1", Value:"token1"},
			{Name:"token12", Value:"token2"},
		},
		Fields: []Field{
			{FieldName:FieldName{"fieldName1"}},
			{FieldName:FieldName{"fieldName2"}},
		},
		Metadata: map[string]string{
			"some_field":"some_value",
			"some_field2":"some_value2",
		},
		PricingTables:[]PricingTable{
			{Name:"pricing1", Sections:[]Section{
				{Title:"section1", Default:true},
				{Title:"section2", Default:false},
			}},
		},
	}

	resp, err := td.Create()
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal(td.Name, resp.Name)
	a.Equal("msFYActMfJHqNTKH8YSvF1", resp.Id)
	a.Equal("document.uploaded", resp.Status)
	a.Equal("2014-10-06T08:42:13.836022Z", resp.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", resp.DateModified)
}

func TestGetDocumentDetails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	file, _ := ioutil.ReadFile("./fixtures/documents/details.json")

	docId := "msFYActMfJHqNTKH8YSvF1"
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s/details", DOCUMENT_API_ENDPOINT, docId),
		httpmock.NewStringResponder(200, string(file[:])))

	details, err := GetDocumentDetails(docId)
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal(docId, details.Id)
	a.Equal("Sample Document", details.Name)
	a.Equal("2014-10-06T08:42:13.836022Z", details.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", details.DateModified)
	a.Equal("document.draft", details.Status)

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", details.CreatedBy.Id)
	a.Equal("john@appleseed.com", details.CreatedBy.Email)
	a.Equal("John", details.CreatedBy.FirstName)
	a.Equal("Appleseed", details.CreatedBy.LastName)
	a.Equal("https://pd-live-media.s3.amazonaws.com/users/FyXaS4SlT2FY7uLPqKD9f2/avatar.jpg", details.CreatedBy.Avatar)

	a.Len(details.Recipients, 1)
	recipient := details.Recipients[0]

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", recipient.Id)
	a.Equal("john@appleseed.com", recipient.Email)
	a.Equal("John", recipient.FirstName)
	a.Equal("Appleseed", recipient.LastName)
	a.Equal("Signer", recipient.RecipientType)
	a.Equal("signer", recipient.Role)
	a.Equal(true, recipient.HasCompleted)

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", details.SentBy.Id)
	a.Equal("john@appleseed.com", details.SentBy.Email)
	a.Equal("John", details.SentBy.FirstName)
	a.Equal("Appleseed", details.SentBy.LastName)
	a.Equal("https://pd-live-media.s3.amazonaws.com/users/FyXaS4SlT2FY7uLPqKD9f2/avatar.jpg", details.SentBy.Avatar)

	a.Len(details.Metadata, 2)
	a.Equal("123456", details.Metadata["salesforce_opp_id"])
	a.Equal("Panda", details.Metadata["my_favorite_pet"])

	a.Len(details.Tokens, 1)
	token := details.Tokens[0]
	a.Equal("Favorite Animal", token.Name)
	a.Equal("Panda", token.Value)

	a.Len(details.Fields, 1)
	field := details.Fields[0]
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

	tables := details.Pricing.Tables

	a.Len(tables, 1)
	table := details.Pricing.Tables[0]
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

	a.Len(details.Tags, 3)
	a.Contains(details.Tags, "test tag")
	a.Contains(details.Tags, "sales")
	a.Contains(details.Tags, "support")

	/**

	 */
}
