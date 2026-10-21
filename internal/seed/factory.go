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
	locations   = []string{"İstanbul", "Ankara", "İzmir", "Remote", ""}
)

// jobsByDepartment lists, for each department, positions and companies that
// fit it: a graduate posts openings in their own field.
var jobsByDepartment = map[string]struct{ titles, companies []string }{
	"Bilgisayar Mühendisliği": {[]string{"Junior Backend Developer", "Frontend Developer", "Gömülü Yazılım Mühendisi", "DevOps Mühendisi"}, []string{"Trendyol", "Getir", "Baykar", "Aselsan", "Peak Games"}},
	"Hukuk":                   {[]string{"Hukuk Stajyeri", "Avukat", "Hukuk Müşaviri"}, []string{"Demir Hukuk Bürosu", "Koç Holding", "Garanti BBVA"}},
	"Tıp":                     {[]string{"Asistan Hekim", "Pratisyen Hekim", "Klinik Araştırma Uzmanı"}, []string{"Acıbadem", "Memorial", "Medipol Hastanesi"}},
	"İktisat":                 {[]string{"Ekonomist", "Veri Analisti", "Risk Analisti"}, []string{"Garanti BBVA", "Türkiye İş Bankası", "Koç Holding"}},
	"Psikoloji":               {[]string{"Klinik Psikolog", "Psikolojik Danışman", "Nöropsikolog"}, []string{"Acıbadem", "Memorial", "Medipol Hastanesi"}},
	"Fizik":                   {[]string{"Araştırma Asistanı", "Ar-Ge Mühendisi", "Veri Analisti"}, []string{"Aselsan", "TÜBİTAK", "Baykar"}},
	"Matematik":               {[]string{"Veri Analisti", "Aktüer", "Kantitatif Analist"}, []string{"Garanti BBVA", "Anadolu Sigorta", "Turkcell"}},
	"İşletme":                 {[]string{"Finans Uzmanı", "Ürün Yöneticisi", "Pazarlama Uzmanı"}, []string{"Koç Holding", "Türk Telekom", "Getir"}},
	"Mimarlık":                {[]string{"Mimari Tasarımcı", "Şantiye Mimarı", "İç Mimar"}, []string{"Rönesans Holding", "Limak", "Tabanlıoğlu Mimarlık"}},
	"Felsefe":                 {[]string{"Editör", "Redaktör", "İçerik Yazarı"}, []string{"İletişim Yayınları", "Can Yayınları", "Yapı Kredi Yayınları"}},
}

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
// The position and company fit the poster's department; a department not in
// jobsByDepartment (one typed in through the API) gets a computer engineering job.
// Some have an apply link or a deadline (1–12 weeks from today), like in real data.
func (f *Factory) JobPosting(poster model.User) model.JobPosting {
	jobs, ok := jobsByDepartment[poster.Department]
	if !ok {
		jobs = jobsByDepartment["Bilgisayar Mühendisliği"]
	}
	p := model.JobPosting{
		UserID:         poster.ID,
		Title:          f.pick(jobs.titles),
		Company:        f.pick(jobs.companies),
		Location:       f.pick(locations),
		EmploymentType: f.pick(model.EmploymentTypes),
		Description:    fmt.Sprintf("Ekibimize %s mezunu bir arkadaş arıyoruz. Ayrıntılar için başvurun ya da bana yazın.", poster.Department),
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
