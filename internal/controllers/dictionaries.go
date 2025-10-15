package controllers

import (
	kafka "bff/pkg/kfk"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/constants"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/customErrors"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	partnercore "gitlab.HIDDEN.com/HIDDEN/sms-core/partner"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/types"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/utils"
	"strconv"
	"time"
)

func (rc *RestController) ContractTypes(c *gin.Context) {
	contractTypes := []string{"Vendor", "Client", "Bilateral"}
	utils.SendOkSliceResponse(c, nil, &contractTypes)
}

func (rc *RestController) PartnerOptions(c *gin.Context) {
	value := partnercore.GetPartnerOptionsPayload{
		Skip:         c.Query("skip"),
		Limit:        c.Query("limit"),
		Search:       c.Query("search"),
		ContractType: c.Query("contractType"),
	}

	err := value.Validate()

	if err != nil {
		logger.Error().Err(err).Msg("value.Validate")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetOptions, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)

	if err != nil {
		logger.Error().Err(err).Msgf("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msgf("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, h, err := utils.ListenReply[[]types.DictionaryOption](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msgf("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) PartnerProductOptions(c *gin.Context) {
	value := partnercore.GetPartnerProductOptionsPayload{
		Skip:         c.Query("skip"),
		Search:       c.Query("search"),
		PartnerId:    c.Param("partnerId"),
		Limit:        c.Query("limit"),
		ContractType: c.Query("contractType"),
	}

	err := value.Validate()

	if err != nil {
		logger.Error().Err(err).Msg("value.Validate")
		utils.SendErrorResponse(c, customErrors.ErrRequestUnprocessableEntityError)
		return
	}

	headers := utils.FormKafkaHeaders(c, constants.Products, constants.ActionGetOptions, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, h, err := utils.ListenReply[[]types.DictionaryOption](chann, rc.avro, time.Second*5)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) AccountManagers(c *gin.Context) {
	value := partnercore.GetPartnerByPayload{
		Skip:     c.Query("skip"),
		Limit:    c.Query("limit"),
		Search:   c.Query("search"),
		SearchBy: []string{"account_manager"},
		Fields:   []string{"_id", "account_manager"},
	}

	headers := utils.FormKafkaHeaders(c, constants.Partners, constants.ActionGetAccountManagers, kafka.Settings.ReplyTopic)

	kfkMsg, err := rc.avro.Encode(value)
	if err != nil {
		logger.Error().Err(err).Msg("rc.avro.Encode")
		utils.SendErrorResponse(c, err)
		return
	}

	chann, err := rc.adapter.ProduceMessage(constants.Request, kfkMsg, headers)

	if err != nil {
		logger.Error().Err(err).Msg("rc.adapter.ProduceMessage")
		utils.SendErrorResponse(c, err)
		return
	}

	d, h, err := utils.ListenReply[[]string](chann, rc.avro, time.Second*30)

	if err != nil {
		logger.Error().Err(err).Msg("utils.ListenReply")
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendOkSliceResponse(c, h, d)
}

func (rc *RestController) BusinessUnits(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"Transit", "Enterprise"})
	return
}

func (rc *RestController) SubUnits(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"Transit", "Channels", "PRS"})
	return
}

func (rc *RestController) UnitsMap(c *gin.Context) {
	utils.SendOkResponse(c, nil, map[string][]string{
		"Transit":    {"PRS"},
		"Enterprise": {"Channels"},
	})
	return
}

func (rc *RestController) CountriesList(c *gin.Context) {
	countries := []string{
		"Abkhazia",
		"Afghanistan",
		"Albania",
		"Algeria",
		"American Samoa",
		"Andorra",
		"Angola",
		"Anguilla",
		"Antigua and Barbuda",
		"Argentina",
		"Armenia",
		"Aruba",
		"Australia",
		"Austria",
		"Azerbaijan",
		"Bahamas",
		"Bahrain",
		"Bangladesh",
		"Barbados",
		"Belarus",
		"Belgium",
		"Belize",
		"Benin",
		"Bermuda",
		"Bhutan",
		"Bolivia",
		"Bosnia and Herzegovina",
		"Botswana",
		"Brazil",
		"British Virgin Islands",
		"Brunei",
		"Bulgaria",
		"Burkina Faso",
		"Burundi",
		"Cambodia",
		"Cameroon",
		"Canada",
		"Cape Verde Islands",
		"Cayman Islands",
		"Central African Republic",
		"Chad",
		"Chile",
		"China",
		"Colombia",
		"Comoros",
		"Congo",
		"Cook Islands",
		"Costa Rica",
		"Croatia",
		"Cuba",
		"Cyprus",
		"Czech Republic",
		"DR of Congo",
		"Denmark",
		"Djibouti",
		"Dominica",
		"Dominican Republic",
		"Ecuador",
		"Egypt",
		"El Salvador",
		"Equatorial Guinea",
		"Eritrea",
		"Estonia",
		"Ethiopia",
		"Falkland Islands",
		"Faroe Islands",
		"Fiji",
		"Finland",
		"France",
		"French Antilles",
		"French Guiana",
		"French Polynesia",
		"Gabon",
		"Gambia",
		"Georgia",
		"Germany",
		"Ghana",
		"Gibraltar",
		"Greece",
		"Greenland",
		"Grenada",
		"Guadeloupe",
		"Guam",
		"Guatemala",
		"Guinea",
		"Guinea-Bissau",
		"Guyana",
		"Haiti",
		"Honduras",
		"Hong Kong",
		"Hungary",
		"Iceland",
		"India",
		"Indonesia",
		"InvalidNumbersHLR",
		"Iran",
		"Iraq",
		"Ireland",
		"Israel",
		"Italy",
		"Ivory Coast",
		"Jamaica",
		"Japan",
		"Jordan",
		"Kazakhstan",
		"Kenya",
		"Kiribati",
		"Kosovo",
		"Kuwait",
		"Kyrgyzstan",
		"Laos",
		"Latvia",
		"Lebanon",
		"Lesotho",
		"Liberia",
		"Libya",
		"Liechtenstein",
		"Lithuania",
		"Luxembourg",
		"Macau",
		"Macedonia",
		"Madagascar",
		"Malawi",
		"Malaysia",
		"Maldives",
		"Mali",
		"Malta",
		"Marshall Islands",
		"Martinique",
		"Mauritania",
		"Mauritius",
		"Mexico",
		"Micronesia",
		"Moldova",
		"Monaco",
		"Mongolia",
		"Montenegro",
		"Montserrat",
		"Morocco",
		"Mozambique",
		"Myanmar",
		"NOT FOR COMMERCIAL USE",
		"Namibia",
		"Nauru",
		"Nepal",
		"Netherlands",
		"Netherlands Antilles",
		"New Caledonia",
		"New Zealand",
		"Nicaragua",
		"Niger",
		"Nigeria",
		"Niue",
		"North America",
		"North Korea",
		"Northern Mariana Islands",
		"Norway",
		"Oman",
		"Pakistan",
		"Palau",
		"Palestine",
		"Panama",
		"Papua New Guinea",
		"Paraguay",
		"Peru",
		"Philippines",
		"Poland",
		"Portugal",
		"Puerto Rico",
		"Qatar",
		"Rest of the world",
		"Reunion Island",
		"Romania",
		"Russia",
		"Rwanda",
		"Saint Helena, Ascension and Tristan da Cunha",
		"Saint Kitts and Nevis",
		"Saint Lucia",
		"Saint Pierre and Miquelon",
		"Saint Vincent and the Grenadines",
		"Samoa",
		"San Marino",
		"Sao Tome and Principe",
		"Saudi Arabia",
		"Senegal",
		"Serbia",
		"Seychelles",
		"Sierra Leone",
		"Singapore",
		"Slovak Republic",
		"Slovenia",
		"Solomon Islands",
		"Somalia",
		"South Africa",
		"South Korea",
		"South Sudan",
		"Spain",
		"Sri Lanka",
		"Sudan",
		"Suriname",
		"Swaziland (Eswatini)",
		"Sweden",
		"Switzerland",
		"Syria",
		"Taiwan",
		"Tajikistan",
		"Tanzania",
		"Thailand",
		"Timor",
		"Togo",
		"Tonga",
		"Trinidad and Tobago",
		"Tunisia",
		"Turkey",
		"Turkmenistan",
		"Turks and Caicos Islands",
		"Tuvalu",
		"Uganda",
		"Ukraine",
		"United Arab Emirates",
		"United Kingdom",
		"United States Virgin Islands",
		"United States of America",
		"Uruguay",
		"Uzbekistan",
		"Vanuatu",
		"Venezuela",
		"Vietnam",
		"Wallis and Futuna Islands",
		"Yemen",
		"Zambia",
		"Zimbabwe",
		"monitoring",
	}

	utils.SendOkSliceResponse(c, nil, &countries)
	return
}

func (rc *RestController) ContractCompanies(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"HIDDEN", "HIDDEN Bulgaria", "TelWhiz"})
	return
}

