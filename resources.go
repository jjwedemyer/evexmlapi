package evexmlapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

// Params is a map for a request's query parameters
type Params map[string][]string

// Resource is a stuct for a resource
type Resource struct {
	path           string
	cacheDuration  int64	// in seconds
	requiredParams []string
	optionalParams []string
	parse          bodyParser
	model
	api
}

type api struct {
	protocol string
	method   string
	baseURL  string
	port     int
}

var (
	xmlAPI = api{
		protocol: "https://",
		method:   "GET",
		baseURL:  "api.eveonline.com/",
		port:     443,
	}
)

type paramError struct {
	params       Params
	resourcePath string
}

func (pe paramError) Error() string {
	return fmt.Sprintf("%q", pe.params) + " not valid for " + pe.resourcePath
}

// VerifyParams verifies the required params are available
func (r Resource) verifyParams(params Params) error {
	if (len(r.requiredParams) + len(r.optionalParams)) < len(params) {
		return paramError{params: params, resourcePath: r.path}
	}
	for _, v := range r.requiredParams {
		if _, ok := params[v]; !ok || params[v][0] == "" {
			return paramError{params: params, resourcePath: r.path}
		}
	}
	return nil
}


// XMLtoJSON takes xml data and a format struct
// then converts it into a json byte array
func XMLtoJSON(xmldata []byte, r Resource) ([]byte, error) {
	if r.format == nil {
		log.Printf("Resource %+v does not have a data format", r)
		return xmldata, nil
	}
	f := r.format
	err := xml.Unmarshal(xmldata, f)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Resource) SetDataFormat(f interface{}) {
	r.format = f
}

var (
	noParams     = []string{}
	accountLevel = []string{"keyID", "vCode"}
	charLevel    = append(accountLevel, "characterID")
	corpLevel	= append(accountLevel, "corporationID")
)

// NewCharSkillQueue is a constructor for the NewCharSkillQueue resource
func NewCharSkillQueue() Resource {
	return Resource{
		path:           "char/SkillQueue.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		model:          model{format: &SkillQueueFormat{}},
		api:            xmlAPI,
	}
}

// NewServerStatus is a constructor for the ServerStatus resource
func NewServerStatus() Resource {
	return Resource{
		path:           "server/ServerStatus.xml.aspx",
		cacheDuration:  180,
		requiredParams: noParams,
		model:          model{format: &ServerStatusFormat{}},
		api:            xmlAPI,
	}
}

// NewAPIKeyInfo is a constructor for the APIKeyInfo resource
func NewAPIKeyInfo() Resource {
	return Resource{
		path:           "account/APIKeyInfo.xml.aspx",
		cacheDuration:  300,
		requiredParams: accountLevel,
		model:          model{format: &APIKeyInfoFormat{}},
		api:            xmlAPI,
	}
}

