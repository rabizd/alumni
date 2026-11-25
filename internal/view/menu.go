package view

import (
	"strings"

	"github.com/rabizd/alumni/internal/model"
)

// Guest is the Role of a menu item only visitors who are not signed in see.
const Guest = "guest"

// MenuItem is one link in the top menu.
type MenuItem struct {
	Label, Href string
	// Role is who sees the link: "" everyone, Guest only signed-out
	// visitors, otherwise at least that role (model.User.HasRole). It matches
	// the guard on the route in internal/routes/web.go, so the menu never
	// shows a link that would answer 403.
	Role string
	// Current is true on the link of the page being shown.
	Current bool
}

// menuItems is the whole menu, in order. Adding a page to the menu is one
// line here; layout.html and main.html both draw it.
var menuItems = []MenuItem{
	{Label: "Mezunlar", Href: "/graduates", Role: model.RoleGraduate},
	{Label: "İlanlar", Href: "/job-postings", Role: model.RoleGraduate},
	{Label: "Başvurular", Href: "/applications", Role: model.RoleGraduate},
	{Label: "Duyurular", Href: "/announcements"},
	{Label: "Kullanıcılar", Href: "/users", Role: model.RoleGraduate},
	{Label: "Hakkında", Href: "/about"},
	{Label: "Yönetim", Href: "/admin", Role: model.RoleAdmin},
	{Label: "Giriş", Href: "/login", Role: Guest},
	{Label: "Kayıt ol", Href: "/register", Role: Guest},
}

// menuFor returns the links the given user may see (nil means a guest),
// marking the one for path as current.
func menuFor(me *model.User, path string) []MenuItem {
	var out []MenuItem
	for _, item := range menuItems {
		switch {
		case item.Role == "":
		case item.Role == Guest:
			if me != nil {
				continue
			}
		case me == nil || !me.HasRole(item.Role):
			continue
		}
		item.Current = path == item.Href || strings.HasPrefix(path, item.Href+"/")
		out = append(out, item)
	}
	return out
}
