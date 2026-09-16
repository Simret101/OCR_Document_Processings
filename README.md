````markdown
# OCR Document Processing

A research-driven **document processing and OCR system implemented in Go**, covering the complete document-processing pipeline from image preprocessing and compression to segmentation, OCR, and postprocessing.

The project was developed through extensive research and experimentation into different techniques for improving OCR results on real-world documents. It combines image-processing algorithms, document segmentation, OCR engines, text correction, normalization, and reconstruction into a modular processing pipeline.

The target of the evaluated OCR pipeline is **95%+ OCR passing accuracy** on the documents used for testing.

---

## Overview

Optical Character Recognition is not only an OCR-engine problem.

The quality of the final extracted text depends on several stages of processing:

- Input image quality
- Image preprocessing
- Noise and artifact removal
- Image normalization
- Thresholding and binarization
- Document segmentation
- OCR engine configuration
- Character recognition
- Text correction
- Encoding normalization
- Paragraph reconstruction
- Document-specific postprocessing

This project therefore treats OCR as a complete pipeline rather than simply sending an image to an OCR engine.

```text
                         Document
                            │
                            ▼
                  ┌───────────────────┐
                  │   Preprocessing   │
                  │                   │
                  │ Image Processing  │
                  │ Filtering         │
                  │ Grayscale         │
                  │ Contrast          │
                  │ Normalization     │
                  │ Line Removal      │
                  │ Localization      │
                  │ Compression       │
                  │    └─ Huffman     │
                  └─────────┬─────────┘
                            │
                            ▼
                  ┌───────────────────┐
                  │    Segmentation   │
                  │                   │
                  │ Otsu              │
                  │ Adaptive          │
                  │ Sauvola           │
                  │ Components        │
                  │ Lines             │
                  │ Regions           │
                  │ Morphology        │
                  └─────────┬─────────┘
                            │
                            ▼
                  ┌───────────────────┐
                  │       OCR         │
                  │                   │
                  │ Tesseract         │
                  │ Amazon Textract   │
                  └─────────┬─────────┘
                            │
                            ▼
                  ┌───────────────────┐
                  │   Postprocessing  │
                  │                   │
                  │ Character Fixing  │
                  │ Correction        │
                  │ Filtering         │
                  │ Normalization     │
                  │ Numeric Handling  │
                  │ Reconstruction    │
                  │ Header/Footer     │
                  └─────────┬─────────┘
                            │
                            ▼
                    Final OCR Output
````

---

# Project Goals

The main goal of this project is to research and implement techniques that improve the reliability of OCR when processing real-world documents.

The project focuses on:

* Improving document image quality before OCR
* Investigating different segmentation techniques
* Supporting multiple OCR engines
* Reducing OCR recognition errors
* Correcting extracted text
* Reconstructing document content
* Normalizing OCR output
* Exploring image compression techniques
* Building reusable document-processing components
* Measuring OCR results instead of relying only on visual inspection
* Maintaining a clean and modular architecture

---

# Processing Pipeline

The project is divided into four major stages:

```text
┌─────────────────┐
│  Preprocessing  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Segmentation  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│      OCR        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Postprocessing │
└─────────────────┘
```

Each stage has its own responsibilities and can be developed, tested, and improved independently.

---

# 1. Preprocessing

The preprocessing package contains the image-processing techniques used to prepare documents before segmentation and OCR.

The purpose of preprocessing is to transform a raw document into a cleaner and more OCR-friendly representation.

The repository contains techniques for:

* Image reading
* Grayscale conversion
* Image filtering
* Noise reduction
* Contrast enhancement
* Image normalization
* Line removal
* Localization augmentation
* Color-space transformation
* Thresholding
* Image compression

```text
Raw Document
     │
     ▼
Image Reading
     │
     ▼
Grayscale / Color Processing
     │
     ▼
Filtering
     │
     ▼
Noise Reduction
     │
     ▼
Contrast Enhancement
     │
     ▼
Line Removal
     │
     ▼
Normalization
     │
     ▼
OCR-ready Image
```

## Compression

Preprocessing also contains a custom image-compression implementation.

The compression component explores **Huffman coding** and the lower-level algorithms required to encode and decode image data.

It includes:

* Bit reader
* Bit writer
* Huffman tree
* Huffman code generation
* Priority queue
* Encoding
* Decoding
* Differential encoding
* YCoCg color-space conversion

```text
Image Data
    │
    ▼
Color/Data Transformation
    │
    ▼
Differential Encoding
    │
    ▼
Frequency Analysis
    │
    ▼
Huffman Tree
    │
    ▼
Huffman Code Generation
    │
    ▼
Bit-level Encoding
    │
    ▼
Compressed Data
```

The decompression process reverses the transformation:

```text
Compressed Data
      │
      ▼
Bit Reader
      │
      ▼
Huffman Decoding
      │
      ▼
Differential Decoding
      │
      ▼
Data Reconstruction
      │
      ▼