// NewAccountStatus is a constructor for the NewAccountStatus resource
func NewAccountStatus() Resource {
	return Resource{
		path:           "account/AccountStatus.xml.aspx",
		cacheDuration:  3600,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewCharacters is a constructor for the NewCharacters resource
func NewCharacters() Resource {
	return Resource{
		path:           "account/Characters.xml.aspx",
		cacheDuration:  3600,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewAPICallList is a constructor for the NewAPICallList resource
func NewAPICallList() Resource {
	return Resource{
		path:           "api/CallList.xml.aspx",
		cacheDuration:  21600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewCharAccountBalance is a constructor for the NewCharAccountBalance resource
func NewCharAccountBalance() Resource {
	return Resource{
		path:           "char/AccountBalance.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharAssetList is a constructor for the NewCharAssetList resource
func NewCharAssetList() Resource {
	return Resource{
		path:           "char/AssetList.xml.aspx",
		cacheDuration:  7200,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharBlueprints is a constructor for the NewCharBlueprints resource
func NewCharBlueprints() Resource {
	return Resource{
		path:           "char/Blueprints.xml.aspx",
		cacheDuration:  43200,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharBookmarks is a constructor for the NewCharBookmarks resource
func NewCharBookmarks() Resource {
	return Resource{
		path:           "char/Bookmarks.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharCalEventAttendees is a constructor for the NewCharCalEventAttendees resource
func NewCharCalEventAttendees() Resource {
	return Resource{
		path:           "char/CalendarEventAttendees.xml.aspx",
		cacheDuration:  600,
		requiredParams: append(charLevel, "eventIDs"),
		api:            xmlAPI,
	}
}

// NewCharCharacterSheet is a constructor for the NewCharCharacterSheet resource
func NewCharCharacterSheet() Resource {
	return Resource{
		path:           "char/CharacterSheet.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharChatChannels is a constructor for the NewCharChatChannels resource
func NewCharChatChannels() Resource {
	return Resource{
		path:           "char/ChatChannels.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharContactList is a constructor for the NewCharContactList resource
func NewCharContactList() Resource {
	return Resource{
		path:           "char/ContactList.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharContactNotifications is a constructor for the NewCharContactNotifications resource
func NewCharContactNotifications() Resource {
	return Resource{
		path:           "char/ContactNotifications.xml.aspx",
		cacheDuration:  1800,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharContractBids is a constructor for the NewCharContractBids resource
func NewCharContractBids() Resource {
	return Resource{
		path:           "char/ContractBids.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharContractItems is a constructor for the NewCharContractItems resource
func NewCharContractItems() Resource {
	return Resource{
		path:           "char/ContractItems.xml.aspx",
		cacheDuration:  315360000,
		requiredParams: append(charLevel, "contractID"),
		api:            xmlAPI,
	}
}

// NewCharContracts is a constructor for the NewCharContracts resource
func NewCharContracts() Resource {
	return Resource{
		path:           "char/Contracts.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(charLevel, "contractID"),
		api:            xmlAPI,
	}
}

// NewCharFacWarStats is a constructor for the NewCharFacWarStats resource
func NewCharFacWarStats() Resource {
	return Resource{
		path:           "char/FacWarStats.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharIndustryJobsHistory is a constructor for the NewCharIndustryJobsHistory resource
func NewCharIndustryJobsHistory() Resource {
	return Resource{
		path:           "char/IndustryJobsHistory.xml.aspx",
		cacheDuration:  21600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharMailBodies is a constructor for the NewCharMailBodies resource
func NewCharMailBodies() Resource {
	return Resource{
		path:           "char/MailBodies.xml.aspx",
		cacheDuration:  315360000,
		requiredParams: append(charLevel, "IDs"),
		api:            xmlAPI,
	}
}

// NewCharMailingLists is a constructor for the NewCharMailingLists resource
func NewCharMailingLists() Resource {
	return Resource{
		path:           "char/mailinglists.xml.aspx",
		cacheDuration:  21600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharMailMessages is a constructor for the NewCharMailMessages resource
func NewCharMailMessages() Resource {
	return Resource{
		path:           "char/MailMessages.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharMarketOrders is a constructor for the NewCharMarketOrders resource
func NewCharMarketOrders() Resource {
	return Resource{
		path:           "char/MarketOrders.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(charLevel, "orderID"),
		api:            xmlAPI,
	}
}

// NewCharNotifications is a constructor for the NewCharNotifications resource
func NewCharNotifications() Resource {
	return Resource{
		path:           "char/Notifications.xml.aspx",
		cacheDuration:  1800,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharNotificationTexts is a constructor for the NewCharNotificationTexts resource
func NewCharNotificationTexts() Resource {
	return Resource{
		path:           "char/NotificationTexts.xml.aspx",
		cacheDuration:  315360000,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharPlanetaryColonies is a constructor for the NewCharPlanetaryColonies resource
func NewCharPlanetaryColonies() Resource {
	return Resource{
		path:           "char/PlanetaryColonies.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharPlanetaryLinks is a constructor for the NewCharPlanetaryLinks resource
func NewCharPlanetaryLinks() Resource {
	return Resource{
		path:           "char/PlanetaryLinks.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(charLevel, "planetID"),
		api:            xmlAPI,
	}
}

// NewCharPlanetaryPins is a constructor for the NewCharPlanetaryPins resource
func NewCharPlanetaryPins() Resource {
	return Resource{
		path:           "char/PlanetaryPins.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(charLevel, "planetID"),
		api:            xmlAPI,
	}
}

// NewCharPlanetaryRoutes is a constructor for the NewCharPlanetaryRoutes resource
func NewCharPlanetaryRoutes() Resource {
	return Resource{
		path:           "char/PlanetaryRoutes.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(charLevel, "planetID"),
		api:            xmlAPI,
	}
}

// NewCharResearch is a constructor for the NewCharResearch resource
func NewCharResearch() Resource {
	return Resource{
		path:           "char/Research.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharSkillInTraining is a constructor for the NewCharSkillInTraining resource
func NewCharSkillInTraining() Resource {
	return Resource{
		path:           "char/SkillInTraining.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharStandings is a constructor for the NewCharStandings resource
func NewCharStandings() Resource {
	return Resource{
		path:           "char/Standings.xml.aspx",
		cacheDuration:  10800,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharUpcomingCalendarEvents is a constructor for the NewCharUpcomingCalendarEvents resource
func NewCharUpcomingCalendarEvents() Resource {
	return Resource{
		path:           "char/UpcomingCalendarEvents.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCharWalletJournal is a constructor for the NewCharWalletJournal resource
func NewCharWalletJournal() Resource {
	return Resource{
		path:           "char/WalletJournal.xml.aspx",
		cacheDuration:  1800,
		requiredParams: append(charLevel, "accountKey", "fromID", "rowCount"),
		api:            xmlAPI,
	}
}

// NewCharWalletTransactions is a constructor for the NewCharWalletTransactions resource
func NewCharWalletTransactions() Resource {
	return Resource{
		path:           "char/WalletTransactions.xml.aspx",
		cacheDuration:  1800,
		requiredParams: append(charLevel, "accountKey", "fromID", "rowCount"),
		api:            xmlAPI,
	}
}

// NewCorpAssetList is a constructor for the NewCorpAssetList resource
func NewCorpAssetList() Resource {
	return Resource{
		path:						"corp/assetlist.xml.aspx",
		cacheDuration:	7200,
		requiredParams: charLevel,
		api:						xmlAPI,
	}
}

// NewCorpAccountBalance is a constructor for the NewCorpAccountBalance resource
func NewCorpAccountBalance() Resource {
	return Resource{
		path:           "corp/AccountBalance.xml.aspx",
		cacheDuration:  900,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewCorpBookmarks is a constructor for the NewCorpBookmarks resource
func NewCorpBookmarks() Resource {
	return Resource{
		path:           "corp/Bookmarks.xml.aspx",
		cacheDuration:  3600,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewCorpCorporationSheet is a constructor for the NewCorpCorporationSheet resource
func NewCorpCorporationSheet() Resource {
	return Resource{
		path:           "corp/CorporationSheet.xml.aspx",
		cacheDuration:  3600,
		requiredParams: []string{"corporationID"},
		optionalParams:	accountLevel,
		api:            xmlAPI,
	}
}

// NewCorpIndustryJobs is a constructor for the NewCorpIndustryJobs resource
func NewCorpIndustryJobs() Resource {
	return Resource{
		path:           "corp/IndustryJobs.xml.aspx",
		cacheDuration:  900,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewCorpIndustryJobsHistory is a constructor for the NewCorpIndustryJobsHistory resource
func NewCorpIndustryJobsHistory() Resource {
	return Resource{
		path:           "corp/IndustryJobsHistory.xml.aspx",
		cacheDuration:  21600,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewCorpMemberTracking is a constructor for the NewCorpMemberTracking resource
func NewCorpMemberTracking() Resource {
	return Resource{
		path:           "corp/MemberTracking.xml.aspx",
		cacheDuration:  19800,
		requiredParams: append(accountLevel, "extended"),
		api:            xmlAPI,
	}
}

// NewCorpStarbaseList is a constructor for the NewCorpStarbaseList resource
func NewCorpStarbaseList() Resource {
	return Resource{
		path:           "corp/StarbaseList.xml.aspx",
		cacheDuration:  21600,
		requiredParams: accountLevel,
		api:            xmlAPI,
	}
}

// NewCorpStarbaseDetail is a constructor for the NewCorpStarbaseDetail resource
func NewCorpStarbaseDetail() Resource {
	return Resource{
		path:           "corp/StarbaseDetail.xml.aspx",
		cacheDuration:  3600,
		requiredParams: append(accountLevel, "itemID"),
		api:            xmlAPI,
	}
}

// NewCorpWalletJournal is a constructor for the NewCorpWalletJournal resource
func NewCorpWalletJournal() Resource {
	return Resource{
		path:           "corp/WalletJournal.xml.aspx",
		cacheDuration:  1800,
		requiredParams: append(accountLevel, "accountKey", "fromID", "rowCount"),
		api:            xmlAPI,
	}
}

// NewCorpWalletTransactions is a constructor for the NewCorpWalletTransactions resource
func NewCorpWalletTransactions() Resource {
	return Resource{
		path:           "corp/WalletTransactions.xml.aspx",
		cacheDuration:  1800,
		requiredParams: append(accountLevel, "accountKey", "fromID", "rowCount"),
		api:            xmlAPI,
	}
}

// NewEveAllianceList is a constructor for the NewEveAllianceList resource
func NewEveAllianceList() Resource {
	return Resource{
		path:           "eve/AllianceList.xml.aspx",
		cacheDuration:  3600,
		optionalParams: []string{"version"},
		api:            xmlAPI,
	}
}

// NewEveCharacterAffiliation is a constructor for the NewEveCharacterAffiliation resource
func NewEveCharacterAffiliation() Resource {
	return Resource{
		path:           "eve/CharacterAffiliation.xml.aspx",
		cacheDuration:  3600,
		requiredParams: []string{"ids"},
		api:            xmlAPI,
	}
}

// NewEveCharacterID is a constructor for the NewEveCharacterID resource
func NewEveCharacterID() Resource {
	return Resource{
		path:           "eve/CharacterID.xml.aspx",
		cacheDuration:  43200,
		requiredParams: []string{"names"},
		api:            xmlAPI,
	}
}

// NewEveCharacterInfo is a constructor for the NewEveCharacterInfo resource
func NewEveCharacterInfo() Resource {
	return Resource{
		path:           "eve/CharacterInfo.xml.aspx",
		cacheDuration:  3600,
		requiredParams: charLevel,
		api:            xmlAPI,
	}
}

// NewEveCharacterName is a constructor for the NewEveCharacterName resource
func NewEveCharacterName() Resource {
	return Resource{
		path:           "eve/CharacterName.xml.aspx",
		cacheDuration:  2592000,
		requiredParams: []string{"ids"},
		api:            xmlAPI,
	}
}

// NewEveConquerableStationList is a constructor for the NewEveConquerableStationList resource
func NewEveConquerableStationList() Resource {
	return Resource{
		path:           "eve/ConquerableStationList.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewEveErrorList is a constructor for the NewEveErrorList resource
func NewEveErrorList() Resource {
	return Resource{
		path:           "eve/ErrorList.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewEveRefTypes is a constructor for the NewEveRefTypes resource
func NewEveRefTypes() Resource {
	return Resource{
		path:           "eve/RefTypes.xml.aspx",
		cacheDuration:  86400,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewEveTypeName is a constructor for the NewEveTypeName resource
func NewEveTypeName() Resource {
	return Resource{
		path:           "eve/TypeName.xml.aspx",
		cacheDuration:  3600,
		requiredParams: []string{"ids"},
		api:            xmlAPI,
	}
}

// NewMapFacWarSystems is a constructor for the NewMapFacWarSystems resource
func NewMapFacWarSystems() Resource {
	return Resource{
		path:           "map/FacWarSystems.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewMapJumps is a constructor for the NewMapJumps resource
func NewMapJumps() Resource {
	return Resource{
		path:           "map/Jumps.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewMapKills is a constructor for the NewMapKills resource
func NewMapKills() Resource {
	return Resource{
		path:           "map/kills.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}

// NewMapSovereignty is a constructor for the NewMapSovereignty resource
func NewMapSovereignty() Resource {
	return Resource{
		path:           "map/Sovereignty.xml.aspx",
		cacheDuration:  3600,
		requiredParams: noParams,
		api:            xmlAPI,
	}
}
