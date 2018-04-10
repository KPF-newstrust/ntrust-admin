package app

import (
	"fmt"
	//"strings"
	"net/http"
	"strconv"
	"errors"
	"time"
	"encoding/json"

	"github.com/labstack/echo"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)





type sanitizedSentence struct {
	SentenceId string	`bson:"sentence_id"`
	Sentence string		`bson:"sentence"`
}

type taggedWord struct {
	Word string		`bson:"word"`
	Tag string		`bson:"tag"`
}

func (tp *taggedWord) UnmarshalJSON(data []byte) error {
    var v []interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        fmt.Printf("Error whilde decoding %v\n", err)
        return err
	}
	
	tp.Word = v[0].(string)
	tp.Tag = v[1].(string)
	fmt.Printf("%s=>%s\n", tp.Word, tp.Tag)

    return nil
}

type sentenceItem struct {
	Word string			`bson:"word"`
	Category string		`bson:"category"`
	Sentence string		`bson:"sentence"`
	SentenceId string	`bson:"sentence_id"`
	SentenceIndex int	`bson:"sentence_index"`
}

type clusterKeyword struct {
	Keyword string	`bson:"keyword"`
	Count int		`bson:"count"`
}

type duplicationItem struct {
	Value float64	`bson:"value"`
	NewsId string	`bson:"news_id"`
}


func (tp *clusterKeyword) UnmarshalJSON(data []byte) error {
    var v []interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        fmt.Printf("Error whilde decoding %v\n", err)
        return err
	}
	
	tp.Keyword = v[0].(string)
	tp.Count = int(v[1].(float64))
    return nil
}

