package gopd

import (
	"testing"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestGetTemplateList(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setTestToken()

	file, _ := ioutil.ReadFile("./fixtures/templates/list.json")
	page := 10
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s?page=%d", TEMPLATE_API_ENDPOINT, page),
		func(req *http.Request) (*http.Response, error) {
			a := assert.New(t)
			a.Equal("Bearer test-token", req.Header.Get("Authorization"))
			a.Equal("application/json", req.Header.Get("Content-Type"))
			return httpmock.NewStringResponse(200, string(file[:])), nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s?page=1", TEMPLATE_API_ENDPOINT),
		func(req *http.Request) (*http.Response, error) {
			a := assert.New(t)
			a.Equal("Bearer test-token", req.Header.Get("Authorization"))
			a.Equal("application/json", req.Header.Get("Content-Type"))
			return httpmock.NewStringResponse(200, string(file[:])), nil
		})

	listData, err := GetTemplateList(page)
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal(12, listData.Count)
	a.Equal("https://api.pandadoc.com/public/v1/documents?page=2", listData.Next)
	a.Equal("", listData.Previous)
	a.Len(listData.Results, 1)
	result := listData.Results[0]

	a.Equal("UgNqHrtsGFqTSk8wtdzqPM", result.Id)
	a.Equal("Sample Template", result.Name)
	a.Equal("2014-10-06T08:42:13.836022Z", result.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", result.DateModified)

	listData, err = GetTemplateList()
	if err != nil {
		t.Error(err.Error())
	}

	a.Equal(12, listData.Count)
	a.Equal("https://api.pandadoc.com/public/v1/documents?page=2", listData.Next)
	a.Equal("", listData.Previous)
	a.Len(listData.Results, 1)
	result = listData.Results[0]

	a.Equal("UgNqHrtsGFqTSk8wtdzqPM", result.Id)
	a.Equal("Sample Template", result.Name)
	a.Equal("2014-10-06T08:42:13.836022Z", result.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", result.DateModified)
}

func TestGetTemplateDetails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setTestToken()

	file, _ := ioutil.ReadFile("./fixtures/templates/details.json")

	templateId := "UgNqHrtsGFqTSk8wtdzqPM"
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s/details", TEMPLATE_API_ENDPOINT, templateId),
		func(req *http.Request) (*http.Response, error) {
			a := assert.New(t)
			a.Equal("Bearer test-token", req.Header.Get("Authorization"))
			a.Equal("application/json", req.Header.Get("Content-Type"))
			return httpmock.NewStringResponse(200, string(file[:])), nil
		})

	details, err := GetTemplateDetails(templateId)
	if err != nil {
		t.Error(err.Error())
	}

	a := assert.New(t)
	a.Equal(templateId, details.Id)
	a.Equal("Sample Template", details.Name)
	a.Equal("2014-10-06T08:42:13.836022Z", details.DateCreated)
	a.Equal("2016-03-04T02:21:13.963750Z", details.DateModified)

	a.Equal("FyXaS4SlT2FY7uLPqKD9f2", details.CreatedBy.Id)
	a.Equal("john@appleseed.com", details.CreatedBy.Email)
	a.Equal("John", details.CreatedBy.FirstName)
	a.Equal("Appleseed", details.CreatedBy.LastName)
	a.Equal("https://pd-live-media.s3.amazonaws.com/users/FyXaS4SlT2FY7uLPqKD9f2/avatar.jpg", details.CreatedBy.Avatar)

	a.Len(details.Roles, 1)
	role := details.Roles[0]
	a.Equal("wHyMMrTeiY35HuLFr7CgAh", role.Id)
	a.Equal("Signer", role.Name)
	a.Equal(User{}, role.PreassignedPerson)

	a.Len(details.Metadata, 2)
	a.Equal("salesforce_opportunity_id", details.Metadata["key"])
	a.Equal("5003000000D8cuI", details.Metadata["value"])

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
	a.Equal("role", field.AssignedTo.Type)
	a.Equal("Signer", field.AssignedTo.Name)
	a.Equal(User{}, field.AssignedTo.PreassignedPerson)

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
	a.Equal("Stuffed Panda", item.Name)
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
}
