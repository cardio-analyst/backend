package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"

	domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/client"
	"github.com/cardio-analyst/backend/internal/gateway/ports/service"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
	"github.com/cardio-analyst/backend/pkg/model"
)

const (
	fontPath      = "/assets/font/"
	fontFile      = "helvetica_1251.json"
	fontFamily    = "Helvetica"
	reportUnicode = "cp1251"
)

var _ service.ReportService = (*PDFReportService)(nil)

type PDFReportService struct {
	recommendations service.RecommendationsService

	analyses        storage.AnalysisRepository
	basicIndicators storage.BasicIndicatorsRepository

	authClient client.Auth
}

func NewPDFReportService(
	recommendations service.RecommendationsService,
	analyses storage.AnalysisRepository,
	basicIndicators storage.BasicIndicatorsRepository,
	authClient client.Auth,
) *PDFReportService {
	return &PDFReportService{
		recommendations: recommendations,
		analyses:        analyses,
		basicIndicators: basicIndicators,
		authClient:      authClient,
	}
}

func (s *PDFReportService) GenerateReport(userID uint64) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", pwd+fontPath)
	pdf.AddFont(fontFamily, "", fontFile)
	pdf.SetFont(fontFamily, "", 16)
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(22)

	html := pdf.HTMLBasicNew()
	writeToHTML := writeToHTMLWithUnicode(reportUnicode, pdf, &html)

	title := "<h1>Отчёт о вашем сердечно-сосудистом здоровье</h1><br></br><br></br>"
	writeToHTML(title)

	var basicIndicatorsDataExists bool
	basicIndicatorsDataExists, err = s.fillBasicIndicatorsReportData(userID, pdf, writeToHTML)
	if err != nil {
		return "", err
	}

	var analysesDataExists bool
	analysesDataExists, err = s.fillAnalysesReportData(userID, pdf, writeToHTML)
	if err != nil {
		return "", err
	}

	var recommendationsDataExists bool
	recommendationsDataExists, err = s.fillRecommendationsReportData(userID, pdf, writeToHTML)
	if err != nil {
		return "", err
	}

	if !basicIndicatorsDataExists && !analysesDataExists && !recommendationsDataExists {
		pdf.Close()
		return "", domain.ErrNotEnoughDataToCompileReport
	}

	pdfFileName := fmt.Sprintf("%v.pdf", time.Now().UTC().Unix())

	if err = pdf.OutputFileAndClose(pdfFileName); err != nil {
		return "", err
	}

	return pdfFileName, nil
}

func (s *PDFReportService) fillBasicIndicatorsReportData(userID uint64, pdf *gofpdf.Fpdf, writeToHTML func(value string)) (bool, error) {
	basicIndicators, err := s.basicIndicators.FindAll(userID)
	if err != nil {
		return false, err
	}
	if len(basicIndicators) == 0 {
		return false, nil
	}

	user, err := s.authClient.GetUser(context.TODO(), model.UserCriteria{
		ID: userID,
	})
	if err != nil {
		return false, err
	}

	scoreData := domain.ExtractScoreDataFrom(basicIndicators)
	scoreData.Age = user.Age()

	weight, height, waistSize, bodyMassIndex := extractBMIIndications(basicIndicators)

	pdf.SetFontSize(16)

	generateRow("Возраст", strconv.Itoa(scoreData.Age), pdf, writeToHTML)

	generateRow("Пол", scoreData.Gender, pdf, writeToHTML)

	generateRow("Вес (кг)", fmt.Sprintf("%.1f", weight), pdf, writeToHTML)

	generateRow("Рост (см)", fmt.Sprintf("%.1f", height), pdf, writeToHTML)

	generateRow("Объем талии (см)", fmt.Sprintf("%.1f", waistSize), pdf, writeToHTML)

	generateRow("Индекс массы тела (ИМТ)", fmt.Sprintf("%.1f", bodyMassIndex), pdf, writeToHTML)

	var smokingStr string
	if scoreData.Smoking {
		smokingStr = "Да"
	} else {
		smokingStr = "Нет"
	}

	generateRow("Статус курения", smokingStr, pdf, writeToHTML)

	generateRow("Уровень систолического АД (мм.рт.ст.)", fmt.Sprintf("%.1f", scoreData.SBPLevel), pdf, writeToHTML)

	generateRow("Общий холестерин (ммоль/л)", fmt.Sprintf("%.1f", scoreData.TotalCholesterolLevel), pdf, writeToHTML)

	var (
		cvEventsRiskValue            int64
		idealCardiovascularAgesRange string
	)
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
		"Риск сердечно-сосудистых событий<br></br>в течение 10 лет по шкале SCORE", fmt.Sprintf("%d", cvEventsRiskValue)+"%",
		pdf, writeToHTML,
	)

	generateRow("Ваш «сердечно-сосудистый возраст»", idealCardiovascularAgesRange, pdf, writeToHTML)

	return true, nil
}