type MdbYear1News struct {
	Id bson.ObjectId	`bson:"_id,omitempty"`
	NewsId string		`bson:"news_id"`
	EmbagoDt int64		`bson:"_embargo_dt"`
	EnvelopTime int64	`bson:"_envelop_time"`
	Time int64			`bson:"_time"`
	Category []string	`bson:"category"`
	CategoryIncident []string	`bson:"category_incident"`
	CharacterCounter string		`bson:"character_counter"`
	CoiCode string		`bson:"coi_code"`
	Date string			`bson:"date"`
	Title string        `bson:"title"`
	Content string  	`bson:"content"`
	Dateline string		`bson:"dateline"`
	DatelineDate string	`bson:"dateline_date"`
	GenreInfo string	`bson:"genre_info"`

	Videos []string		`bson:"videos"`
	Articles []string	`bson:"articles"`
	Tags []string		`bson:"tags"`
	Images []string		`bson:"images"`
	KsCompany []string	`bson:"ks_company"`
	KsOrgan []string	`bson:"ks_organ"`
	KsPeople []string	`bson:"ks_people"`

	PageCategory string	`bson:"page_category"`
	PaperEdition string	`bson:"paper_edition"`
	PrintingPage string	`bson:"printing_page"`
	PrintingPageNo string	`bson:"printing_page_no"`
	ProviderCode string	`bson:"provider_code"`
	ProviderLinkExist string	`bson:"provider_link_exist"`
	ProviderLinkPage string	`bson:"provider_link_page"`
	ProviderMlinkPage string	`bson:"provider_mlink_page"`
	ProviderNewsId string	`bson:"provider_news_id"`
	ProviderSubject []string	`bson:"provider_subject"`

	PublisherCode string	`bson:"publisher_code"`
	PublishingStatus string	`bson:"publishing_status"`
	Revision string		`bson:"revision"`
	RightsLine string	`bson:"rightsline"`
	SeriesLine string	`bson:"seriesline"`

	SubTitle []string	`bson:"sub_title"`
	SubjectInfo1 []string	`bson:"subject_info1"`
	SubjectInfo2 []string	`bson:"subject_info2"`
	SubjectInfo3 []string	`bson:"subject_info3"`
	Urgency string	`bson:"urgency"`

	TmsNeLocation string	`bson:"TMS_NE_LOCATION"`
	TmsNeOrganization string	`bson:"TMS_NE_ORGANIZATION"`
	TmsNePerson string		`bson:"TMS_NE_PERSON"`
	TmsNeStream string		`bson:"TMS_NE_STREAM"`
	TmsRawStream string		`bson:"TMS_RAW_STREAM"`
	TmsSimilarity string	`bson:"TMS_SIMILARITY"`

	ContentsLength int	`bson:"contents_length"`	
	TitleLength int	`bson:"title_length"`
	WordsInTitle int	`bson:"words_in_title"`
	AttachmentLinks int	`bson:"attachment_links"`
	AttachmentMultimedia int	`bson:"attachment_multimedia"`
	AttachmentTables int	`bson:"attachment_tables"`
	AttachmentCardnews int	`bson:"attachment_cardnews"`	
	Figures []string	`bson:"figures"`

	Objects int	`bson:"objects"`
	AbuseTitles int	`bson:"abuse_titles"`	
	ArticleClustering int	`bson:"article_clustering"`
	PairQuotesInTitle int	`bson:"pair_quotes_in_title"`	

	SanitizedContent string	`bson:"sanitized_content"`
	SanitizedSentences []sanitizedSentence	`bson:"sanitized_sentences"`
	WordsPerSentence float64	`bson:"words_per_sentence"`
	TaggedContent []([]string)	`bson:"tagged_content"`
	TaggedTitle []taggedWord	`bson:"tagged_title"`
	NewsUrl string	`bson:"news_url"`
	Score int	`bson:"score"`
	Original bool	`bson:"original"`
	ClusterKeyword []clusterKeyword `bson:"clutster_keyword"`
	SanitizedBigkindsProvider string	`bson:"sanitized_bigkinds_provider"`
	SanitizedBigkindsCategory string	`bson:"sanitized_bigkinds_category"`
	Location []string	`bson:"location"`

	Role int	`bson:"role"`
	Job int	`bson:"job"`
	ConjunctionsInContent int	`bson:"conjunctions_in_content"`
	AdverbsInContent int	`bson:"adverbs_in_content"`
	AdjectivesInContent int	`bson:"adjectives_in_content"`
	WordNum int	`bson:"word_num"`
	PersonCount int	`bson:"person_count"`
	PersonSentence []sentenceItem `bson:"person"`
	LocationCount int	`bson:"location_count"`
	FiguresCount int	`bson:"figures_count"`
	HomonymCount int	`bson:"homonym_count"`
	EntityCount int	`bson:"entity_count"`
	PredictedCategory string	`bson:"predicted_category"`
	FillteringSanitizedByline string	`bson:"filltering_sanitized_byline"`
	OrganizationCount float64	`bson:"organization_count"`
	PoliticsCount float64	`bson:"politics_count"`
	ProductCount float64	`bson:"product_count"`
	EventCount float64	`bson:"event_count"`
	RoleCount float64	`bson:"role_count"`
	JobCount float64	`bson:"job_count"`
	ReplaceCount float64	`bson:"replace_count"`
	BigkindsByline []string	`bson:"bigkinds_byline"`
	BylineForSubmit string	`bson:"byline_for_submit"`	
	TaggedContentWithSS []taggedWord	`bson:"tagged_content_with_SS_"`
	TaggedTitleWithSS []taggedWord	`bson:"tagged_title_with_SS_"`
	TaggedContentNNP []sentenceItem	`bson:"tagged_content_NNP"`
	TaggedContentNNPSentence []sentenceItem	`bson:"tagged_content_NNP_sentence"`
	TaggedContentEntity []sentenceItem	`bson:"tagged_content_entity"`
	During int	`bson:"during"`
	TaggedContentEntitySentence []sentenceItem	`bson:"tagged_content_entity_sentence"`
	PersonWord []string	`bson:"person_word"`
	OrganizationWord []string	`bson:"organization_word"`
	PoliticsWord []string	`bson:"politics_word"`
	ProductWord []string	`bson:"product_word"`
	EventWord []string	`bson:"event_word"`
	LocationWord []string	`bson:"location_word"`

	EntitiesCount float64	`bson:"entities_count"`
	SanitizedTitle string	`bson:"sanitized_title"`	

	SanitizedTitleLength int	`bson:"sanitized_title_length"`
	TitleLengthForSubmit int	`bson:"title_length_for_submit"`

	Politics []string	`bson:"politics"`
	Event []string	`bson:"event"`
	Organization []string	`bson:"organization"`
	Product []string	`bson:"product"`

	Quotes []sentenceItem	`bson:"quotes"`
	Duplications []duplicationItem	`bson:"duplications"`
	ClusterId bson.ObjectId	`bson:"cluster_id"`
	ClusterIds []bson.ObjectId	`bson:"cluster_ids"`

	BigkindsFirstCategory string	`bson:"bigkinds_first_category"`
	ScoreSanitizedTitleLength float64	`bson:"score_sanitized_title_length"`
	ScoreSanitizedContentLength int	`bson:"score_sanitized_content_length"`
	ScoreSentencesNum int	`bson:"score_sentences_num"`
	ScoreAdjectivesPerSentence int	`bson:"score_adjectives_per_sentence"`
	ScoreConjunctionsPerSentence float64	`bson:"score_conjunctions_per_sentence"`
	ScoreAdverbsPerSentence float64	`bson:"score_adverbs_per_sentence"`
	ScoreAdverbsInTitle int	`bson:"score_adverbs_in_title"`
	ScoreAttachmentImages float64	`bson:"score_attachment_images"`
	ScoreQuotesInTitle int	`bson:"score_quotes_in_title"`
	ScoreQuestionInTitle int	`bson:"score_question_in_title"`
	ScoreExclamationInTitle int	`bson:"score_exclamation_in_title"`
	ScoreShockInTitle int	`bson:"score_shock_in_title"`
	ScoreArticleSimilarity float64	`bson:"score_article_similarity"`
	ScoreQuotesNum float64	`bson:"score_quotes_num"`
	ScoreQuotesMany int	`bson:"score_quotes_many"`
	ScoreQuotesLength float64	`bson:"score_quotes_length"`
	ScoreBulletin int	`bson:"score_bulletin"`
	ScorePlan int	`bson:"score_plan"`
	ScoreExclusive int	`bson:"score_exclusive"`
	ScoreByline float64	`bson:"score_byline"`
	SanitizedContentLength float64	`bson:"sanitized_content_length"`
	SentencesNum float64	`bson:"sentences_num"`
	AdjectivesPerSentence float64	`bson:"adjectives_per_sentence"`
	ConjunctionsPerSentence float64	`bson:"conjunctions_per_sentence"`
	AdverbsPerSentence float64	`bson:"adverbs_per_sentence"`
	AdverbsInTitle float64	`bson:"adverbs_in_title"`
	AttachmentImages float64	`bson:"attachment_images"`
	ArticleSimilarity float64	`bson:"article_similarity"`
	QuotesInTitle float64	`bson:"quotes_in_title"`
	QuestionInTitle float64	`bson:"question_in_title"`
	ExclamationInTitle float64	`bson:"exclamation_in_title"`
	ShockInTitle float64	`bson:"shock_in_title"`
	QuotesNum float64	`bson:"quotes_num"`
	QuotesMany float64	`bson:"quotes_many"`
	QuotesLength float64	`bson:"quotes_length"`
	Bulletin float64	`bson:"bulletin"`
	Plan float64	`bson:"plan"`
	Exclusive float64	`bson:"exclusive"`
	ScoreSum float64	`bson:"score_sum"`
	Byline []string       `bson:"byline"`
	DuplicationsMaxValue float64	`bson:"duplications_max_value"`
	DuplicationsMinValue float64	`bson:"duplications_min_value"`
	DbviewerUrl string	`bson:"dbviewer_url"`
}

