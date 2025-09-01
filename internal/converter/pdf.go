package converter

import (
	"fmt"
	"os"
	"path/filepath"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFConverter struct {
	options *wkhtmltopdf.PDFGenerator
}

func NewPDFConverter() *PDFConverter {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.MarginTop.Set(10)

	return &PDFConverter{
		options: pdfg,
	}
}

func (p *PDFConverter) ConvertHTMLToPDF(htmlContent []byte, outputPath string) error {
	// Create temp HTML file
	tempDir := os.TempDir()
	tempHTMLFile := filepath.Join(tempDir, fmt.Sprintf("pagemail_temp_%d.html", os.Getpid()))
	
	// Write HTML content to temp file
	if err := os.WriteFile(tempHTMLFile, htmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write temp HTML file: %w", err)
	}
	defer os.Remove(tempHTMLFile)

	// Add page to PDF generator
	page := wkhtmltopdf.NewPage(tempHTMLFile)
	
	// Set page options
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(9)
	page.Zoom.Set(0.95)
	
	p.options.AddPage(page)

	// Generate PDF
	if err := p.options.Create(); err != nil {
		return fmt.Errorf("failed to create PDF: %w", err)
	}

	// Write PDF to output file
	if err := p.options.WriteFile(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

func (p *PDFConverter) ConvertURLToPDF(url string, outputPath string) error {
	// Reset options for new conversion
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed to create PDF generator: %w", err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.MarginBottom.Set(10)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.MarginTop.Set(10)

	// Add page from URL
	page := wkhtmltopdf.NewPage(url)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(9)
	page.Zoom.Set(0.95)
	page.LoadErrorHandling.Set("ignore")
	page.LoadMediaErrorHandling.Set("ignore")
	
	pdfg.AddPage(page)

	// Generate PDF
	if err := pdfg.Create(); err != nil {
		return fmt.Errorf("failed to create PDF from URL: %w", err)
	}

	// Write PDF to output file
	if err := pdfg.WriteFile(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}