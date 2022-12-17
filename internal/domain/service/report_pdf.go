package service

import (
	"fmt"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/service"
	"github.com/cardio-analyst/backend/internal/ports/smtp"
	"github.com/cardio-analyst/backend/internal/ports/storage"
	"github.com/jung-kurt/gofpdf"
	"os"
	"strconv"
)

var _ service.ReportService = (*pdfReportService)(nil)

type pdfReportService struct {
	sender          smtp.Client
	recommendations service.RecommendationsService
	analyses        storage.AnalysisRepository
	basicIndicators storage.BasicIndicatorsRepository
	users           storage.UserRepository
}

func NewPDFReportService(
	sender smtp.Client,
	recommendations service.RecommendationsService,
	analyses storage.AnalysisRepository,
	basicIndicators storage.BasicIndicatorsRepository,
	users storage.UserRepository) *pdfReportService {
	return &pdfReportService{
		sender:          sender,
		recommendations: recommendations,
		analyses:        analyses,
		basicIndicators: basicIndicators,
		users:           users,
	}
}

func (s *pdfReportService) GenerateReport(userID uint64) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", pwd+"/font/")

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")

	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 16)

	html := pdf.HTMLBasicNew()

	htmlWhite := whiteToHTMLWithUnicode("cp1251", pdf, &html)

	title := "<h1>Отчёт о вашем сердечно-сосудистом здоровье</h1><br></br><br></br>"
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(22)

	htmlWhite(title)

	err = s.generateReportByBasicIndicators(userID, pdf, htmlWhite)

	if err != nil {
		return "", err
	}

	err = s.generateReportByAnalyses(userID, pdf, htmlWhite)

	if err != nil {
		return "", err
	}

	err = s.generateReportByRecommendations(userID, pdf, htmlWhite)

	if err != nil {
		return "", err
	}

	if err := pdf.OutputFileAndClose("report.pdf"); err != nil {
		return "", err
	}

	return "report.pdf", nil
}

func (s *pdfReportService) generateReportByBasicIndicators(userID uint64, pdf *gofpdf.Fpdf, htmlWhite func(value string)) error {
	basicIndicators, err := s.basicIndicators.FindAll(userID)

	if err != nil {
		return err
	}

	if len(basicIndicators) == 0 {
		return nil
	}

	pdf.SetFontSize(16)

	scoreData := models.ExtractScoreDataFrom(basicIndicators)

	criteria := models.UserCriteria{
		ID: &userID,
	}

	user, err := s.users.GetByCriteria(criteria)

	if err != nil {
		return err
	}

	scoreData.Age = user.Age()

	weight, height, waistSize, bodyMassIndex := GetBMIIndications(basicIndicators)

	generateRow(
		"Возвраст",
		strconv.FormatInt(int64(scoreData.Age), 10),
		pdf,
		htmlWhite,
	)

	generateRow(
		"Пол",
		scoreData.Gender,
		pdf,
		htmlWhite,
	)

	generateRow(
		"Вес",
		fmt.Sprintf("%.1f", weight),
		pdf,
		htmlWhite,
	)

	generateRow(
		"Рост",
		fmt.Sprintf("%.1f", height),
		pdf,
		htmlWhite,
	)

	generateRow(
		"Объем талии(см)",
		fmt.Sprintf("%.1f", waistSize),
		pdf,
		htmlWhite,
	)

	generateRow(
		"Индекс массы тела (ИМТ)",
		fmt.Sprintf("%.1f", bodyMassIndex),
		pdf,
		htmlWhite,
	)

	var smokingStr string

	if scoreData.Smoking {
		smokingStr = "Да"
	} else {
		smokingStr = "Нет"
	}

	generateRow(
		"Статус курения",
		smokingStr,
		pdf,
		htmlWhite,
	)

	generateRow(
		"Уровень систолического АД",
		fmt.Sprintf("%.1f", scoreData.SBPLevel),
		pdf,
		htmlWhite,
	)

	generateRow(
		"Общий холестерин",
		fmt.Sprintf("%.1f", scoreData.TotalCholesterolLevel),
		pdf,
		htmlWhite,
	)

	var cvEventsRiskValue int64
	var idealCardiovascularAgesRange string

	for _, indicators := range basicIndicators {
		if indicators.CVEventsRiskValue != nil && cvEventsRiskValue == 0 {
			cvEventsRiskValue = *indicators.CVEventsRiskValue
		}
		if indicators.IdealCardiovascularAgesRange != nil && idealCardiovascularAgesRange == "" {
			idealCardiovascularAgesRange = *indicators.IdealCardiovascularAgesRange
		}

		// fastest break condition
		if cvEventsRiskValue != 0 && idealCardiovascularAgesRange != "" {
			break
		}
	}

	generateRow(
		"Риск сердечно-сосудистых событий<br></br>в течение 10 лет по шкале SCORE",
		fmt.Sprintf("%d", cvEventsRiskValue)+"%",
		pdf,
		htmlWhite,
	)

	generateRow(
		"Ваш идеальный «сердечно-сосудистый возраст»",
		idealCardiovascularAgesRange,
		pdf,
		htmlWhite,
	)

	return nil
}