type Year1ListItem struct {
	Id string			`json:"DT_RowId"`
	PubDate string  	`json:"pubDate"`
	MediaName string   	`json:"mediaName"`
	Title string        `json:"title"`
	//Byline string       `json:"byline"`
	Category string `json:"category"`
}

func getFirstStringOrEmpty(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	} else {
		return ""
	}
}

func (src *MdbYear1News) MakeListItem() Year1ListItem {
	return Year1ListItem{
		src.Id.Hex(),
		src.Date,
		src.SanitizedBigkindsProvider,
		src.Title,
		//src.FillteringSanitizedByline,
		src.BigkindsFirstCategory,
	}
}

func jsonYear1List(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	
	skipRows,_ := strconv.Atoi(ctx.FormValue("start"))
    fetchRows,_ := strconv.Atoi(ctx.FormValue("length"))
    if fetchRows <= 0 {
        fetchRows = 10
	}
	
    orderDir := ctx.FormValue("order[0][dir]")
    orderColumn,_ := strconv.Atoi(ctx.FormValue("order[0][column]"))

	filterMediaId := ctx.FormValue("columns[1][search][value]")
	filterCategory := ctx.FormValue("columns[2][search][value]")
	filterTitle := ctx.FormValue("columns[3][search][value]")
	filterNewsId := ctx.FormValue("search[value]")

	var json JqDataTable
    json.Draw, _ = strconv.Atoi(ctx.FormValue("draw"))

	// 검색
	var err error
	var findParam = make(bson.M)

	if filterNewsId != "" {
		findParam["news_id"] = bson.RegEx{filterNewsId, ""}
	} else {
		ymd := ctx.FormValue("ymd")
		dtBegin, err := time.ParseInLocation("2006-01-02", ymd, tzLocation)
		if err == nil {
			date8 := dtBegin.Format("20060102")
			findParam["date"] = bson.RegEx{date8, ""}
		}

		if filterMediaId != "" {
			findParam["sanitized_bigkinds_provider"] = bson.RegEx{filterMediaId, ""}
		}
		if filterCategory != "" {
			findParam["bigkinds_first_category"] = bson.RegEx{filterCategory, ""}
		}
		if filterTitle != "" {
			findParam["title"] = bson.RegEx{filterTitle, ""}
		}
	}
	
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoOldDB.With(sessionCopy).C(COLLNAME_YEAR1NEWS)

	json.Total, err = coll.Count()
	if err != nil {
		panic(err)
	}

	qry := coll.Find(&findParam).Select(bson.M{"title":1, "date":1,
		"sanitized_bigkinds_provider":1, "bigkinds_first_category":1 })

	json.Filtered, err = qry.Count()
	if err != nil {
		panic(err)
	}

	// 정렬
	var orderFieldName string
	if orderDir == "desc" {
		orderFieldName = "-"
	} else {
		orderFieldName = ""
	}

	switch (orderColumn) {
		case 0:
			orderFieldName += "date"
		case 1:
			orderFieldName += "sanitized_bigkinds_provider"
		case 2:
			orderFieldName += "bigkinds_first_category"
		//case 3:
		//	orderFieldName += "filltering_sanitized_byline"
		case 3:
			fallthrough
		default:
			orderFieldName += "title"
	}

	var results []MdbYear1News
	err = qry.Sort(orderFieldName).Skip(skipRows).Limit(fetchRows).All(&results)
	if err != nil {
		panic(err)
	}

    json.Rows = make([]interface{}, len(results))
	for i,v := range results {
		json.Rows[i] = v.MakeListItem();
	}
	
	return ctx.JSON(http.StatusOK, json)
}

