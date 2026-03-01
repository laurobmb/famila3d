package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PageData struct {
	Input      InputData
	Result     ResultData
	Calculated bool
}

type InputData struct {
	FilamentPrice float64
	Weight        float64
	Hours         float64
	LightPrice    float64
	FailureMargin float64
	PrinterMargin float64
	Markup        float64
}

type ResultData struct {
	FilamentCost float64
	LightCost    float64
	Cost1        float64
	Cost2        float64
	SalePrice    float64
}

// responseRecorder intercepta o ResponseWriter para capturar o HTTP Status Code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Intercepta o Write para garantir que captures status 200 implícitos do Go
func (rec *responseRecorder) Write(b []byte) (int, error) {
	if rec.statusCode == 0 {
		rec.statusCode = http.StatusOK
	}
	return rec.ResponseWriter.Write(b)
}

// loggingMiddleware registra os detalhes de cada requisição HTTP
func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Inicializa com 200 OK por padrão
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logger.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status", rec.statusCode),
			slog.Duration("duration", duration),
		)
	})
}

func main() {
	// Configura o logger para formato JSON no stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	funcMap := template.FuncMap{
		"sub": func(a, b float64) float64 {
			return a - b
		},
	}

	mux := http.NewServeMux()

	// Rota de arquivos estáticos
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Rota principal
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Input: InputData{
				FilamentPrice: 130.0,
				LightPrice:    1.0,
				FailureMargin: 0.1,
				PrinterMargin: 0.1,
				Markup:        2.0,
			},
		}

		if r.Method == http.MethodPost {
			// Limite de 1MB para prevenir OOM
			r.Body = http.MaxBytesReader(w, r.Body, 1048576)

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Requisição inválida ou muito grande", http.StatusBadRequest)
				return
			}

			// Verifica se algum campo excedeu 10 caracteres
			for _, val := range r.PostForm {
				if len(val[0]) > 10 {
					http.Error(w, "Valor de input excede o limite permitido (10 caracteres)", http.StatusBadRequest)
					return
				}
			}
			
			// Função auxiliar para converter e validar os inputs numéricos de forma segura
			parseFloatSafe := func(key string) (float64, bool) {
				valStr := r.FormValue(key)
				if valStr == "" {
					return 0, true // Campo vazio, ignora se não for obrigatório na sua regra
				}
				valFloat, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					return 0, false
				}
				return valFloat, true
			}

			var ok bool
			if data.Input.FilamentPrice, ok = parseFloatSafe("filament_price"); !ok {
				http.Error(w, "Preço do filamento inválido", http.StatusBadRequest)
				return
			}
			if data.Input.Weight, ok = parseFloatSafe("weight"); !ok {
				http.Error(w, "Peso inválido", http.StatusBadRequest)
				return
			}
			if data.Input.Hours, ok = parseFloatSafe("hours"); !ok {
				http.Error(w, "Horas inválidas", http.StatusBadRequest)
				return
			}
			if data.Input.LightPrice, ok = parseFloatSafe("light_price"); !ok {
				http.Error(w, "Preço da luz inválido", http.StatusBadRequest)
				return
			}
			if data.Input.Markup, ok = parseFloatSafe("markup"); !ok {
				http.Error(w, "Markup inválido", http.StatusBadRequest)
				return
			}

			fCost := (data.Input.FilamentPrice / 1000) * data.Input.Weight
			lCost := data.Input.Hours * data.Input.LightPrice
			c1 := fCost + lCost
			c2 := c1 * (1 + data.Input.FailureMargin + data.Input.PrinterMargin)
			sale := c2 * data.Input.Markup

			data.Result = ResultData{
				FilamentCost: fCost,
				LightCost:    lCost,
				Cost1:        c1,
				Cost2:        c2,
				SalePrice:    sale,
			}
			data.Calculated = true

			// Log estruturado de negócio
			slog.Info("Orçamento calculado",
				slog.Float64("peso_g", data.Input.Weight),
				slog.Float64("tempo_h", data.Input.Hours),
				slog.Float64("valor_venda", sale),
			)
		}

		tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
		if err != nil {
			slog.Error("Falha ao processar template", slog.String("erro", err.Error()))
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			slog.Error("Falha ao renderizar template", slog.String("erro", err.Error()))
		}
	})

	loggedMux := loggingMiddleware(logger, mux)

	porta := "8080"
	slog.Info("Iniciando servidor web", slog.String("porta", porta))
	
	if err := http.ListenAndServe(":"+porta, loggedMux); err != nil {
		slog.Error("Servidor encerrado com erro", slog.String("erro", err.Error()))
		os.Exit(1)
	}
}