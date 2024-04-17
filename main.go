package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func AddTable(client *lark.Client, ChatId string, Name string, Doc string) {
	req := larkim.NewCreateChatTabReqBuilder().
		ChatId(ChatId).
		Body(larkim.NewCreateChatTabReqBodyBuilder().
			ChatTabs([]*larkim.ChatTab{
				larkim.NewChatTabBuilder().
					TabName(Name).
					TabType(`doc`).
					TabContent(larkim.NewChatTabContentBuilder().
						Doc(Doc).
						Build()).
					Build(),
			}).
			Build()).
		Build()
	resp, err := client.Im.ChatTab.Create(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
func CreateGroup(client *lark.Client, Name string, OwnerId string, UserIdList []string, Description string) string {
	req := larkim.NewCreateChatReqBuilder().
		SetBotManager(true).
		Body(larkim.NewCreateChatReqBodyBuilder().
			Name(Name).
			Description(Description).
			OwnerId(OwnerId).
			BotIdList([]string{"cli_a136ebe23478900b"}).
			UserIdList(UserIdList).
			Build()).
		Build()

	resp, err := client.Im.Chat.Create(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return ""
	}
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.ChatId
}
func CopyFile(client *lark.Client, FatherFolder string, FileToken string, Name string, Type string) (string, string) {
	req := larkdrive.NewCopyFileReqBuilder().
		FileToken(FileToken).
		Body(larkdrive.NewCopyFileReqBodyBuilder().
			Name(Name).
			Type(Type).
			FolderToken(FatherFolder).
			Build()).
		Build()
	resp, err := client.Drive.File.Copy(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return "", ""
	}
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.File.Url, *resp.Data.File.Token
}
func CreateDocx(client *lark.Client, FolderToken string, Title string) string {
	req := larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken(FolderToken).
			Title(Title).
			Build()).
		Build()
	resp, err := client.Docx.Document.Create(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return ""
	}
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.Document.DocumentId
}
func AddFilePrem(client *lark.Client, FileToken string, Type string, MemberType string, MemberId string, Prem string) {
	req := larkdrive.NewCreatePermissionMemberReqBuilder().
		Token(FileToken).
		Type(Type).
		BaseMember(larkdrive.NewBaseMemberBuilder().
			MemberType(MemberType).
			MemberId(MemberId).
			Perm(Prem).
			Build()).
		Build()
	resp, err := client.Drive.PermissionMember.Create(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
func CreateFolder(client *lark.Client, FatherFolder string, Name string) string {
	req := larkdrive.NewCreateFolderFileReqBuilder().
		Body(larkdrive.NewCreateFolderFileReqBodyBuilder().
			Name(Name).
			FolderToken(FatherFolder).
			Build()).
		Build()
	resp, err := client.Drive.File.CreateFolder(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return ""
	}
	fmt.Println(larkcore.Prettify(resp))
	return *resp.Data.Token
}
func SearchRecord(client *lark.Client, AppToken string, TableId string, RecordId string) map[string]interface{} {
	req := larkbitable.NewGetAppTableRecordReqBuilder().
		AppToken(AppToken).
		TableId(TableId).
		RecordId(RecordId).
		Build()
	resp, err := client.Bitable.AppTableRecord.Get(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil
	}
	fmt.Println(larkcore.Prettify(resp))
	return resp.Data.Record.Fields
}

func UpdateRecord(client *lark.Client, AppToken string, TableId string, RecordId string, Fields map[string]interface{}) {
	req := larkbitable.NewUpdateAppTableRecordReqBuilder().
		AppToken(AppToken).
		TableId(TableId).
		RecordId(RecordId).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(Fields).
			Build()).
		Build()
	resp, err := client.Bitable.AppTableRecord.Update(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}

func getUserInfo(client *lark.Client, UserId string) *larkcontact.User {
	req := larkcontact.NewGetUserReqBuilder().
		UserId(UserId).
		Build()
	resp, err := client.Contact.User.Get(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil
	}
	return resp.Data.User
}

var FolderList map[string]map[string]string

func GetFolderList(client *lark.Client, FolderToken string) {
	req := larkdrive.NewListFileReqBuilder().
		FolderToken(FolderToken).
		OrderBy(`EditedTime`).
		Direction(`DESC`).
		Build()
	resp, err := client.Drive.File.List(context.Background(), req)

	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(resp))
	if FolderList[FolderToken] == nil {
		FolderList[FolderToken] = make(map[string]string)
	}
	for _, value := range resp.Data.Files {
		FolderList[FolderToken][*value.Name] = *value.Token
	}
	if len(resp.Data.Files) > 0 {
		FolderList[FolderToken]["ParentToken"] = *resp.Data.Files[0].ParentToken
	}
}
func DelFile(client *lark.Client, Token string, Type string) {
	req := larkdrive.NewDeleteFileReqBuilder().
		FileToken(Token).
		Type(Type).
		Build()

	resp, err := client.Drive.File.Delete(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}

func getLarkClient() *lark.Client {
	return lark.NewClient(os.Getenv("appId"), os.Getenv("appSecret"))
}
func newFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	AppToken := r.Form.Get("AppToken")
	TableId := r.Form.Get("TableId")
	RecordId := r.Form.Get("RecordId")
	ChallType := r.Form.Get("ChallType")
	ChallName := r.Form.Get("ChallName")
	if AppToken == "" || TableId == "" || RecordId == "" || ChallType == "" || ChallName == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	var GroupID, FolderToken string
	err = db.QueryRow("SELECT GroupID, FolderToken FROM `game` WHERE FileToken=?", AppToken).Scan(&GroupID, &FolderToken)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	client := getLarkClient()
	FileToken := CreateDocx(client, FolderToken, ChallType+"-"+ChallName)
	Url := "https://szuaurora.feishu.cn/docx/" + FileToken
	AddFilePrem(client, FileToken, "docx", "openchat", GroupID, "full_access")
	Fields := SearchRecord(client, AppToken, TableId, RecordId)
	Fields["题目文档"] = Url
	UpdateRecord(client, AppToken, TableId, RecordId, Fields)
}
func newGameHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	AppToken := r.Form.Get("AppToken")
	TableId := r.Form.Get("TableId")
	RecordId := r.Form.Get("RecordId")
	if AppToken == "" || TableId == "" || RecordId == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	client := getLarkClient()
	Fields := SearchRecord(client, AppToken, TableId, RecordId)
	GroupID := CreateGroup(client, Fields["比赛名称"].(string), os.Getenv("masterId"), []string{}, "")
	Groups, _ := Fields["比赛群"].([]interface{})
	if len(Groups) == 0 {
		Fields["比赛群"] = []interface{}{
			map[string]interface{}{
				"id": GroupID,
			},
		}
	} else {
		Group := Groups[0].(map[string]interface{})
		Group["id"] = GroupID
	}
	UpdateRecord(client, AppToken, TableId, RecordId, Fields)
	FolderToken := CreateFolder(client, FolderList[""]["比赛数据"], Fields["比赛名称"].(string))
	FileUrl, FileToken := CopyFile(client, FolderToken, FolderList[FolderList[""]["模板"]]["比赛表格模板"], Fields["比赛名称"].(string), "bitable")
	AddFilePrem(client, FileToken, "bitable", "openchat", GroupID, "full_access")
	AddTable(client, GroupID, "比赛表格", FileUrl)
	_, _ = db.Exec("INSERT INTO `game` (GroupId,FolderToken,FileToken) VALUES (?,?,?)", GroupID, FolderToken, FileToken)
}

func callBackhandler() *dispatcher.EventDispatcher {
	handler := dispatcher.NewEventDispatcher(os.Getenv("verToken"), os.Getenv("eventKey"))
	handler.OnP2CardNewProtocalURLPreviewGet(func(ctx context.Context, event *dispatcher.URLPreviewGetEvent) (*dispatcher.URLPreviewGetResponse, error) {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.Event.Context.URL)
		u, _ := url.Parse(event.Event.Context.URL)
		queryParams := u.Query()
		tParam := queryParams.Get("t")
		if strings.Contains(tParam, "$name") {
			client := getLarkClient()
			Users := getUserInfo(client, event.Event.Operator.OpenID)
			Name := ""
			if Users != nil {
				Name = *Users.Name
			}
			tParam = strings.ReplaceAll(tParam, "$name", Name)
		}
		return &dispatcher.URLPreviewGetResponse{
			Inline: &dispatcher.Inline{
				Title: tParam,
			},
		}, nil
	})
	return handler
}
func urlRedirect(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("u")
	if u != "" {
		http.Redirect(w, r, u, http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func main() {
	{
		if err := init_db(); err != nil {
			fmt.Println("初始化数据库失败")
			return
		}
		if os.Getenv("appId") == "" || os.Getenv("appSecret") == "" || os.Getenv("masterId") == "" {
			fmt.Println("缺失appId或appSecret或masterId")
			return
		}
		client := getLarkClient()
		FolderList = make(map[string]map[string]string)
		GetFolderList(client, "")
		for _, value := range []string{"比赛数据", "模板"} {
			if FolderList[""][value] == "" {
				CreateFolder(client, FolderList[""]["ParentToken"], value)
			}
		}
		GetFolderList(client, "")
		GetFolderList(client, FolderList[""]["模板"])
	}
	http.HandleFunc("/newfile", newFileHandler)
	http.HandleFunc("/newgame", newGameHandler)
	http.HandleFunc("/callback", httpserverext.NewEventHandlerFunc(callBackhandler(),
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	http.HandleFunc("/url", urlRedirect)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
