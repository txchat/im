package http

import "github.com/gin-gonic/gin"

func Statics(c *gin.Context) {
	var arg struct {
		Drop bool `form:"isDrop"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	res := map[string]interface{}{
		"buckets": groupsInfo(arg.Drop),
	}
	result(c, res, OK)
}

func groupsInfo(isDrop bool) []map[string]interface{} {
	var res = make([]map[string]interface{}, len(srv.Buckets()))
	for i, bucket := range srv.Buckets() {
		if !isDrop {
			if bucket.GroupCount() == 0 {
				continue
			}
			if len(bucket.GroupsCount()) == 0 {
				continue
			}
		}
		item := map[string]interface{}{
			"counts":        bucket.GroupCount(),
			"group-members": bucket.GroupsCount(),
		}
		res[i] = item
	}
	return res
}

func GroupDetails(c *gin.Context) {
	var arg struct {
		Groups []string `form:"groups" binding:"required"`
	}
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}

	var groups = make([]interface{}, 0)
	for _, gid := range arg.Groups {
		gInfo := map[string]interface{}{
			"gid":     gid,
			"members": groupsMembers(gid),
		}
		groups = append(groups, gInfo)
	}
	res := map[string]interface{}{
		"groups": groups,
	}
	result(c, res, OK)
}

func groupsMembers(gid string) map[string]string {
	members := make(map[string]string, 0)
	for _, bucket := range srv.Buckets() {
		g := bucket.Group(gid)
		if g == nil {
			continue
		}
		mems, ips := g.Members()
		for i, mem := range mems {
			members[mem] = ips[i]
		}
	}
	return members
}
