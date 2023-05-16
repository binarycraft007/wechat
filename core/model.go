package core

type User struct {
	Uin               int    `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgURL        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
}

type InitRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
}

type StatusNotifyRequest struct {
	BaseRequest  BaseRequest `json:"BaseRequest"`
	Code         int         `json:"Code"`
	FromUserName string      `json:"FromUserName"`
	ToUserName   string      `json:"ToUserName"`
	ClientMsgId  int64       `json:"ClientMsgId"`
}

type BatchGetContactRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
	Count       int         `json:"Count"`
	List        []Contact   `json:"List"`
}

type SyncRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
	SyncKey     SyncKey     `json:"SyncKey"`
	RR          int64       `json:"rr"`
}

type BaseRequest struct {
	Uin      int64  `json:"Uin"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

type BaseResponse struct {
	Ret    int    `json:"Ret"`
	ErrMsg string `json:"ErrMsg"`
}

type Contact struct {
	Uin              int    `json:"Uin"`
	UserName         string `json:"UserName"`
	NickName         string `json:"NickName"`
	HeadImgURL       string `json:"HeadImgUrl"`
	ContactFlag      int    `json:"ContactFlag"`
	MemberCount      int    `json:"MemberCount"`
	MemberList       []any  `json:"MemberList"`
	RemarkName       string `json:"RemarkName"`
	HideInputBarFlag int    `json:"HideInputBarFlag"`
	Sex              int    `json:"Sex"`
	Signature        string `json:"Signature"`
	VerifyFlag       int    `json:"VerifyFlag"`
	OwnerUin         int    `json:"OwnerUin"`
	PYInitial        string `json:"PYInitial"`
	PYQuanPin        string `json:"PYQuanPin"`
	RemarkPYInitial  string `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string `json:"RemarkPYQuanPin"`
	StarFriend       int    `json:"StarFriend"`
	AppAccountFlag   int    `json:"AppAccountFlag"`
	Statues          int    `json:"Statues"`
	AttrStatus       int    `json:"AttrStatus"`
	Province         string `json:"Province"`
	City             string `json:"City"`
	Alias            string `json:"Alias"`
	SnsFlag          int    `json:"SnsFlag"`
	UniFriend        int    `json:"UniFriend"`
	DisplayName      string `json:"DisplayName"`
	ChatRoomID       int    `json:"ChatRoomId"`
	KeyWord          string `json:"KeyWord"`
	EncryChatRoomID  string `json:"EncryChatRoomId"`
	IsOwner          int    `json:"IsOwner"`
}

type BatchGetContactResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	Count        int          `json:"Count"`
	ContactList  []Contact    `json:"ContactList"`
}

type SyncKey struct {
	Count int `json:"Count"`
	List  []struct {
		Key int `json:"Key"`
		Val int `json:"Val"`
	} `json:"List"`
}

type InitResponse struct {
	BaseResponse        BaseResponse `json:"BaseResponse"`
	Count               int          `json:"Count"`
	ContactList         []Contact    `json:"ContactList"`
	SyncKey             SyncKey      `json:"SyncKey"`
	User                User         `json:"User"`
	ChatSet             string       `json:"ChatSet"`
	SKey                string       `json:"SKey"`
	ClientVersion       int          `json:"ClientVersion"`
	SystemTime          int          `json:"SystemTime"`
	GrayScale           int          `json:"GrayScale"`
	InviteStartCount    int          `json:"InviteStartCount"`
	MPSubscribeMsgCount int          `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []any        `json:"MPSubscribeMsgList"`
	ClickReportInterval int          `json:"ClickReportInterval"`
}

type StatusNotifyResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MsgID        string       `json:"MsgID"`
}

type SyncResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	AddMsgCount  int          `json:"AddMsgCount"`
	AddMsgList   []struct {
		MsgID                string `json:"MsgId"`
		FromUserName         string `json:"FromUserName"`
		ToUserName           string `json:"ToUserName"`
		MsgType              int    `json:"MsgType"`
		Content              string `json:"Content"`
		Status               int    `json:"Status"`
		ImgStatus            int    `json:"ImgStatus"`
		CreateTime           int    `json:"CreateTime"`
		VoiceLength          int    `json:"VoiceLength"`
		PlayLength           int    `json:"PlayLength"`
		FileName             string `json:"FileName"`
		FileSize             string `json:"FileSize"`
		MediaID              string `json:"MediaId"`
		URL                  string `json:"Url"`
		AppMsgType           int    `json:"AppMsgType"`
		StatusNotifyCode     int    `json:"StatusNotifyCode"`
		StatusNotifyUserName string `json:"StatusNotifyUserName"`
		RecommendInfo        struct {
			UserName   string `json:"UserName"`
			NickName   string `json:"NickName"`
			QQNum      int    `json:"QQNum"`
			Province   string `json:"Province"`
			City       string `json:"City"`
			Content    string `json:"Content"`
			Signature  string `json:"Signature"`
			Alias      string `json:"Alias"`
			Scene      int    `json:"Scene"`
			VerifyFlag int    `json:"VerifyFlag"`
			AttrStatus int    `json:"AttrStatus"`
			Sex        int    `json:"Sex"`
			Ticket     string `json:"Ticket"`
			OpCode     int    `json:"OpCode"`
		} `json:"RecommendInfo"`
		ForwardFlag int `json:"ForwardFlag"`
		AppInfo     struct {
			AppID string `json:"AppID"`
			Type  int    `json:"Type"`
		} `json:"AppInfo"`
		HasProductID  int    `json:"HasProductId"`
		Ticket        string `json:"Ticket"`
		ImgHeight     int    `json:"ImgHeight"`
		ImgWidth      int    `json:"ImgWidth"`
		SubMsgType    int    `json:"SubMsgType"`
		NewMsgID      int64  `json:"NewMsgId"`
		OriContent    string `json:"OriContent"`
		EncryFileName string `json:"EncryFileName"`
	} `json:"AddMsgList"`
	ModContactCount        int   `json:"ModContactCount"`
	ModContactList         []any `json:"ModContactList"`
	DelContactCount        int   `json:"DelContactCount"`
	DelContactList         []any `json:"DelContactList"`
	ModChatRoomMemberCount int   `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []any `json:"ModChatRoomMemberList"`
	Profile                struct {
		BitFlag  int `json:"BitFlag"`
		UserName struct {
			Buff string `json:"Buff"`
		} `json:"UserName"`
		NickName struct {
			Buff string `json:"Buff"`
		} `json:"NickName"`
		BindUin   int `json:"BindUin"`
		BindEmail struct {
			Buff string `json:"Buff"`
		} `json:"BindEmail"`
		BindMobile struct {
			Buff string `json:"Buff"`
		} `json:"BindMobile"`
		Status            int    `json:"Status"`
		Sex               int    `json:"Sex"`
		PersonalCard      int    `json:"PersonalCard"`
		Alias             string `json:"Alias"`
		HeadImgUpdateFlag int    `json:"HeadImgUpdateFlag"`
		HeadImgURL        string `json:"HeadImgUrl"`
		Signature         string `json:"Signature"`
	} `json:"Profile"`
	ContinueFlag int     `json:"ContinueFlag"`
	SyncKey      SyncKey `json:"SyncKey"`
	SKey         string  `json:"SKey"`
	SyncCheckKey SyncKey `json:"SyncCheckKey"`
}

type GetContactResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []any        `json:"MemberList"`
	Seq          int          `json:"Seq"`
}