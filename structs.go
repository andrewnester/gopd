package gopd

type Recipient struct {
	Id            string `json:"id,omitempty"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Role          string `json:"role,omitempty"`
	RecipientType string `json:"recipient_type,omitempty"`
	HasCompleted  bool `json:"has_completed,omitempty"`
	Type          string `json:"type,omitempty"`
}

type Token struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type FieldName struct {
	Value string `json:"value"`
}

type Field struct {
	FieldName FieldName `json:"field_name"`
}

type SectionRowOptions struct {
	Optional        bool `json:"optional"`
	MayEditQuantity bool `json:"may_edit_quantity"`
}

type SectionRowData struct {
	Qty         int `json:"qty"`
	Name        string `json:"name"`
	Cost        string `json:"cost"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Discount    int `json:"discount,omitempty"`
}

type SectionRow struct {
	Options       SectionRowOptions `json:"options"`
	Data          SectionRowData `json:"data"`
	CustomFields  map[string]string `json:"custom_fields,omitempty"`
	CustomColumns map[string]string `json:"custom_columns,omitempty"`
}

type Section struct {
	Title   string `json:"title"`
	Default bool `json:"default"`
	Rows    []SectionRow `json:"rows,omitempty"`
}

type PricingTable struct {
	Name     string `json:"name"`
	Sections []Section `json:"sections"`
}

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
}

type DocumentField struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Value      interface{} `json:"value"`
	AssignedTo Recipient `json:"assigned_to"`
}

type PricingTableSummary struct {
	Discount float32 `json:"discount,string"`
	Tax      float32 `json:"tax,string"`
	Total    float32 `json:"total,string"`
	Subtotal float32 `json:"subtotal,string"`
}

type PricingTableItem struct {
	Id            string `json:"id"`
	Qty           int `json:"qty"`
	Name          string `json:"name"`
	Cost          float32 `json:"cost,string"`
	Price         float32 `json:"price,string"`
	Description   string `json:"description"`
	CustomFields  map[string]string `json:"custom_fields"`
	CustomColumns map[string]string `json:"custom_columns"`
	Discount      float32 `json:"discount,string"`
	Subtotal      float32 `json:"subtotal,string"`
}

type DocumentPricingTable struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Total             float32 `json:"total,string"`
	IsIncludedInTotal bool `json:"is_included_in_total"`
	Summary           PricingTableSummary `json:"summary"`
	Items             []PricingTableItem `json:"items"`
}

type Pricing struct {
	Tables []DocumentPricingTable `json:"tables"`
}

type Document struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	Status       string `json:"status"`
	CreatedBy    User `json:"created_by"`
	Recipients   []Recipient `json:"recipients"`
	SentBy       User `json:"sent_by"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Tokens       []Token `json:"tokens,omitempty"`
	Fields       []DocumentField `json:"fields,omitempty"`
	Pricing      Pricing `json:"pricing,omitempty"`
	Tags         []string `json:"tags"`
}

type DocumentStatus struct {
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
	Status       string `json:"status"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type Role struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	PreassignedPerson User `json:"preassigned_person,omitempty"`
}

type TemplateAssignee struct {
	Id                string `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	PreassignedPerson User `json:"preassigned_person,omitempty"`
}

type TemplateField struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Value      interface{} `json:"value"`
	AssignedTo TemplateAssignee `json:"assigned_to"`
}

type Template struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	CreatedBy    User `json:"created_by"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Tokens       []Token `json:"tokens,omitempty"`
	Fields       []TemplateField `json:"fields,omitempty"`
	Pricing      Pricing `json:"pricing,omitempty"`
	Tags         []string `json:"tags"`
	Roles        []Role `json:"roles"`
}

// Short version of template data shown in template list
type TemplateShort struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type TemplateList struct {
	Count    int `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []TemplateShort `json:"results"`
}