Reconstructed Image
```

This part of the project provides a lower-level implementation of compression concepts rather than treating compression as a black-box operation.

---

# 2. Segmentation

Segmentation focuses on separating useful information within a document and preparing regions of the image for OCR.

The project researches several segmentation and binarization techniques because different document conditions can require different approaches.

## Otsu Thresholding

The project implements Otsu thresholding using GoCV.

Otsu determines a global threshold that separates foreground and background based on the image histogram.

It is particularly useful for documents with relatively consistent illumination and background characteristics.

---

# 3. OCR

The project supports and researches multiple OCR approaches.

## Tesseract

[Tesseract OCR](https://github.com/tesseract-ocr/tesseract) is used for local OCR processing.

The project investigates how different combinations of:

* preprocessing
* thresholding
* segmentation
* image quality
* OCR configuration
* postprocessing

affect the final OCR output.

The objective is to optimize the complete pipeline rather than relying solely on the default OCR configuration.

---

## Amazon Textract

The project also contains an `amazontextract` package for working with Amazon Textract.

```text
amazontextract/
├── async/
├── sync/
├── amazonteextract.go
├── extract.go
├── helper.go
├── populate.go
└── textract.init.go
```

The implementation includes both synchronous and asynchronous processing approaches.

Amazon Textract provides another approach for extracting text and document information, allowing the project to investigate document processing using both local OCR and cloud-based document analysis.

```text
                    Document
                       │
              ┌────────┴────────┐
              │                 │
              ▼                 ▼
         Tesseract          Textract
              │                 │
              └────────┬────────┘
                       ▼
                 OCR Output
                       │
                       ▼
                 Postprocessing
```

---

# 4. Postprocessing

OCR output frequently contains errors even when the majority of characters have been recognized correctly.

The postprocessing package is responsible for cleaning, correcting, normalizing, and reconstructing the extracted content.

It contains components for:

* Character correction
* General text correction
* Filtering
* Encoding normalization
* Whitespace normalization
* Numeric processing
* Paragraph reconstruction
* Regex-based correction rules
* Header/footer removal

```text
OCR Output
    │
    ▼
Character Correction
    │
    ▼
Text Correction
    │
    ▼
Filtering
    │
    ▼
Numeric Processing
    │
    ▼
Encoding Normalization
    │
    ▼
Whitespace Normalization
    │
    ▼
Paragraph Reconstruction
    │
    ▼
Header/Footer Removal
    │
    ▼
Final Text
```

---

## Character Correction

OCR engines can confuse visually similar characters.

The character-correction layer provides rules for handling known recognition errors.

---

## Numeric Processing

Numbers can be particularly sensitive to OCR errors.

The numeric processing layer handles document-specific recognition problems involving numeric content.

---

## Encoding Normalization

OCR output can contain inconsistent or unexpected character encoding.

The encoding normalization layer converts the extracted content into a consistent representation.

---

## Whitespace Normalization

OCR can produce:

* Multiple spaces
* Incorrect line breaks
* Empty lines
* Inconsistent spacing

The whitespace normalization stage cleans these artifacts.

---

## Paragraph Reconstruction

OCR engines can split text incorrectly across lines or regions.

Paragraph reconstruction attempts to rebuild the logical structure of extracted text.

---

## Header and Footer Removal

Documents often contain recurring headers and footers that should not be treated as part of the main document content.

The pipeline includes processing specifically for identifying and removing these sections when required.

---

# Research Approach

The project was developed through research, implementation, experimentation, and evaluation.

The general workflow is:

```text
Research
   │
   ▼
Identify OCR Problem
   │
   ▼
Research Possible Techniques
   │
   ▼
Implement Technique
   │
   ▼
Run OCR
   │
   ▼
Analyze Result
   │
   ▼
Postprocess Output
   │
   ▼
Measure Accuracy
   │
   ▼
Iterate
```

The research covers the entire OCR lifecycle instead of focusing on only one component.

---

# OCR Accuracy

The project was developed with a target of achieving **95%+ OCR passing accuracy** on the documents evaluated during development.

The approach is based on improving the complete processing pipeline:

```text
                    OCR Quality
                        │
       ┌────────────────┼────────────────┐
       │                │                │
       ▼                ▼                ▼
 Preprocessing      Segmentation        OCR
       │                │                │
       └────────────────┼────────────────┘
                        │
                        ▼
                 Postprocessing
                        │
                        ▼
                  Final Result
```

The measured result depends on:

* The document dataset
* Document quality
* Document type
* Evaluation criteria
* OCR engine
* Preprocessing configuration
* Segmentation strategy
* Postprocessing rules

Therefore, the 95%+ result refers to the documents and evaluation methodology used during testing rather than representing a universal OCR accuracy guarantee.

---

# Architecture

The project is implemented entirely around **Go-based modular components** and follows **Clean Architecture principles**.

The architecture separates document-processing logic from infrastructure and external OCR providers.

```text
                 ┌─────────────────────┐
                 │     Application     │
                 │       Use Cases     │
                 └──────────┬──────────┘
                            │
                            ▼
                 ┌─────────────────────┐
                 │       Domain       │
                 │ Document Processing │
                 │     Interfaces     │
                 └──────────┬──────────┘
                            │
             ┌──────────────┼──────────────┐
             │              │              │
             ▼              ▼              ▼
       Preprocessing    OCR Services   Postprocessing
             │              │              │
             │        ┌─────┴─────┐        │
             │        ▼           ▼        │
             │   Tesseract   Textract      │
             │                             │
             └──────────────┬──────────────┘
                            │
                            ▼
                       Final Output