func (rc *RestController) Timezones(c *gin.Context) {
	var timezones []string

	for i := -12; i <= 14; i++ {
		timezone := fmt.Sprintf("UTC/GMT %+d", i)

		timezones = append(timezones, timezone)
	}

	utils.SendOkSliceResponse(c, nil, &timezones)
	return
}

func (rc *RestController) HIDDENModels(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"SUBMIT IN", "SUBMIT OUT", "DELIVERY IN"})
	return

}

func (rc *RestController) HIDDENPeriods(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"monthly", "weekly", "biweekly"})
	return
}

func (rc *RestController) ContactTypes(c *gin.Context) {
	utils.SendOkSliceResponse(c, nil, &[]string{"Account Manager", "Contract Manager", "Rates Manager", "Support Manager"})
	return
}

func (rc *RestController) Tags(c *gin.Context) {
	var tags []partnercore.Tag

	for i := 1; i <= 5; i++ {
		tag := partnercore.Tag{
			ID:    strconv.Itoa(i),
			Label: fmt.Sprintf("Tag %d", i),
		}

		tags = append(tags, tag)
	}

	utils.SendOkSliceResponse(c, nil, &tags)
	return
}

func (rc *RestController) Currencies(c *gin.Context) {
	currencies := []string{"EUR", "USD", "GBP", "ILS"}

	utils.SendOkSliceResponse(c, nil, &currencies)
	return
}
