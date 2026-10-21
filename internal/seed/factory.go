// Package seed fills the database with data: factories make random but valid
// records, and Run stores them together with a few fixed samples.
package seed

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rabizd/alumni/internal/model"
)

var (
	firstNames  = []string{"Ayşe", "Mehmet", "Elif", "Can", "Zeynep", "Burak", "Selin", "Emre", "Deniz", "Kerem", "Ece", "Murat", "İrem", "Oğuz", "Şule"}
	lastNames   = []string{"Yılmaz", "Demir", "Kaya", "Şahin", "Çelik", "Arslan", "Doğan", "Kılıç", "Aydın", "Öztürk", "Koç", "Kurt"}
	departments = []string{"Bilgisayar Mühendisliği", "Hukuk", "Tıp", "İktisat", "Psikoloji", "Fizik", "Matematik", "İşletme", "Mimarlık", "Felsefe"}
	titles      = []string{"Mezunlar buluşması", "Kariyer günleri", "Mentorluk programı", "Staj ilanı", "Seminer duyurusu", "Burs başvuruları"}
	authors     = []string{"Mezunlar Ofisi", "Kariyer Merkezi", "Öğrenci İşleri", "Rektörlük"}
	jobTitles   = []string{"Junior Backend Developer", "Frontend Developer", "Veri Analisti", "Hukuk Stajyeri", "Asistan Hekim", "Finans Uzmanı", "Klinik Psikolog", "Mimari Tasarımcı", "Ürün Yöneticisi", "Araştırma Asistanı"}
	companies   = []string{"Trendyol", "Getir", "Aselsan", "Koç Holding", "Türk Telekom", "Garanti BBVA", "Acıbadem", "Turkcell", "Baykar", "Peak Games"}
	locations   = []string{"İstanbul", "Ankara", "İzmir", "Remote", ""}
)

// Factory makes random records. The same seed gives the same records, so a
// run can be repeated exactly.
type Factory struct {
	r *rand.Rand
	n int // counter that keeps the generated emails unique
}

func NewFactory(seed int64) *Factory {
	return &Factory{r: rand.New(rand.NewSource(seed))}
}

func (f *Factory) pick(list []string) string {
	return list[f.r.Intn(len(list))]
}

func (f *Factory) pickUser(users []model.User) model.User {
	return users[f.r.Intn(len(users))]
}

func (f *Factory) name() string {
	return f.pick(firstNames) + " " + f.pick(lastNames)
}

// User returns a user that passes User.Validate: started 4–20 years ago,
// graduated 4–6 years later, never in the future.
func (f *Factory) User() model.User {
	f.n++
	name := f.name()
	start := time.Now().Year() - 4 - f.r.Intn(17)
	grad := min(start+4+f.r.Intn(3), time.Now().Year())

	u := model.User{
		Name:           name,
		Email:          fmt.Sprintf("%s.%d@example.com", emailPart(name), f.n),
		Department:     f.pick(departments),
		StartYear:      start,
		GraduationYear: grad,
		PrepExempt:     f.r.Intn(2) == 0,
	}
	// the optional fields are filled only sometimes, like in real data
	if f.r.Intn(4) == 0 {
		u.DoubleMajor = f.pick(departments)
	}
	if f.r.Intn(3) == 0 {
		u.Minor = f.pick(departments)
	}
	if f.r.Intn(2) == 0 {
		u.Advisor = "Prof. Dr. " + f.name()
	}
	return u
}

func (f *Factory) Announcement() model.Announcement {
	return model.Announcement{
		Title:  f.pick(titles),
		Body:   fmt.Sprintf("%s bölümü mezunları için yeni bir duyuru. Ayrıntılar için Mezunlar Ofisi ile iletişime geçebilirsiniz.", f.pick(departments)),
		Author: f.pick(authors),
	}
}

func (f *Factory) Alumni() model.Alumni {
	return model.Alumni{
		Name:           f.name(),
		Department:     f.pick(departments),
		GraduationYear: time.Now().Year() - f.r.Intn(30),
	}
}

// JobPosting returns a posting by the given user that passes JobPosting.Validate.
// Some have an apply link or a deadline (1–12 weeks from today), like in real data.
func (f *Factory) JobPosting(userID int) model.JobPosting {
	p := model.JobPosting{
		UserID:         userID,
		Title:          f.pick(jobTitles),
		Company:        f.pick(companies),
		Location:       f.pick(locations),
		EmploymentType: f.pick(model.EmploymentTypes),
		Description:    fmt.Sprintf("%s mezunlarını ekibimizde görmek isteriz. Ayrıntılar için başvurun ya da ilanı açan mezuna yazın.", f.pick(departments)),
	}
	if f.r.Intn(2) == 0 {
		p.ApplyURL = fmt.Sprintf("https://example.com/jobs/%d", 1000+f.r.Intn(9000))
	}
	if f.r.Intn(2) == 0 {
		d := time.Now().AddDate(0, 0, 7*(1+f.r.Intn(12))).Format("2006-01-02")
		p.Deadline = &d
	}
	return p
}

// emailPart turns "Şule Öztürk" into "sule.ozturk".
func emailPart(name string) string {
	r := strings.NewReplacer("ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u", "İ", "i", " ", ".")
	return r.Replace(strings.ToLower(strings.ReplaceAll(name, "İ", "i")))
}
