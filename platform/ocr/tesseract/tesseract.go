package tesseract

import (
	"aidoc/internal/entity"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	BinaryPath  string
	Languages   string
	Timeout     time.Duration
	PageSegMode int
}

type Reader struct {
	cfg Config
}

func New(cfg Config) (*Reader, error) {
	if cfg.Languages == "" {
		cfg.Languages = "eng"
	}
	if cfg.PageSegMode == 0 {
		cfg.PageSegMode = 6
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60 * time.Second
	}
	return &Reader{cfg: cfg}, nil
}

func (r *Reader) AnalyzeExpense(ctx context.Context, input *entity.ReadInput) (*entity.OCRDocument, error) {
	doc := &entity.OCRDocument{
		ID:          uuid.NewString(),
		ProcessedAt: time.Now(),
		Status:      "success",
	}

	bin := r.cfg.BinaryPath
	if bin == "" {
		bin = "tesseract"
	}
	if _, err := exec.LookPath(bin); err != nil {
		doc.Status = "failed"
		doc.ErrorMessage = fmt.Sprintf("tesseract binary not found: %v", err)
		return doc, fmt.Errorf("%s", doc.ErrorMessage)
	}

	ext, err := detectFormat(input.Bytes)
	if err != nil {
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		return doc, err
	}

	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	imageFiles, workDir, isPDF, err := prepareInput(ctx, input.Bytes, ext)
	if err != nil {
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		return doc, err
	}
	defer func() {
		for _, f := range imageFiles {
			os.Remove(f)
		}
		if workDir != "" {
			os.RemoveAll(workDir)
		}
	}()

	var texts []string
	var totalConf float64
	var wordCount int
	for _, img := range imageFiles {
		raw, conf, err := r.ocrImage(ctx, img)
		if err != nil {
			doc.Status = "failed"
			doc.ErrorMessage = err.Error()
			return doc, err
		}
		if raw != "" {
			texts = append(texts, raw)
		}
		totalConf += conf
		wordCount++
	}

	raw := strings.Join(texts, "\n\n")
	avgConf := 0.0
	if wordCount > 0 {
		avgConf = totalConf / float64(wordCount)
	}

	doc.RawText = raw
	doc.Confidence = avgConf
	if isPDF {
		doc.PageCount = len(imageFiles)
	} else {
		doc.PageCount = 1
	}
	return doc, nil
}

func (r *Reader) ocrImage(ctx context.Context, imgPath string) (string, float64, error) {
	outBase := imgPath + "-out"
	defer func() {
		for _, e := range []string{".txt", ".tsv"} {
			os.Remove(outBase + e)
		}
	}()

	cmd := exec.CommandContext(ctx, r.binary(),
		imgPath, outBase,
		"-l", r.cfg.Languages,
		"--psm", strconv.Itoa(r.cfg.PageSegMode),
		"tsv",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", 0, fmt.Errorf("tesseract failed: %v: %s", err, string(out))
	}
	return parseTSV(outBase + ".tsv")
}

func (r *Reader) binary() string {
	if r.cfg.BinaryPath != "" {
		return r.cfg.BinaryPath
	}
	return "tesseract"
}

func prepareInput(ctx context.Context, data []byte, ext string) ([]string, string, bool, error) {
	if ext != ".pdf" {
		inFile, err := os.CreateTemp("", "tess-in-*"+ext)
		if err != nil {
			return nil, "", false, err
		}
		defer inFile.Close()
		if _, err := inFile.Write(data); err != nil {
			return nil, "", false, err
		}
		return []string{inFile.Name()}, "", false, nil
	}

	pdfFile, err := os.CreateTemp("", "tess-in-*.pdf")
	if err != nil {
		return nil, "", false, err
	}
	defer pdfFile.Close()
	if _, err := pdfFile.Write(data); err != nil {
		return nil, "", false, err
	}
	pdfFile.Close()

	dir, err := os.MkdirTemp("", "tess-pdf-")
	if err != nil {
		return nil, "", false, err
	}
	pagePattern := filepath.Join(dir, "page")
	cmd := exec.CommandContext(ctx, "pdftoppm", "-png", "-r", "200", pdfFile.Name(), pagePattern)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(dir)
		return nil, "", false, fmt.Errorf("pdftoppm failed: %v: %s", err, string(out))
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		os.RemoveAll(dir)
		return nil, "", false, err
	}
	var images []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".png") {
			images = append(images, filepath.Join(dir, e.Name()))
		}
	}
	if len(images) == 0 {
		os.RemoveAll(dir)
		return nil, "", false, fmt.Errorf("pdftoppm produced no images")
	}
	return images, dir, true, nil
}

func detectFormat(data []byte) (string, error) {
	if len(data) < 4 {
		return "", fmt.Errorf("input too small to detect format")
	}
	switch {
	case data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47:
		return ".png", nil
	case data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
		return ".jpg", nil
	case (data[0] == 0x49 && data[1] == 0x49 && data[2] == 0x2A && data[3] == 0x00) ||
		(data[0] == 0x4D && data[1] == 0x4D && data[2] == 0x00 && data[3] == 0x2A):
		return ".tiff", nil
	case data[0] == 0x25 && data[1] == 0x50 && data[2] == 0x44 && data[3] == 0x46:
		return ".pdf", nil
	}
	return "", fmt.Errorf("unsupported input format for tesseract")
}

func parseTSV(path string) (string, float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, fmt.Errorf("failed reading tesseract tsv: %w", err)
	}

	var lines []string
	var totalConf float64
	var count int

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) < 12 {
			continue
		}
		level := fields[0]
		if level != "5" {
			continue
		}
		confStr := fields[10]
		conf, err := strconv.ParseFloat(confStr, 64)
		if err != nil {
			continue
		}
		if conf < 0 {
			continue
		}
		totalConf += conf
		count++

		text := strings.TrimSpace(fields[11])
		if text != "" {
			lines = append(lines, text)
		}
	}

	avg := 0.0
	if count > 0 {
		avg = totalConf / float64(count)
	}
	return strings.Join(lines, "\n"), avg, nil
}
