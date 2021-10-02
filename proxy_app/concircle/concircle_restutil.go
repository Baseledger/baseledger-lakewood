package systemofrecord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"

	"github.com/spf13/viper"
	"github.com/unibrightio/proxy-api/logger"
)

type PostBusinessObjectDto struct {
	ID                 string          `json:"id"`
	Type               string          `json:"type"`
	OrganizationId     string          `json:"organization_id"`
	ObjectConnectionId string          `json:"object_connection_id"`
	MessageId          string          `json:"message_id"`
	TransactionId      string          `json:"transaction_id"`
	Payload            json.RawMessage `json:"payload"`
}

type PostStatusUpdateDto struct {
	Type               string        `json:"type"`
	MessageId          string        `json:"message_id"`
	Status             string        `json:"status"`
	Errors             []interface{} `json:"errors"`
	ObjectConnectionID string        `json:"object_connection_id"`
	OrganizationId     string        `json:"organization_id"`
	TransactionId      string        `json:"transaction_id"`
}

// not sure if this is ok approach, should store auth state in package and i keep refreshing it below if this is not defined
type CurrentAuthState struct {
	Token   string
	Cookies []*http.Cookie
}

var client http.Client
var currentAuthState CurrentAuthState

func InitClient() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client = http.Client{
		Jar: jar,
	}
}

func Auth() {
	authUrl := formatReqUrl("ubc/ubc/auth?sap-client=100")
	req, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		logger.Errorf("Concircle new request error %v\n", err.Error())
	}

	req.Header.Add("X-CSRF-Token", "Fetch")
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Concircle auth req error %v\n", err.Error())
	}

	if resp.StatusCode != http.StatusNoContent {
		logger.Errorf("Concircle auth req error, wrong status code %v\n", resp.StatusCode)
	}

	fmt.Printf("COOKIES %v\n", resp.Cookies())
	currentAuthState = CurrentAuthState{
		Token:   resp.Header.Get("X-CSRF-Token"),
		Cookies: resp.Cookies(),
	}
}

func PostBusinessObject(
	baseledgerObjectId string,
	businessObjectType string,
	recipientOrgId string,
	offchainMessageId string,
	baseledgerTransactionId string,
	businessObjectPayload string,
	trustmeshId string) {

	Auth()

	postBOUrl := formatReqUrl("ubc/ubc/business_objects")

	logger.Infof("businessObjectPayload: %s\n", businessObjectPayload)

	postBoDto := PostBusinessObjectDto{
		ID:                 baseledgerObjectId,
		Type:               businessObjectType,
		OrganizationId:     recipientOrgId,
		ObjectConnectionId: trustmeshId,
		TransactionId:      baseledgerTransactionId,
		MessageId:          offchainMessageId,
		Payload:            json.RawMessage(businessObjectPayload), // this avoids automatic escaptin of quotes by JSON marshaler
	}

	bodyBytes, _ := json.Marshal(postBoDto)
	reader := bytes.NewReader(bodyBytes)

	logger.Infof("POSTING NEW OBJECT URL: %s\n", postBOUrl)
	logger.Infof("POSTING NEW OBJECT BODY: %s\n", postBoDto)

	req, err := http.NewRequest("POST", postBOUrl, reader)
	if err != nil {
		logger.Errorf("Concircle new request error %v\n", err.Error())
	}

	req.Header.Add("X-CSRF-Token", currentAuthState.Token)
	for _, c := range currentAuthState.Cookies {
		req.AddCookie(c)
	}

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		logger.Errorf("DumpRequest error %s\n", err.Error())
	}

	logger.Infof(string(requestDump))

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Concircle post bo req error %v\n", err.Error())
	}

	logger.Infof("BO RESP %v\n", resp)
	logger.Infof("BO RESP BODY %v\n", resp.Body)
}

func PutStatusUpdate(
	baseledgerObjectId string,
	businessObjectType string,
	sorBusinessObjectId string,
	status string,
	baseledgerTransactionId string,
	organizationId string,
	trustmeshId string,
) {
	Auth()

	var resourceUrl = "ubc/ubc/business_objects/" + baseledgerObjectId + "/status"
	postStatusUpdateUrl := formatReqUrl(resourceUrl)

	postStatusUpdateDto := PostStatusUpdateDto{
		Type:               businessObjectType,
		MessageId:          sorBusinessObjectId,
		Status:             status,
		ObjectConnectionID: trustmeshId,
		OrganizationId:     organizationId,
		TransactionId:      baseledgerTransactionId,
	}

	bodyBytes, _ := json.Marshal(postStatusUpdateDto)
	reader := bytes.NewReader(bodyBytes)

	logger.Infof("POSTING STATUS UPDATE URL: %v\n", postStatusUpdateUrl)
	logger.Infof("POSTING STATUS UPDATE BODY: %v\n", postStatusUpdateDto)

	req, err := http.NewRequest("PUT", postStatusUpdateUrl, reader)
	if err != nil {
		logger.Errorf("Concircle new request error %v\n", err.Error())
	}

	req.Header.Add("X-CSRF-Token", currentAuthState.Token)
	for _, c := range currentAuthState.Cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Concircle post status update error %v\n", err.Error())
	}

	logger.Infof("STATUS UPDATE RESP %v\n", resp)
	logger.Infof("STATUS UPDATE RESP BODY %v\n", resp.Body)
}

func formatReqUrl(endpoint string) string {
	concircleUrl := viper.GetString("API_CONCIRCLE_URL") // s4h.rp.concircle.com
	concircleUser := viper.GetString("API_CONCIRCLE_USER")
	concirclePwd := viper.GetString("API_CONCIRCLE_PWD")

	return "https://" + concircleUser + ":" + concirclePwd + "@" + concircleUrl + "/" + endpoint
}

func (t *PostBusinessObjectDto) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