```

The architecture is designed around:

* Separation of concerns
* Dependency inversion
* Interfaces
* Independent processing stages
* Testability
* Reusable components
* Replaceable implementations
* Maintainability

---

# Technology Stack

| Technology              | Purpose                              |
| ----------------------- | ------------------------------------ |
| **Go**                  | Primary programming language         |
| **GoCV**                | Computer vision and image processing |
| **OpenCV**              | Image-processing algorithms          |
| **Tesseract**           | Local OCR                            |
| **Amazon Textract**     | Cloud-based document/OCR processing  |
| **Huffman Coding**      | Image/data compression               |
| **Clean Architecture**  | Application architecture             |
| **Regular Expressions** | OCR correction and filtering         |

---

# Key Areas of Research

The project covers several areas of document processing and OCR:

### Image Processing

* Image filtering
* Grayscale conversion
* Noise reduction
* Contrast enhancement
* Image normalization
* Line removal
* Color-space conversion
* Localization

### Compression

* Huffman coding
* Huffman tree construction
* Code generation
* Priority queues
* Bit-level encoding
* Bit-level decoding
* Differential encoding
* YCoCg conversion

### Segmentation

* Otsu thresholding
* Adaptive thresholding
* Sauvola thresholding
* Connected components
* Projection profiles
* Text-line segmentation
* Region segmentation
* Contour detection
* Morphological processing

### OCR

* Tesseract
* Amazon Textract
* OCR input optimization
* OCR configuration
* OCR output evaluation

### Postprocessing

* Character correction
* Numeric correction
* Text filtering
* Encoding normalization
* Whitespace normalization
* Paragraph reconstruction
* Header/footer removal
* Regex-based correction

---

# Why Go?

Go was chosen as the implementation language because it provides a good foundation for building efficient and maintainable document-processing systems and it is fast.

The project benefits from Go's:

* Strong type system
* Simple concurrency model
* Efficient execution
* Low-level control when required
* Straightforward package structure
* Small deployment footprint
* Interface-based design
* Good support for building processing pipelines

The language also makes it possible to implement lower-level algorithms, such as the Huffman compression components, while keeping the higher-level OCR pipeline modular.

---

# Design Principles

## Modular Processing

Each processing stage has a defined responsibility.

```text
Preprocessing
      ↓
Segmentation
      ↓
OCR
      ↓
Postprocessing
```

Individual algorithms can be tested and improved without rewriting the entire pipeline.

## Research-Driven Implementation

Algorithms are implemented based on research into document processing and OCR rather than treating the OCR engine as a black box.

## Replaceable OCR Providers

The architecture allows OCR implementations to be separated from the rest of the processing pipeline.

This makes it possible to work with different OCR approaches such as Tesseract and Amazon Textract.

## Clean Architecture

Domain and processing logic remain separated from external infrastructure and providers.

This keeps the project easier to test, maintain, and extend.

---

# Example Processing Strategy

A typical document-processing workflow can look like:

```text
Input Document
      │
      ▼
Image Reading
      │
      ▼
Grayscale Conversion
      │
      ▼
Noise Reduction
      │
      ▼
Contrast Enhancement
      │
      ▼
Thresholding / Binarization
      │
      ▼
Document Segmentation
      │
      ▼
Text Regions
      │
      ▼
Tesseract / Textract
      │
      ▼
Character Correction
      │
      ▼
Filtering
      │
      ▼
Encoding Normalization
      │
      ▼
Whitespace Normalization
      │
      ▼
Paragraph Reconstruction
      │
      ▼
Final Structured Text
```

The exact processing sequence can be adapted depending on the characteristics of the document being processed.

---

# Project Status

This repository represents an ongoing collection of **OCR research, document-processing algorithms, experiments, and production-oriented implementations**.

It is intentionally broader than an OCR wrapper.

The project explores how:

```text
Image Processing
       +
Compression
       +
Segmentation
       +
OCR
       +
Postprocessing
       =
Document Processing Pipeline
```

can be combined to produce more reliable OCR results.

---

# Future Work

Potential areas for further development include:

* Automated preprocessing selection
* Document-type classification
* Advanced document-layout analysis
* Improved table detection
* More OCR providers
* OCR confidence analysis
* Automated benchmarking
* Larger evaluation datasets
* Parallel document processing
* GPU-accelerated image processing
* More advanced document binarization techniques
* Improved structured-data extraction

---

# Author

**Simret Belete**

Software Engineer focused on Go, backend systems, document processing, and system design.

---
