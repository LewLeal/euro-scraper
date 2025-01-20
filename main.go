package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getEuroPrice se encarga de:
// 1. Hacer request a la página que muestra el valor del euro.
// 2. Parsear el HTML con goquery.
// 3. Devolver el valor como float64.
func getEuroPrice(url string) (float64, error) {
	// Realizamos la petición a la URL dada
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Usamos goquery para parsear el HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}

	// Aquí debemos buscar el elemento que contiene el valor del euro.
	// Según tu screenshot, está en <div class="col-sm-6 col-xs-6 centrado"> <strong>56,28013760</strong> </div>
	// Ajusta el selector CSS según corresponda. Ejemplo:
	selection := doc.Find("div.col-sm-6.col-xs-6.centrado strong")

	// Obtenemos el texto del <strong>
	euroText := strings.TrimSpace(selection.Text()) // "56,28013760" (ejemplo)
	if euroText == "" {
		return 0, fmt.Errorf("no se encontró el valor del euro en la página")
	}

	// Reemplazamos la coma por punto para poder convertir a float
	euroText = strings.ReplaceAll(euroText, ",", ".") // "56.28013760"

	// Convertimos a float64
	var euroValue float64
	_, err = fmt.Sscanf(euroText, "%f", &euroValue)
	if err != nil {
		return 0, fmt.Errorf("error convirtiendo el texto a float: %v", err)
	}

	return euroValue, nil
}

// euroHandler es un manejador HTTP que, al hacer GET, devuelve el precio del euro en JSON.
func euroHandler(w http.ResponseWriter, r *http.Request) {
	// Aquí pones la URL real de tu página
	urlBCV := "https://www.bcv.org.ve/"

	precio, err := getEuroPrice(urlBCV)
	if err != nil {
		// Si hay error, respondemos con un 500 y el mensaje
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respondemos con JSON simple
	w.Header().Set("Content-Type", "application/json")
	// Por ejemplo, {"euro_price": 56.28013760}
	fmt.Fprintf(w, `{"euro_price": %.8f}`, precio)
}

func main() {
	// Mapeamos la ruta "/euro" a nuestra función euroHandler
	http.HandleFunc("/euro", euroHandler)

	// Iniciamos el servidor en el puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