func (s *PDFReportService) fillAnalysesReportData(userID uint64, pdf *gofpdf.Fpdf, writeToHTML func(value string)) (bool, error) {
	analyses, err := s.analyses.FindAll(userID)
	if err != nil {
		return false, err
	}
	if len(analyses) == 0 {
		return false, nil
	}

	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(18)

	title := "<h2>Результаты лабораторных и инструментальных исследований</h2><br></br><br></br>"
	writeToHTML(title)

	for _, analysis := range analyses {
		pdf.SetTextColor(89, 89, 89)
		pdf.SetFontSize(16)

		date := analysis.CreatedAt.Format("02.01.2006 15:04")
		dateStr := fmt.Sprintf(`<h3>Дата: %s</h3><br></br><br></br>`, date)
		writeToHTML(dateStr)

		pdf.SetTextColor(0, 0, 0)

		if analysis.HighDensityCholesterol != nil {
			generateRow(
				"Холестерин высокой плотности (ЛПВП) (ммоль/л)", fmt.Sprintf("%.1f", *analysis.HighDensityCholesterol),
				pdf, writeToHTML,
			)
		}

		if analysis.LowDensityCholesterol != nil {
			generateRow(
				"Холестерин низкой плотности (ЛПНП) (ммоль/л)", fmt.Sprintf("%.1f", *analysis.LowDensityCholesterol),
				pdf, writeToHTML,
			)
		}

		if analysis.Triglycerides != nil {
			generateRow("Триглицериды (ммоль/л)", fmt.Sprintf("%.1f", *analysis.Triglycerides), pdf, writeToHTML)
		}

		if analysis.Lipoprotein != nil {
			generateRow("Липопротеин (г/л)", fmt.Sprintf("%.1f", *analysis.Lipoprotein), pdf, writeToHTML)
		}

		if analysis.HighlySensitiveCReactiveProtein != nil {
			generateRow(
				"Высокочувствительный С-реактивный белок (кардио) (мг/л)", fmt.Sprintf("%.1f", *analysis.HighlySensitiveCReactiveProtein),
				pdf, writeToHTML,
			)
		}

		if analysis.AtherogenicityCoefficient != nil {
			generateRow("Коэффициент атерогенности", fmt.Sprintf("%.1f", *analysis.AtherogenicityCoefficient), pdf, writeToHTML)
		}

		if analysis.Creatinine != nil {
			generateRow("Креатинин (ммоль/л)", fmt.Sprintf("%.1f", *analysis.Creatinine), pdf, writeToHTML)
		}

		if analysis.AtheroscleroticPlaquesPresence != nil {
			var result string
			if *analysis.AtheroscleroticPlaquesPresence {
				result = "Да"
			} else {
				result = "Нет"
			}

			generateRow(
				"Результаты УЗДМАГ. Наличие атеросклеротических бляшек", result,
				pdf, writeToHTML,
			)
		}

		writeToHTML("<br></br><br></br>")
	}

	return true, nil
}

func generateRow(label string, value string, pdf *gofpdf.Fpdf, htmlWhite func(value string)) {
	pdf.SetTextColor(0, 0, 0)
	htmlWhite(fmt.Sprintf(`<p>%s</p>: `, label))

	pdf.SetTextColor(89, 89, 89)
	htmlWhite(fmt.Sprintf(`<p>%s</p><br></br><br></br>`, value))
}

func (s *PDFReportService) fillRecommendationsReportData(userID uint64, pdf *gofpdf.Fpdf, htmlWrite func(value string)) (bool, error) {
	recommendations, err := s.recommendations.GetRecommendations(userID)
	if err != nil {
		return false, err
	}
	if len(recommendations) == 0 {
		return false, nil
	}

	pdf.AddPage()
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(18)

	title := "<h2>Рекомендации</h2><br></br><br></br>"
	htmlWrite(title)

	for _, recommendation := range recommendations {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFontSize(16)

		whatStr := fmt.Sprintf(`<h3>%s</h3><br></br><br></br>`, recommendation.What)
		htmlWrite(whatStr)

		pdf.SetFontSize(14)
		pdf.SetTextColor(89, 89, 89)

		whyStr := fmt.Sprintf(`<p>%s</p><br></br><br></br>`, recommendation.Why)
		htmlWrite(whyStr)

		pdf.SetTextColor(0, 0, 0)

		htmlWrite("<h3>Что нужно делать?<h3><br></br>")

		pdf.SetTextColor(89, 89, 89)

		howStr := fmt.Sprintf(`<p>%s</p><br></br><br></br>`, recommendation.How)
		htmlWrite(howStr)
	}

	return true, nil
}

func writeToHTMLWithUnicode(unicode string, pdf *gofpdf.Fpdf, html *gofpdf.HTMLBasicType) func(value string) {
	tr := pdf.UnicodeTranslatorFromDescriptor(unicode)

	_, lineHt := pdf.GetFontSize()

	return func(value string) {
		html.Write(lineHt, tr(value))
	}
}