func showYear1List(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)
	ctx.model["lmenu"] = "t3a"
	ctx.model["targetDate"] = "2016-06-01"
	return ctx.renderTemplate("year1_list.html")
}

func showYear1Detail(_ctx echo.Context) error {
	ctx := _ctx.(*AuthContext)

	var where bson.M
	objId := ctx.QueryParam("obj")
	newsId := ctx.QueryParam("nws")
	if objId != "" {
		where = bson.M{"_id": bson.ObjectIdHex(objId)}
	} else if newsId != "" {
		where = bson.M{"news_id": newsId}
	} else {
		return ctx.showAdminError(errors.New("Invalid ID"))
	}

	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()
	coll := mgoOldDB.With(sessionCopy).C(COLLNAME_YEAR1NEWS)

	var data MdbYear1News
	err := coll.Find(where).One(&data)
	if err != nil {
		return ctx.showAdminError(errors.New("Object not found"))
	}

	ctx.model["item"] = data
	ctx.model["lmenu"] = "t3a"
	
	return ctx.renderTemplate("year1_view.html")
}


func setupYear1Routes(grp *echo.Group) {

	// 수집 DB
	grp.GET("/year1/list", showYear1List)
	grp.POST("/year1/list.json", jsonYear1List)
	grp.GET("/year1/view", showYear1Detail)
}
