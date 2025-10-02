package utility

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertPDFToDocx converts PDF (text or scanned) to DOCX
func ConvertPDFToDocx(inputPath, outputDir, targetFormat string) (string, error) {
	return runLibreOfficeConversion(inputPath, outputDir, targetFormat)
}

// ConvertDocxToPDF converts DOCX to PDF
func ConvertDocxToPDF(inputPath, outputDir, targetFormat string) (string, error) {
	return runLibreOfficeConversion(inputPath, outputDir, targetFormat)
}

// ConvertDocToPDF converts DOC to PDF
func ConvertDocToPDF(inputPath, outputDir, targetFormat string) (string, error) {
	return runLibreOfficeConversion(inputPath, outputDir, targetFormat)
}

// ConvertPPTXToPDF converts PPTX to PDF
func ConvertPPTXToPDF(inputPath, outputDir, targetFormat string) (string, error) {
	return runLibreOfficeConversion(inputPath, outputDir, targetFormat)
}

// ConvertXLSXToPDF converts XLSX to PDF
func ConvertXLSXToPDF(inputPath, outputDir, targetFormat string) (string, error) {
	return runLibreOfficeConversion(inputPath, outputDir, targetFormat)
}

// Generic method using LibreOffice CLI
func runLibreOfficeConversion(inputPath, outputDir, targetFormat string) (string, error) {
	ext := strings.ToLower(filepath.Ext(inputPath))
	ext = strings.TrimPrefix(ext, ".")

	if ext == targetFormat {
		return "", errors.New("source and target formats are the same")
	}

	cmd := exec.Command(
		"libreoffice",
		"--headless",
		"--convert-to", targetFormat,
		"--outdir", outputDir,
		inputPath,
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	outputFile := filepath.Join(outputDir, getOutputFileName(inputPath, targetFormat))
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return "", errors.New("conversion failed, output file not found")
	}

	return outputFile, nil
}

// Helper: derive output filename
func getOutputFileName(inputPath, targetFormat string) string {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return base + "." + targetFormat
}
