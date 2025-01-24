package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ovaphlow.com/crate/data/schema"
	"ovaphlow.com/crate/data/service"
	"ovaphlow.com/crate/data/utility"
)

func LoadMySQLRouter(mux *http.ServeMux, prefix string, service *service.ApplicationServiceImpl) {
	route := &RouteMySQL{service: service}

	mux.HandleFunc("DELETE "+prefix+"/mysql/{st}/{id}", func(w http.ResponseWriter, r *http.Request) {
		route.delete(w, r)
	})

	mux.HandleFunc("PUT "+prefix+"/mysql/{st}/{id}", func(w http.ResponseWriter, r *http.Request) {
		route.put(w, r)
	})

	mux.HandleFunc("GET "+prefix+"/mysql/{st}/{id}", func(w http.ResponseWriter, r *http.Request) {
		route.get(w, r)
	})

	mux.HandleFunc("GET "+prefix+"/mysql/{st}", func(w http.ResponseWriter, r *http.Request) {
		route.getMany(w, r)
	})

	mux.HandleFunc("POST "+prefix+"/mysql/{st}", func(w http.ResponseWriter, r *http.Request) {
		route.post(w, r)
	})

}

type RouteMySQL struct {
	service *service.ApplicationServiceImpl
}

func (route *RouteMySQL) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	st := r.PathValue("st")
	id := r.PathValue("id")

	err := route.service.Remove(st, "id='"+id+"'")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response := schema.CreateHTTPResponseRFC9457("删除失败", http.StatusInternalServerError, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := schema.CreateHTTPResponseRFC9457("删除成功", http.StatusOK, r)
	json.NewEncoder(w).Encode(response)
}

func (route *RouteMySQL) put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	st := r.PathValue("st")
	id := r.PathValue("id")
	d := r.URL.Query().Get("d")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := schema.CreateHTTPResponseRFC9457("无效的请求体", http.StatusBadRequest, r)
		json.NewEncoder(w).Encode(response)
		return
	}
	data["id"] = id

	deprecated := false
	if d == "1" || d == "true" {
		deprecated = true
	}
	err := route.service.Update(st, data, "id='"+id+"'", deprecated)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response := schema.CreateHTTPResponseRFC9457("更新失败", http.StatusInternalServerError, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := schema.CreateHTTPResponseRFC9457("更新成功", http.StatusOK, r)
	json.NewEncoder(w).Encode(response)
}

func (route *RouteMySQL) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	st := r.PathValue("st")
	id := r.PathValue("id")

	result, err := route.service.Get(st, [][]string{{"id='" + id + "'"}}, "")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response := schema.CreateHTTPResponseRFC9457("内部服务器错误", http.StatusInternalServerError, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (route *RouteMySQL) getMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	st := r.PathValue("st")
	last := r.URL.Query().Get("l")
	filter := r.URL.Query().Get("f")
	f, err := utility.ConvertQueryStringToDefaultFilter(filter)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := schema.CreateHTTPResponseRFC9457("无效的查询参数", http.StatusBadRequest, r)
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println(f)
	columns := r.URL.Query().Get("c")
	var c []string
	if columns == "" {
		c = []string{}
	} else {
		c = strings.Split(columns, ",")
	}

	result, err := route.service.GetMany(st, c, f, last)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response := schema.CreateHTTPResponseRFC9457("内部服务器错误", http.StatusInternalServerError, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (route *RouteMySQL) post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	st := r.PathValue("st")

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := schema.CreateHTTPResponseRFC9457("无效的请求体", http.StatusBadRequest, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := route.service.Create(st, data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response := schema.CreateHTTPResponseRFC9457("创建失败", http.StatusInternalServerError, r)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := schema.CreateHTTPResponseRFC9457("创建成功", http.StatusCreated, r)
	json.NewEncoder(w).Encode(response)
}