func (s *pdfReportService) generateReportByAnalyses(userID uint64, pdf *gofpdf.Fpdf, htmlWhite func(value string)) error {
	analyses, err := s.analyses.FindAll(userID)

	if err != nil {
		return err
	}

	if len(analyses) == 0 {
		return nil
	}

	pdf.AddPage()

	title := "<h2>Результаты лабораторных и инструментальных исследований</h2><br></br><br></br>"

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(18)
	htmlWhite(title)

	for _, analysis := range analyses {
		pdf.SetTextColor(89, 89, 89)

		var monthStr string

		var monthNumber = int64(analysis.CreatedAt.Month())

		if monthNumber > 10 {
			monthStr = strconv.FormatInt(monthNumber, 10)
		} else {
			monthStr = fmt.Sprintf(`0%d`, monthNumber)
		}

		date := fmt.Sprintf(`%d.%d.%s %d:%d`,
			analysis.CreatedAt.Day(),
			analysis.CreatedAt.Year(),
			monthStr,
			analysis.CreatedAt.Hour(),
			analysis.CreatedAt.Minute(),
		)
		dateStr := fmt.Sprintf(`<h3>Дата: %s</h3><br></br><br></br>`, date)
		pdf.SetFontSize(16)
		htmlWhite(dateStr)

		pdf.SetTextColor(0, 0, 0)

		if analysis.HighDensityCholesterol != nil {
			generateRow(
				"Холестерин высокой плотности (ЛПВП)",
				fmt.Sprintf("%.1f", *analysis.HighDensityCholesterol),
				pdf,
				htmlWhite,
			)
		}

		if analysis.LowDensityCholesterol != nil {
			generateRow(
				"Холестерин низкой плотности (ЛПНП)",
				fmt.Sprintf("%.1f", *analysis.LowDensityCholesterol),
				pdf,
				htmlWhite,
			)
		}

		if analysis.Triglycerides != nil {
			generateRow(
				"Триглицериды",
				fmt.Sprintf("%.1f", *analysis.Triglycerides),
				pdf,
				htmlWhite,
			)
		}

		if analysis.Lipoprotein != nil {
			generateRow(
				"Липопротеин",
				fmt.Sprintf("%.1f", *analysis.Lipoprotein),
				pdf,
				htmlWhite,
			)
		}

		if analysis.HighlySensitiveCReactiveProtein != nil {
			generateRow(
				"Высокочувствительный С-реактивный белок (кардио)",
				fmt.Sprintf("%.1f", *analysis.HighlySensitiveCReactiveProtein),
				pdf,
				htmlWhite,
			)
		}

		if analysis.AtherogenicityCoefficient != nil {
			generateRow(
				"Коэффициент атерогенности",
				fmt.Sprintf("%.1f", *analysis.AtherogenicityCoefficient),
				pdf,
				htmlWhite,
			)
		}

		if analysis.Creatinine != nil {
			generateRow(
				"Креатинин",
				fmt.Sprintf("%.1f", *analysis.Creatinine),
				pdf,
				htmlWhite,
			)
		}

		if analysis.AtheroscleroticPlaquesPresence != nil {

			var result string

			if *analysis.AtheroscleroticPlaquesPresence {
				result = "Да"
			} else {
				result = "Нет"
			}

			generateRow(
				"Результаты УЗДМАГ. Наличие атеросклеротических бляшек",
				result,
				pdf,
				htmlWhite,
			)
		}

		htmlWhite("<br></br><br></br>")
	}

	return nil
}

func generateRow(label string, value string, pdf *gofpdf.Fpdf, htmlWhite func(value string)) {
	pdf.SetTextColor(0, 0, 0)
	labelStr := fmt.Sprintf(`<p>%s</p>: `, label)
	htmlWhite(labelStr)

	pdf.SetTextColor(89, 89, 89)
	valueStr := fmt.Sprintf(`<p>%s</p><br></br><br></br>`, value)
	htmlWhite(valueStr)
}

func (s *pdfReportService) generateReportByRecommendations(userID uint64, pdf *gofpdf.Fpdf, htmlWhite func(value string)) error {
	recommendations, err := s.recommendations.GetRecommendations(userID)

	if err != nil {
		return err
	}

	if len(recommendations) == 0 {
		return nil
	}

	pdf.AddPage()

	title := "<h2>Рекомендации</h2><br></br><br></br>"

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(18)
	htmlWhite(title)

	for _, recommendation := range recommendations {
		pdf.SetTextColor(0, 0, 0)
		whatStr := fmt.Sprintf(`<h3>%s</h3><br></br><br></br>`, recommendation.What)
		pdf.SetFontSize(16)
		htmlWhite(whatStr)

		pdf.SetFontSize(14)
		pdf.SetTextColor(89, 89, 89)
		whyStr := fmt.Sprintf(`<p>%s</p><br></br><br></br>`, recommendation.Why)
		htmlWhite(whyStr)

		pdf.SetTextColor(0, 0, 0)
		htmlWhite("<h3>Что нужно делать?<h3><br></br>")

		pdf.SetTextColor(89, 89, 89)
		howStr := fmt.Sprintf(`<p>%s</p><br></br><br></br>`, recommendation.How)
		htmlWhite(howStr)
	}

	return nil
}

func whiteToHTMLWithUnicode(unicode string, pdf *gofpdf.Fpdf, html *gofpdf.HTMLBasicType) func(value string) {
	tr := pdf.UnicodeTranslatorFromDescriptor(unicode)

	_, lineHt := pdf.GetFontSize()

	return func(value string) {
		html.Write(lineHt, tr(value))
	}
}
