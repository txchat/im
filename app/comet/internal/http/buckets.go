package http

import (
	"github.com/gin-gonic/gin"
)

type Bucket struct {
	UserNum  int
	GroupNum int
	Users    []string
	Groups   []string
}

type Numbers struct {
	BucketsNumber int `json:"buckets_number"`
	UsersNumber   int `json:"users_number"`
	GroupsNumber  int `json:"groups_number"`
}

func StaticsNumb(c *gin.Context) {
	usersCount := 0
	for _, bucket := range svcCtx.Buckets() {
		usersCount += bucket.ChannelCount()
	}
	groupsCount := 0
	for _, bucket := range svcCtx.Buckets() {
		groupsCount += bucket.GroupCount()
	}
	ret := Numbers{
		BucketsNumber: len(svcCtx.Buckets()),
		UsersNumber:   usersCount,
		GroupsNumber:  groupsCount,
	}
	c.JSON(200, ret)
}

func BucketsInfo(c *gin.Context) {
	buckets := make([]Bucket, 0)
	for _, bucket := range svcCtx.Buckets() {
		groups := make([]string, 0)

		keys := bucket.ChannelsCount()
		gps := bucket.GroupsCount()
		for gid := range gps {
			groups = append(groups, gid)
		}

		item := Bucket{
			UserNum:  len(keys),
			GroupNum: len(groups),
			Users:    keys,
			Groups:   groups,
		}
		buckets = append(buckets, item)
	}
	c.JSON(200, buckets)
}

type User struct {
	ChIndex int      `json:"ch_index"`
	Seq     int32    `json:"seq"`
	Key     string   `json:"key"`
	IP      string   `json:"ip"`
	Port    string   `json:"port"`
	Groups  []string `json:"groups"`
}

func UserInfo(c *gin.Context) {
	var arg struct {
		UID string `uri:"uid" binding:"required"`
	}
	if err := c.ShouldBindUri(&arg); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, getUser(arg.UID))
}

func getUser(uid string) *User {
	for i, bucket := range svcCtx.Buckets() {
		ch := bucket.Channel(uid)
		if ch != nil {
			return &User{
				ChIndex: i,
				Seq:     ch.GetSeq(),
				Key:     ch.GetKey(),
				IP:      ch.GetIP(),
				Port:    ch.GetPort(),
				Groups:  ch.GetGroups(),
			}
		}
	}
	return nil
}

type Group struct {
	GID          string            `json:"gid"`
	Members      map[string]string `json:"members"`
	MemberNumber int               `json:"member_number"`
}

func GroupsInfo(c *gin.Context) {
	c.JSON(200, getGroups())
}

func getGroups() []*Group {
	groups := make([]*Group, 0)
	for _, bucket := range svcCtx.Buckets() {
		gps := bucket.GroupsCount()
		for gid := range gps {
			groups = append(groups, getGroup(gid))
		}
	}
	return groups
}

func GroupInfo(c *gin.Context) {
	var arg struct {
		GID string `uri:"gid" binding:"required"`
	}
	if err := c.ShouldBindUri(&arg); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	c.JSON(200, getGroup(arg.GID))
}

func getGroup(gid string) *Group {
	members := make(map[string]string, 0)
	for _, bucket := range svcCtx.Buckets() {
		g := bucket.Group(gid)
		if g == nil {
			continue
		}
		mems, ips := g.Members()
		for i, mem := range mems {
			members[mem] = ips[i]
		}
	}
	return &Group{
		GID:     gid,
		Members: members,
	}
}
