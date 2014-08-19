package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/smtc/JustCms/utils"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

const (
	ADMIN_ROUTE = "/admin/"
)

func AdminMux() *web.Mux {
	mux := web.New()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get(ADMIN_ROUTE, indexHandler)
	mux.Get(ADMIN_ROUTE+"menu", menuHandler)

	mux.Get(regexp.MustCompile(`^/admin/(?P<class>.+)\.(?P<model>.+)\.(?P<fn>.+)$`), templateHandler)

	// 这里有点乱，待重新规划
	mux.Get(ADMIN_ROUTE+"account/user/", AccountList)

	mux.NotFound(utils.NotFound)
	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	utils.RenderHtml("/admin/main.html", w, r)
}

/*
模板页暂时以 class.model_fn 分级
*/
func templateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	temp := fmt.Sprintf("/admin/%s/%s_%s.html", c.URLParams["class"], c.URLParams["model"], c.URLParams["fn"])
	utils.RenderHtml(temp, w, r)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("./admin/menu.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
