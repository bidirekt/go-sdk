package catalogapi

import (
	"encoding/json"
	"net/http"
	"sync"

	"example.com/acme/shared/catalog"
	"github.com/bidirekt/go-sdk"
)

type Server struct {
	mutex    sync.RWMutex
	products map[string]catalog.Product
}

func NewServer() *Server {
	return &Server{products: map[string]catalog.Product{}}
}

func (this *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", this.listProducts)
	mux.HandleFunc("GET /products/{sku}", this.getProduct)
	mux.HandleFunc("PUT /products/{sku}", this.putProduct)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

var _ = bidirekt.Provides(http.MethodGet, "/products",
	bidirekt.Response(http.StatusOK, []catalog.Product{}, bidirekt.Named("ProductPage")),
)

func (this *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	page := make([]catalog.Product, 0, len(this.products))
	for _, product := range this.products {
		if category := r.URL.Query().Get("category"); category == "" || catalog.Category(category) == product.Category {
			page = append(page, product)
		}
	}
	writeJSON(w, http.StatusOK, page)
}

var _ = bidirekt.Provides(http.MethodGet, "/products/*",
	bidirekt.Response(http.StatusOK, catalog.Product{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

func (this *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	product, found := this.products[r.PathValue("sku")]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

var _ = bidirekt.Provides(http.MethodPut, "/products/*",
	bidirekt.Request(catalog.Product{}),
	bidirekt.Response(http.StatusOK, catalog.Product{}),
	bidirekt.Response(http.StatusBadRequest, nil),
)

func (this *Server) putProduct(w http.ResponseWriter, r *http.Request) {
	var product catalog.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.SKU = r.PathValue("sku")
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.products[product.SKU] = product
	writeJSON(w, http.StatusOK, product)
}
