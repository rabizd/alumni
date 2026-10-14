package seed

import (
	"fmt"

	"github.com/rabizd/alumni/internal/model"
)

// Fixed samples: the same records on every machine, so the demo pages and
// requests.http always have something to show.
var (
	sampleUsers = []model.User{
		{Name: "Ayşe Yılmaz", Email: "ayse@example.com", Department: "Bilgisayar Mühendisliği",
			StartYear: 2017, GraduationYear: 2021, DoubleMajor: "Matematik", Advisor: "Prof. Dr. Ahmet Kaya", PrepExempt: true},
		{Name: "Mehmet Demir", Email: "mehmet@example.com", Department: "Hukuk",
			StartYear: 2015, GraduationYear: 2019, Minor: "İşletme"},
	}
	sampleAnnouncements = []model.Announcement{
		{Title: "Mezunlar buluşması", Author: "Mezunlar Ofisi",
			Body: "Bu yılki mezunlar buluşması 15 Kasım'da Beyazıt kampüsünde. Kayıt için profilinizin güncel olması yeterli."},
		{Title: "Kariyer günleri başvuruları açıldı", Author: "Kariyer Merkezi",
			Body: "Şirketinizle kariyer günlerine katılmak ve öğrencilerle tanışmak isteyen mezunlarımız 31 Ekim'e kadar başvurabilir."},
	}
	sampleAlumni = []model.Alumni{
		{Name: "Ayşe Yılmaz", Department: "Bilgisayar Mühendisliği", GraduationYear: 2021},
		{Name: "Mehmet Demir", Department: "Hukuk", GraduationYear: 2019},
	}
)

// Counts is how many random records the factories add on top of the samples.
type Counts struct {
	Users, Announcements, Alumni int
}

// Run stores the samples, then the random records. The samples go into empty
// tables only, so running the seeder twice does not add them twice. Every
// record goes through Validate and the normal store, exactly like one sent to
// the API.
func Run(f *Factory, c Counts) error {
	existingUsers, err := model.Users.List()
	if err != nil {
		return err
	}
	var users []model.User
	if len(existingUsers) == 0 {
		users = append(users, sampleUsers...)
	}
	for range c.Users {
		users = append(users, f.User())
	}
	for _, u := range users {
		if err := u.Validate(); err != nil {
			return fmt.Errorf("user %q: %w", u.Name, err)
		}
		if _, err := model.Users.Add(u); err != nil {
			return err
		}
	}

	existingAnnouncements, err := model.Announcements.List()
	if err != nil {
		return err
	}
	var announcements []model.Announcement
	if len(existingAnnouncements) == 0 {
		announcements = append(announcements, sampleAnnouncements...)
	}
	for range c.Announcements {
		announcements = append(announcements, f.Announcement())
	}
	for _, a := range announcements {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("announcement %q: %w", a.Title, err)
		}
		if _, err := model.Announcements.Add(a); err != nil {
			return err
		}
	}

	existingAlumni, err := model.Graduates.List()
	if err != nil {
		return err
	}
	var alumni []model.Alumni
	if len(existingAlumni) == 0 {
		alumni = append(alumni, sampleAlumni...)
	}
	for range c.Alumni {
		alumni = append(alumni, f.Alumni())
	}
	for _, a := range alumni {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("alumni %q: %w", a.Name, err)
		}
		if _, err := model.Graduates.Add(a); err != nil {
			return err
		}
	}
	return nil
}
